// Package flowspec provides a YAML DSL for defining and executing Temporal workflows.
//
// Workflows are declared as YAML files containing a name, optional trigger, and an
// ordered list of steps. Each step references an agent by name and can pass data from
// the previous step's output using $prev.* syntax. Steps may also fan out into
// parallel sub-steps or execute conditionally.
//
// The package provides parsing ([Parse], [ParseFile]) and validation ([Validate])
// but does not include execution logic. A separate engine compiles parsed specs into
// Temporal workflow definitions at runtime.
//
// Example YAML:
//
//	name: morning-briefing
//	trigger:
//	  cron: "0 8 * * *"
//	steps:
//	  - name: ingest
//	    agent: ingestor
//	    config:
//	      sources: [hn, rss-arstechnica]
//	  - name: deliver
//	    agent: broadcaster
//	    input: "$prev.briefing"
package flowspec

// Workflow represents a complete workflow definition parsed from YAML.
//
// A Workflow has a unique name, an optional trigger that determines when it runs,
// and an ordered list of steps that execute sequentially. The OnError field sets
// the default error strategy for all steps in the workflow.
//
// Example:
//
//	wf := &flowspec.Workflow{
//	    Name: "my-pipeline",
//	    Trigger: flowspec.Trigger{Cron: "0 8 * * *"},
//	    Steps: []flowspec.Step{
//	        {Name: "fetch", Agent: "fetcher"},
//	        {Name: "process", Agent: "processor", Input: "$prev.data"},
//	    },
//	    OnError: "abort",
//	}
type Workflow struct {
	// Name is the unique identifier for this workflow. Required.
	Name string `yaml:"name"`

	// Trigger defines when the workflow should be executed.
	// Optional -- workflows without triggers must be started explicitly.
	Trigger Trigger `yaml:"trigger,omitempty"`

	// Steps is the ordered list of steps to execute. At least one step is required.
	Steps []Step `yaml:"steps"`

	// OnError sets the default error handling strategy for all steps.
	// Valid values are "retry", "skip", and "abort". Individual steps can
	// override this with their own OnError field.
	OnError string `yaml:"on_error,omitempty"`
}

// Trigger defines when a workflow should be executed.
//
// At most one field should be set. A Trigger with no fields set means the
// workflow must be started manually or programmatically.
//
// Example (cron):
//
//	trigger:
//	  cron: "0 8 * * *"
//
// Example (event):
//
//	trigger:
//	  event: gitea.pull_request.opened
type Trigger struct {
	// Cron is a cron expression that schedules the workflow (e.g., "0 8 * * *").
	Cron string `yaml:"cron,omitempty"`

	// Event is an event name that triggers the workflow (e.g., "gitea.pull_request.opened").
	Event string `yaml:"event,omitempty"`

	// Webhook is a webhook endpoint path that triggers the workflow.
	Webhook string `yaml:"webhook,omitempty"`
}

// Step represents a single unit of work in a workflow.
//
// A Step must have either an Agent (for a single activity) or Parallel sub-steps
// (for fan-out), but not neither. Step names must be unique within a workflow.
//
// Example (agent step):
//
//	step := flowspec.Step{
//	    Name:  "fetch-data",
//	    Agent: "fetcher",
//	    Config: map[string]any{"url": "https://example.com"},
//	    Timeout: "30s",
//	}
//
// Example (parallel fan-out):
//
//	step := flowspec.Step{
//	    Name: "fan-out",
//	    Parallel: []flowspec.Step{
//	        {Name: "task-a", Agent: "worker-a"},
//	        {Name: "task-b", Agent: "worker-b"},
//	    },
//	}
type Step struct {
	// Name is the unique identifier for this step within the workflow. Required.
	Name string `yaml:"name"`

	// Agent is the name of the agent that executes this step.
	// Required unless Parallel is set.
	Agent string `yaml:"agent"`

	// Input is an expression referencing output from a previous step.
	// Typically uses $prev.* syntax (e.g., "$prev.items").
	Input string `yaml:"input,omitempty"`

	// Config is arbitrary key-value configuration passed to the agent.
	Config map[string]any `yaml:"config,omitempty"`

	// Timeout is the maximum duration for this step (e.g., "5m", "1h").
	Timeout string `yaml:"timeout,omitempty"`

	// OnError sets the error handling strategy for this step, overriding
	// the workflow-level default. Valid values: "retry", "skip", "abort".
	OnError string `yaml:"on_error,omitempty"`

	// Parallel contains sub-steps to execute concurrently (fan-out pattern).
	// When set, Agent should be empty.
	Parallel []Step `yaml:"parallel,omitempty"`

	// Condition is a boolean expression for conditional execution
	// (e.g., "$prev.urgent == true"). If the condition evaluates to false,
	// the step is skipped.
	Condition string `yaml:"condition,omitempty"`
}
