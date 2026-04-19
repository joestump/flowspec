// Package flowspec provides a YAML DSL for defining and executing Temporal workflows.
package flowspec

// Workflow represents a complete workflow definition.
type Workflow struct {
	Name    string  `yaml:"name"`
	Trigger Trigger `yaml:"trigger,omitempty"`
	Steps   []Step  `yaml:"steps"`
	OnError string  `yaml:"on_error,omitempty"` // "retry", "skip", "abort"
}

// Trigger defines when a workflow should be executed.
type Trigger struct {
	Cron    string `yaml:"cron,omitempty"`
	Event   string `yaml:"event,omitempty"`
	Webhook string `yaml:"webhook,omitempty"`
}

// Step represents a single step in a workflow.
type Step struct {
	Name      string         `yaml:"name"`
	Agent     string         `yaml:"agent"`
	Input     string         `yaml:"input,omitempty"`     // "$prev.items" or literal
	Config    map[string]any `yaml:"config,omitempty"`
	Timeout   string         `yaml:"timeout,omitempty"`   // "5m", "1h"
	OnError   string         `yaml:"on_error,omitempty"`
	Parallel  []Step         `yaml:"parallel,omitempty"`  // fan-out
	Condition string         `yaml:"condition,omitempty"` // "$prev.urgent == true"
}
