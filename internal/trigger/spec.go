package trigger

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

type TestPrompt struct {
	Prompt     string `yaml:"prompt"`
	Reason     string `yaml:"reason"`
	Confidence string `yaml:"confidence"`
}

type TestSpec struct {
	Skill                   string       `yaml:"skill"`
	ShouldTriggerPrompts    []TestPrompt `yaml:"should_trigger_prompts"`
	ShouldNotTriggerPrompts []TestPrompt `yaml:"should_not_trigger_prompts"`
}

var validConfidences = map[string]bool{"": true, "high": true, "medium": true}

func ParseSpec(data []byte) (*TestSpec, error) {
	var s TestSpec
	if err := yaml.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parsing trigger tests: %w", err)
	}
	if s.Skill == "" {
		return nil, errors.New("trigger tests missing required 'skill' field")
	}
	if len(s.ShouldTriggerPrompts)+len(s.ShouldNotTriggerPrompts) == 0 {
		return nil, errors.New("trigger tests must define at least one prompt")
	}
	for _, p := range s.ShouldTriggerPrompts {
		if p.Prompt == "" {
			return nil, errors.New("trigger test prompt missing required 'prompt' field")
		}
		if !validConfidences[p.Confidence] {
			return nil, fmt.Errorf("prompt %q has unrecognized confidence %q (valid: high, medium)", p.Prompt, p.Confidence)
		}
	}
	for _, p := range s.ShouldNotTriggerPrompts {
		if p.Prompt == "" {
			return nil, errors.New("trigger test prompt missing required 'prompt' field")
		}
		if !validConfidences[p.Confidence] {
			return nil, fmt.Errorf("prompt %q has unrecognized confidence %q (valid: high, medium)", p.Prompt, p.Confidence)
		}
	}
	return &s, nil
}
