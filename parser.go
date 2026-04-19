package flowspec

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Parse reads a YAML workflow definition from raw bytes.
func Parse(data []byte) (*Workflow, error) {
	var wf Workflow
	if err := yaml.Unmarshal(data, &wf); err != nil {
		return nil, fmt.Errorf("flowspec: parse error: %w", err)
	}
	return &wf, nil
}

// ParseFile reads a YAML workflow definition from a file path.
func ParseFile(path string) (*Workflow, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("flowspec: read file: %w", err)
	}
	return Parse(data)
}

// Validate checks a parsed workflow for structural correctness.
func Validate(w *Workflow) error {
	if w.Name == "" {
		return fmt.Errorf("flowspec: workflow name is required")
	}
	if len(w.Steps) == 0 {
		return fmt.Errorf("flowspec: workflow %q must have at least one step", w.Name)
	}

	seen := make(map[string]bool)
	for i, step := range w.Steps {
		if err := validateStep(step, i, seen); err != nil {
			return err
		}
	}
	return nil
}

func validateStep(s Step, index int, seen map[string]bool) error {
	if s.Name == "" {
		return fmt.Errorf("flowspec: step %d must have a name", index)
	}
	if seen[s.Name] {
		return fmt.Errorf("flowspec: duplicate step name %q", s.Name)
	}
	seen[s.Name] = true

	// A step must have either an agent or parallel sub-steps, but not neither.
	if s.Agent == "" && len(s.Parallel) == 0 {
		return fmt.Errorf("flowspec: step %q must specify an agent or parallel sub-steps", s.Name)
	}

	if s.OnError != "" {
		switch s.OnError {
		case "retry", "skip", "abort":
			// valid
		default:
			return fmt.Errorf("flowspec: step %q has invalid on_error value %q (must be retry, skip, or abort)", s.Name, s.OnError)
		}
	}

	for i, ps := range s.Parallel {
		if err := validateStep(ps, i, seen); err != nil {
			return err
		}
	}
	return nil
}
