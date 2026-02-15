package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spboyer/waza/internal/cache"
	"github.com/spboyer/waza/internal/config"
	"github.com/spboyer/waza/internal/execution"
	"github.com/spboyer/waza/internal/models"
	"github.com/spboyer/waza/internal/orchestration"
	"github.com/spboyer/waza/internal/reporting"
	"github.com/spboyer/waza/internal/utils"
	"github.com/spf13/cobra"
)

var (
	contextDir     string
	outputPath     string
	verbose        bool
	transcriptDir  string
	taskFilters    []string
	parallel       bool
	workers        int
	interpret      bool
	format         string
	enableCache    bool
	disableCache   bool
	runCacheDir    string
	modelOverrides []string
)

// modelResult pairs a model identifier with its evaluation outcome.
type modelResult struct {
	modelID string
	outcome *models.EvaluationOutcome
}

func newRunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run <eval.yaml>",
		Short: "Run an evaluation benchmark",
		Long: `Run an evaluation benchmark from a spec file.

The spec file defines the benchmark configuration, test cases, and validation rules.
Resources are loaded from the context directory (defaults to ./fixtures).`,
		Args: cobra.ExactArgs(1),
		RunE: runCommandE,
	}

	cmd.Flags().StringVar(&contextDir, "context-dir", "", "Context directory for fixtures (default: ./fixtures relative to spec)")
	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output JSON file for results")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output with detailed progress")
	cmd.Flags().StringVar(&transcriptDir, "transcript-dir", "", "Directory to save per-task transcript JSON files")
	cmd.Flags().StringArrayVar(&taskFilters, "task", nil, "Filter tasks by name/ID glob pattern (can be repeated)")
	cmd.Flags().BoolVar(&parallel, "parallel", false, "Run tasks concurrently")
	cmd.Flags().IntVar(&workers, "workers", 0, "Number of concurrent workers (default: 4, requires --parallel)")
	cmd.Flags().BoolVar(&interpret, "interpret", false, "Print a plain-language interpretation of the results")
	cmd.Flags().StringVar(&format, "format", "default", "Output format: default, github-comment")
	cmd.Flags().BoolVar(&enableCache, "cache", false, "Enable result caching (default: false)")
	cmd.Flags().BoolVar(&disableCache, "no-cache", false, "Disable result caching (default)")
	cmd.Flags().StringVar(&runCacheDir, "cache-dir", ".waza-cache", "Cache directory for storing results")
	cmd.Flags().StringArrayVar(&modelOverrides, "model", nil, "Model to use (overrides spec config, can be repeated for comparison)")

	return cmd
}

func runCommandE(cmd *cobra.Command, args []string) error {
	specPath := args[0]

	// Load spec
	spec, err := models.LoadBenchmarkSpec(specPath)
	if err != nil {
		return fmt.Errorf("failed to load spec: %w", err)
	}

	// CLI flags override spec config
	if parallel {
		spec.Config.Concurrent = true
	}
	if workers > 0 {
		spec.Config.Workers = workers
	}

	// Determine the list of models to evaluate
	modelsToRun := []string{spec.Config.ModelID}
	if len(modelOverrides) > 0 {
		modelsToRun = modelOverrides
	}
	multiModel := len(modelsToRun) > 1

	// Run evaluation for each model, collecting results
	var allResults []modelResult
	var lastErr error

	for _, modelID := range modelsToRun {
		// Override spec model for this iteration
		spec.Config.ModelID = modelID

		outcome, err := runSingleModel(cmd, spec, specPath)
		if err != nil {
			var testErr *TestFailureError
			if errors.As(err, &testErr) {
				// Test failures are recorded but don't stop a multi-model run
				allResults = append(allResults, modelResult{modelID: modelID, outcome: outcome})
				lastErr = err
				continue
			}
			return err
		}
		allResults = append(allResults, modelResult{modelID: modelID, outcome: outcome})
	}

	// Print comparison table when multiple models were evaluated
	if multiModel && len(allResults) > 0 {
		printModelComparison(allResults)
	}

	// Save per-model results when --output is specified with multiple models
	if outputPath != "" && multiModel {
		ext := filepath.Ext(outputPath)
		base := strings.TrimSuffix(outputPath, ext)
		for _, mr := range allResults {
			perModelPath := fmt.Sprintf("%s_%s%s", base, sanitizeModelName(mr.modelID), ext)
			if err := saveOutcome(mr.outcome, perModelPath); err != nil {
				return fmt.Errorf("failed to save output for model %s: %w", mr.modelID, err)
			}
			fmt.Printf("Results saved to: %s\n", perModelPath)
		}
	}

	if lastErr != nil {
		return lastErr
	}

	return nil
}

// runSingleModel executes a benchmark for one model and returns the outcome.
// It prints the per-model summary and saves output for single-model runs.
func runSingleModel(_ *cobra.Command, spec *models.BenchmarkSpec, specPath string) (*models.EvaluationOutcome, error) {
	// Get spec directory for resolving relative paths
	specDir := filepath.Dir(specPath)
	if !filepath.IsAbs(specDir) {
		absSpecDir, err := filepath.Abs(specDir)
		if err == nil {
			specDir = absSpecDir
		}
	}

	// Resolve fixture/context dir relative to spec file if not absolute
	fixtureDir := contextDir
	if fixtureDir == "" {
		fixtureDir = filepath.Join(specDir, "fixtures")
	} else if !filepath.IsAbs(fixtureDir) {
		absFixtureDir, err := filepath.Abs(fixtureDir)
		if err == nil {
			fixtureDir = absFixtureDir
		}
	}

	// Create config with both directories
	cfg := config.NewBenchmarkConfig(spec,
		config.WithSpecDir(specDir),
		config.WithFixtureDir(fixtureDir),
		config.WithVerbose(verbose),
		config.WithOutputPath(outputPath),
		config.WithTranscriptDir(transcriptDir),
	)

	// Setup cache if enabled
	var resultCache *cache.Cache
	useCaching := enableCache && !disableCache

	if useCaching && cache.HasNonDeterministicGraders(spec) {
		if verbose {
			fmt.Println("Note: Caching disabled due to non-deterministic graders (behavior, prompt)")
		}
		useCaching = false
	}

	if useCaching {
		absCacheDir, err := filepath.Abs(runCacheDir)
		if err != nil {
			return nil, fmt.Errorf("resolving cache directory: %w", err)
		}
		resultCache = cache.New(absCacheDir)
		if verbose {
			fmt.Printf("Cache enabled: %s\n", absCacheDir)
		}
	}

	// Create engine based on spec
	var engine execution.AgentEngine

	switch spec.Config.EngineType {
	case "mock":
		engine = execution.NewMockEngine(spec.Config.ModelID)
	case "copilot-sdk":
		engine = execution.NewCopilotEngineBuilder(spec.Config.ModelID).Build()
	default:
		return nil, fmt.Errorf("unknown engine type: %s", spec.Config.EngineType)
	}

	// Create runner with optional task filters and cache
	runnerOpts := []orchestration.RunnerOption{
		orchestration.WithTaskFilters(taskFilters...),
	}
	if resultCache != nil {
		runnerOpts = append(runnerOpts, orchestration.WithCache(resultCache))
	}
	runner := orchestration.NewTestRunner(cfg, engine, runnerOpts...)

	// Add progress listener
	if verbose {
		runner.OnProgress(verboseProgressListener)
	} else {
		runner.OnProgress(simpleProgressListener)
	}

	// Run benchmark
	ctx := context.Background()

	fmt.Printf("Running benchmark: %s\n", spec.Name)
	fmt.Printf("Skill: %s\n", spec.SkillName)
	fmt.Printf("Engine: %s\n", spec.Config.EngineType)
	fmt.Printf("Model: %s\n", spec.Config.ModelID)
	if spec.Config.Concurrent {
		w := spec.Config.Workers
		if w <= 0 {
			w = 4
		}
		fmt.Printf("Parallel: %d workers\n", w)
	}

	if verbose && len(spec.Config.SkillPaths) > 0 {
		fmt.Printf("Skill Directories:\n")
		resolvedPaths := utils.ResolvePaths(spec.Config.SkillPaths, specDir)
		for _, path := range resolvedPaths {
			fmt.Printf("  - %s\n", path)
		}
	}

	fmt.Println()

	outcome, err := runner.RunBenchmark(ctx)
	if err != nil {
		return nil, fmt.Errorf("benchmark failed: %w", err)
	}

	// Print results based on format
	switch format {
	case "github-comment":
		fmt.Print(FormatGitHubComment(outcome))
	case "default":
		printSummary(outcome)

		if interpret {
			fmt.Println()
			fmt.Print(reporting.FormatSummaryReport(outcome))
		}
	default:
		return nil, fmt.Errorf("unknown output format: %s (supported: default, github-comment)", format)
	}

	// Save output for single-model runs (multi-model saves are handled by the caller)
	if outputPath != "" && len(modelOverrides) <= 1 {
		if err := saveOutcome(outcome, outputPath); err != nil {
			return nil, fmt.Errorf("failed to save output: %w", err)
		}
		fmt.Printf("\nResults saved to: %s\n", outputPath)
	}

	// Return test failure as error so caller can decide how to handle it
	if outcome.Digest.Failed > 0 || outcome.Digest.Errors > 0 {
		return outcome, &TestFailureError{
			Message: fmt.Sprintf("benchmark completed with %d failed and %d error(s)", outcome.Digest.Failed, outcome.Digest.Errors),
		}
	}

	return outcome, nil
}

// printModelComparison renders a comparison table for multi-model runs.
func printModelComparison(results []modelResult) {
	fmt.Println()
	fmt.Println("═" + strings.Repeat("═", 54))
	fmt.Println(" MODEL COMPARISON")
	fmt.Println("═" + strings.Repeat("═", 54))
	fmt.Println()
	fmt.Printf("%-20s %-8s %-10s %s\n", "Model", "Score", "Pass Rate", "Duration")
	fmt.Println("─" + strings.Repeat("─", 54))

	for _, mr := range results {
		score := 0.0
		passRate := 0.0
		durationMs := int64(0)
		if mr.outcome != nil {
			score = mr.outcome.Digest.AggregateScore
			passRate = mr.outcome.Digest.SuccessRate * 100
			durationMs = mr.outcome.Digest.DurationMs
		}
		duration := time.Duration(durationMs) * time.Millisecond
		passStr := fmt.Sprintf("%.1f%%", passRate)
		fmt.Printf("%-20s %-8.2f %-10s %v\n", mr.modelID, score, passStr, duration)
	}
	fmt.Println()
}

// sanitizeModelName replaces characters that are invalid in filenames.
func sanitizeModelName(name string) string {
	r := strings.NewReplacer("/", "-", "\\", "-", ":", "-", " ", "-")
	return r.Replace(name)
}

func verboseProgressListener(event orchestration.ProgressEvent) {
	switch event.EventType {
	case orchestration.EventBenchmarkStart:
		fmt.Printf("Starting benchmark with %d test(s)...\n\n", event.TotalTests)
	case orchestration.EventTestStart:
		fmt.Printf("[%d/%d] Running test: %s\n", event.TestNum, event.TotalTests, event.TestName)
	case orchestration.EventTestCached:
		fmt.Printf("[%d/%d] Test: %s [cached]\n\n", event.TestNum, event.TotalTests, event.TestName)
	case orchestration.EventRunStart:
		fmt.Printf("  Run %d/%d...", event.RunNum, event.TotalRuns)
	case orchestration.EventRunComplete:
		duration := time.Duration(event.DurationMs) * time.Millisecond
		fmt.Printf(" %s (%v)\n", event.Status, duration)
	case orchestration.EventTestComplete:
		fmt.Printf("  Test %s: %s\n\n", event.TestName, event.Status)
	case orchestration.EventBenchmarkComplete:
		duration := time.Duration(event.DurationMs) * time.Millisecond
		fmt.Printf("Benchmark completed in %v\n\n", duration)
	case orchestration.EventAgentPrompt:
		if msg, ok := event.Details["message"].(string); ok {
			fmt.Printf("  [PROMPT] %s\n", msg)
		}
	case orchestration.EventAgentResponse:
		if output, ok := event.Details["output"].(string); ok && output != "" {
			fmt.Printf("  [RESPONSE] %s\n", truncate(output, 200))
		}
		if tc, ok := event.Details["tool_calls"].(int); ok && tc > 0 {
			fmt.Printf("  [TOOLS] %d tool call(s)\n", tc)
		}
		if e, ok := event.Details["error"].(string); ok && e != "" {
			fmt.Printf("  [ERROR] %s\n", e)
		}
	case orchestration.EventGraderResult:
		name := fmt.Sprintf("%v", event.Details["grader"])
		passed, ok := event.Details["passed"].(bool)
		if !ok {
			passed = false
		}
		score, ok := event.Details["score"].(float64)
		if !ok {
			score = 0
		}
		feedback := fmt.Sprintf("%v", event.Details["feedback"])
		icon := "✗"
		if passed {
			icon = "✓"
		}
		duration := time.Duration(event.DurationMs) * time.Millisecond
		fmt.Printf("  [GRADER] %s %s score=%.2f (%v)", icon, name, score, duration)
		if feedback != "" {
			fmt.Printf(" — %s", feedback)
		}
		fmt.Println()
	}
}

// truncate shortens s to maxLen characters, appending "..." if truncated.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func simpleProgressListener(event orchestration.ProgressEvent) {
	switch event.EventType {
	case orchestration.EventTestCached:
		fmt.Printf("✓ [%d/%d] %s [cached]\n", event.TestNum, event.TotalTests, event.TestName)
	case orchestration.EventTestComplete:
		status := "✓"
		if event.Status != models.StatusPassed {
			status = "✗"
		}
		fmt.Printf("%s [%d/%d] %s\n", status, event.TestNum, event.TotalTests, event.TestName)
	}
}

func printSummary(outcome *models.EvaluationOutcome) {
	fmt.Println("=" + strings.Repeat("=", 50))
	fmt.Println(" BENCHMARK RESULTS")
	fmt.Println("=" + strings.Repeat("=", 50))
	fmt.Println()

	digest := outcome.Digest

	fmt.Printf("Total Tests:    %d\n", digest.TotalTests)
	fmt.Printf("Succeeded:      %d\n", digest.Succeeded)
	fmt.Printf("Failed:         %d\n", digest.Failed)
	fmt.Printf("Errors:         %d\n", digest.Errors)
	fmt.Printf("Success Rate:   %.1f%%\n", digest.SuccessRate*100)
	fmt.Printf("Aggregate Score: %.2f\n", digest.AggregateScore)
	fmt.Printf("Min Score:      %.2f\n", digest.MinScore)
	fmt.Printf("Max Score:      %.2f\n", digest.MaxScore)
	fmt.Printf("Std Dev:        %.4f\n", digest.StdDev)

	duration := time.Duration(digest.DurationMs) * time.Millisecond
	fmt.Printf("Duration:       %v\n", duration)
	fmt.Println()

	// Per-task breakdown
	fmt.Println("-" + strings.Repeat("-", 50))
	fmt.Println(" PER-TASK BREAKDOWN")
	fmt.Println("-" + strings.Repeat("-", 50))
	for _, to := range outcome.TestOutcomes {
		icon := "✓"
		if to.Status != models.StatusPassed {
			icon = "✗"
		}
		fmt.Printf("  %s %s [%s]\n", icon, to.DisplayName, to.Status)
		if to.Stats != nil {
			fmt.Printf("      pass_rate=%.1f%%  avg=%.2f  min=%.2f  max=%.2f  stddev=%.4f  avg_dur=%dms\n",
				to.Stats.PassRate*100, to.Stats.AvgScore,
				to.Stats.MinScore, to.Stats.MaxScore,
				to.Stats.StdDevScore, to.Stats.AvgDurationMs)
		}
	}
	fmt.Println()

	// Show failed tests
	if digest.Failed > 0 || digest.Errors > 0 {
		fmt.Println("Failed Tests:")
		for _, to := range outcome.TestOutcomes {
			if to.Status != models.StatusPassed {
				fmt.Printf("  - %s (%s)\n", to.DisplayName, to.Status)

				// Show validation failures
				if len(to.Runs) > 0 {
					for _, run := range to.Runs {
						for _, val := range run.Validations {
							if !val.Passed {
								fmt.Printf("    • %s: %s\n", val.Name, val.Feedback)
							}
						}
					}
				}
			}
		}
		fmt.Println()
	}

	// Show flaky tasks
	var flakyTasks []models.TestOutcome
	for _, to := range outcome.TestOutcomes {
		if to.Stats != nil && to.Stats.Flaky {
			flakyTasks = append(flakyTasks, to)
		}
	}
	if len(flakyTasks) > 0 {
		fmt.Println("\u26a0 Flaky Tasks (inconsistent pass/fail across trials):")
		for _, to := range flakyTasks {
			fmt.Printf("  - %s  pass_rate=%.0f%%  score=%.2f\u00b1%.2f  CI95=[%.2f, %.2f]\n",
				to.DisplayName,
				to.Stats.PassRate*100,
				to.Stats.AvgScore,
				to.Stats.StdDevScore,
				to.Stats.CI95Lo,
				to.Stats.CI95Hi,
			)
		}
		fmt.Println()
	}
}

func saveOutcome(outcome *models.EvaluationOutcome, path string) error {
	data, err := json.MarshalIndent(outcome, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
