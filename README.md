# flowspec

![CI](https://github.com/joestump/flowspec/actions/workflows/ci.yml/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/joestump/flowspec.svg)](https://pkg.go.dev/github.com/joestump/flowspec)

**A YAML DSL for defining and executing Temporal workflows.**

Declarative workflow definitions that compile to Temporal workflows. Define pipelines as YAML, execute them with a generic Go interpreter. Designed for AI agent orchestration but general-purpose.

> **Status:** Early development

**[Documentation](https://joestump.github.io/flowspec/docs/intro)** | **[pkg.go.dev](https://pkg.go.dev/github.com/joestump/flowspec)**

## Example

```yaml
name: morning-briefing
trigger:
  cron: "0 8 * * *"
steps:
  - name: ingest
    agent: ingestor
    config:
      sources: [hn, rss-arstechnica, rss-lobsters]
  - name: read-and-score
    agent: reader
    input: "$prev.items"
    config:
      max_items: 20
  - name: compose
    agent: editor
    input: "$prev.scored_items"
    config:
      layout: editorial-hero
  - name: deliver
    agent: broadcaster
    input: "$prev.briefing"
    config:
      channels: [telegram, web]
```

## Architecture

flowspec has three layers:

1. **Spec** (`spec.go`) --- Core types representing workflows, steps, triggers, and configuration. Pure data structures with YAML tags.
2. **Parser** (`parser.go`) --- Reads YAML files into spec types and validates them. No execution logic.
3. **Engine** (planned) --- Compiles parsed specs into Temporal workflow definitions and executes them. Steps map to Temporal activities, `parallel` blocks become fan-out patterns, and `$prev.*` references wire step outputs to inputs.

### Step References

Steps can reference outputs from previous steps using `$prev.*` syntax in their `input` field:

- `$prev.items` --- the `items` field from the previous step's output
- `$prev.scored_items` --- the `scored_items` field from the previous step's output

### Error Handling

Each step (and the workflow itself) can specify an `on_error` strategy:

- `retry` --- retry the step with backoff
- `skip` --- skip the step and continue
- `abort` --- stop the workflow immediately

### Parallel Execution

Steps can contain a `parallel` block for fan-out:

```yaml
steps:
  - name: fan-out
    parallel:
      - name: task-a
        agent: worker-a
      - name: task-b
        agent: worker-b
```

## Installation

```bash
go get github.com/joestump/flowspec
```

## Usage

```go
package main

import (
    "fmt"
    "log"

    "github.com/joestump/flowspec"
)

func main() {
    wf, err := flowspec.ParseFile("examples/morning-briefing.yaml")
    if err != nil {
        log.Fatal(err)
    }

    if err := flowspec.Validate(wf); err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Workflow: %s (%d steps)\n", wf.Name, len(wf.Steps))
}
```

## License

MIT
