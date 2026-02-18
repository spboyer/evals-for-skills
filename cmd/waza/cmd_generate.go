package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spboyer/waza/internal/generate"
	"github.com/spboyer/waza/internal/workspace"
	"github.com/spf13/cobra"
)

var generateOutputDir string

func newGenerateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate <skill-name | SKILL.md>",
		Short: "Generate an eval suite from a SKILL.md file",
		Long: `Generate evaluation files from a SKILL.md file.

Parses the YAML frontmatter (name, description) from the given SKILL.md and
creates an eval.yaml, starter task files, and a fixtures directory.

If the argument looks like a skill name (no path separators or file extension),
it is resolved via workspace detection to find the SKILL.md path.`,
		Args: cobra.ExactArgs(1),
		RunE: generateCommandE,
	}

	cmd.Flags().StringVarP(&generateOutputDir, "output-dir", "d", "", "Output directory (default: ./eval-{skill-name}/)")

	return cmd
}

func generateCommandE(_ *cobra.Command, args []string) error {
	skillPath := args[0]

	// If arg looks like a skill name (not a path), resolve via workspace
	if !workspace.LooksLikePath(skillPath) {
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("getting working directory: %w", err)
		}
		ctx, err := workspace.DetectContext(wd)
		if err != nil {
			return fmt.Errorf("detecting workspace: %w", err)
		}
		si, err := workspace.FindSkill(ctx, skillPath)
		if err != nil {
			return err
		}
		skillPath = si.SkillPath
	}

	skill, err := generate.ParseSkillMD(skillPath)
	if err != nil {
		return fmt.Errorf("failed to parse SKILL.md: %w", err)
	}

	outDir := generateOutputDir
	if outDir == "" {
		outDir = filepath.Join(".", fmt.Sprintf("eval-%s", skill.Name))
	}

	fmt.Printf("Generating eval suite for skill: %s\n", skill.Name)
	fmt.Printf("Output directory: %s\n", outDir)

	if err := generate.GenerateEvalSuite(skill, outDir); err != nil {
		return fmt.Errorf("failed to generate eval suite: %w", err)
	}

	fmt.Println()
	fmt.Println("Generated files:")
	fmt.Printf("  %s/eval.yaml\n", outDir)
	fmt.Printf("  %s/tasks/%s-basic.yaml\n", outDir, skill.Name)
	fmt.Printf("  %s/fixtures/sample.txt\n", outDir)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Printf("  1. Edit the task files in %s/tasks/\n", outDir)
	fmt.Printf("  2. Add real fixtures to %s/fixtures/\n", outDir)
	fmt.Printf("  3. Run: waza run %s/eval.yaml\n", outDir)

	return nil
}
