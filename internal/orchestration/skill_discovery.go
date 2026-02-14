package orchestration

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spboyer/waza/internal/skill"
)

// discoverSkills scans the provided directories for SKILL.md files
// and returns a map of skill names to their file paths.
func discoverSkills(directories []string) (map[string]string, error) {
	discovered := make(map[string]string)

	for _, dir := range directories {
		// Check if directory exists
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			// Directory doesn't exist, skip it
			continue
		}

		// Look for SKILL.md in this directory
		skillPath := filepath.Join(dir, "SKILL.md")
		if _, err := os.Stat(skillPath); err == nil {
			// SKILL.md exists, parse it
			skillName, err := parseSkillName(skillPath)
			if err != nil {
				// Log warning but continue - we want to discover as many skills as possible
				fmt.Fprintf(os.Stderr, "warning: failed to parse %s: %v\n", skillPath, err)
				continue
			}
			if skillName != "" {
				discovered[skillName] = skillPath
			}
		}
	}

	return discovered, nil
}

// parseSkillName reads a SKILL.md file and extracts the skill name from frontmatter.
func parseSkillName(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("reading file: %w", err)
	}

	var s skill.Skill
	if err := s.UnmarshalText(data); err != nil {
		return "", fmt.Errorf("parsing SKILL.md: %w", err)
	}

	return strings.TrimSpace(s.Frontmatter.Name), nil
}

// validateRequiredSkills checks that all required skills are discovered.
// Returns nil if validation passes, or an error describing what's missing.
func validateRequiredSkills(requiredSkills []string, discoveredSkills map[string]string, searchedDirs []string) error {
	if len(requiredSkills) == 0 {
		// No required skills specified, validation passes
		return nil
	}

	var missing []string
	for _, required := range requiredSkills {
		if _, found := discoveredSkills[required]; !found {
			missing = append(missing, required)
		}
	}

	if len(missing) > 0 {
		// Build error message
		var sb strings.Builder
		sb.WriteString("required skills not found:\n")
		for _, skillName := range missing {
			sb.WriteString(fmt.Sprintf("  - %s\n", skillName))
		}
		sb.WriteString("\nSearched directories:\n")
		for _, dir := range searchedDirs {
			sb.WriteString(fmt.Sprintf("  - %s\n", dir))
		}
		if len(discoveredSkills) > 0 {
			sb.WriteString("\nFound skills:\n")
			for skillName := range discoveredSkills {
				sb.WriteString(fmt.Sprintf("  - %s\n", skillName))
			}
		} else {
			sb.WriteString("\nNo skills were found in the searched directories.\n")
		}
		return errors.New(sb.String())
	}

	return nil
}
