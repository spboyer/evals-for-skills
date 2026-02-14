package main

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/spboyer/waza/internal/models"
)

// formatDuration formats a duration in a consistent, human-readable way.
// This ensures stable output regardless of Go version changes.
func formatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	// Use the built-in formatting but ensure we control it
	return d.String()
}

// FormatGitHubComment formats an EvaluationOutcome as a markdown comment for GitHub PRs
func FormatGitHubComment(outcome *models.EvaluationOutcome) string {
	var b strings.Builder

	digest := outcome.Digest
	duration := time.Duration(digest.DurationMs) * time.Millisecond

	// Header with overall status
	b.WriteString("## 🧪 Waza Eval Results\n\n")

	// Overall status badge
	statusIcon := "✅ Passed"
	if digest.Failed > 0 || digest.Errors > 0 {
		statusIcon = "❌ Failed"
	}

	b.WriteString(fmt.Sprintf("**Status:** %s | **Score:** %.2f | **Duration:** %s\n\n",
		statusIcon, digest.AggregateScore, formatDuration(duration)))

	// Summary stats
	b.WriteString(fmt.Sprintf("- **Tests:** %d total, %d passed, %d failed, %d errors\n",
		digest.TotalTests, digest.Succeeded, digest.Failed, digest.Errors))
	b.WriteString(fmt.Sprintf("- **Success Rate:** %.1f%%\n", digest.SuccessRate*100))
	b.WriteString(fmt.Sprintf("- **Score Range:** %.2f - %.2f (σ=%.4f)\n\n",
		digest.MinScore, digest.MaxScore, digest.StdDev))

	// Per-task breakdown table
	b.WriteString("### Task Results\n\n")
	b.WriteString("| Task | Score | Status | Graders |\n")
	b.WriteString("|------|-------|--------|----------|\n")

	for _, to := range outcome.TestOutcomes {
		statusIcon := "✅"
		if to.Status != models.StatusPassed {
			statusIcon = "❌"
		}

		// Calculate average score across runs
		avgScore := 0.0
		if to.Stats != nil {
			avgScore = to.Stats.AvgScore
		} else if len(to.Runs) > 0 {
			// Fallback if stats not available
			for _, run := range to.Runs {
				avgScore += run.ComputeRunScore()
			}
			avgScore /= float64(len(to.Runs))
		}

		// Collect grader names from first run
		graderNames := []string{}
		if len(to.Runs) > 0 {
			for name := range to.Runs[0].Validations {
				graderNames = append(graderNames, name)
			}
		}
		// Sort grader names for consistent output
		sort.Strings(graderNames)
		graders := strings.Join(graderNames, ", ")
		if graders == "" {
			graders = "-"
		}

		b.WriteString(fmt.Sprintf("| %s | %.2f | %s | %s |\n",
			to.DisplayName, avgScore, statusIcon, graders))
	}

	b.WriteString("\n")

	// Flaky tasks warning
	var flakyTasks []models.TestOutcome
	for _, to := range outcome.TestOutcomes {
		if to.Stats != nil && to.Stats.Flaky {
			flakyTasks = append(flakyTasks, to)
		}
	}
	if len(flakyTasks) > 0 {
		b.WriteString("### ⚠️ Flaky Tasks\n\n")
		b.WriteString("The following tasks showed inconsistent results across runs:\n\n")
		for _, to := range flakyTasks {
			b.WriteString(fmt.Sprintf("- **%s**: %.0f%% pass rate, score=%.2f±%.2f\n",
				to.DisplayName,
				to.Stats.PassRate*100,
				to.Stats.AvgScore,
				to.Stats.StdDevScore,
			))
		}
		b.WriteString("\n")
	}

	// Grader breakdown for failed tasks
	if digest.Failed > 0 || digest.Errors > 0 {
		b.WriteString("### Failed Task Details\n\n")
		for _, to := range outcome.TestOutcomes {
			if to.Status != models.StatusPassed {
				b.WriteString(fmt.Sprintf("#### %s\n\n", to.DisplayName))

				// Show validation failures from runs
				if len(to.Runs) > 0 {
					for runIdx, run := range to.Runs {
						if run.Status != models.StatusPassed {
							b.WriteString(fmt.Sprintf("**Run %d/%d** (%s):\n",
								runIdx+1, len(to.Runs), run.Status))

							// Collect and sort validation names for consistent output
							valNames := make([]string, 0, len(run.Validations))
							for name := range run.Validations {
								valNames = append(valNames, name)
							}
							sort.Strings(valNames)

							// Print validations in sorted order
							for _, name := range valNames {
								val := run.Validations[name]
								icon := "✅"
								if !val.Passed {
									icon = "❌"
								}
								b.WriteString(fmt.Sprintf("- %s **%s** (%.2f): %s\n",
									icon, val.Name, val.Score, val.Feedback))
							}
							b.WriteString("\n")
						}
					}
				}
			}
		}
	}

	// Footer with metadata
	b.WriteString("---\n\n")
	b.WriteString(fmt.Sprintf("**Benchmark:** %s | **Skill:** %s | **Model:** %s\n",
		outcome.BenchName, outcome.SkillTested, outcome.Setup.ModelID))

	return b.String()
}
