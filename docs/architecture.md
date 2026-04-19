---
layout: default
title: Architecture
---

# Architecture

How flowspec compiles YAML definitions to Temporal workflows.

## Layers

flowspec has three layers:

```
YAML files  -->  Parser  -->  Spec types  -->  Engine  -->  Temporal workflows
```

### 1. Spec (`spec.go`)

Core data structures representing workflows, steps, triggers, and configuration. Pure Go types with YAML struct tags -- no logic, no dependencies beyond the standard library.

- `Workflow` -- top-level container with name, trigger, steps, and error strategy
- `Trigger` -- cron, event, or webhook trigger definitions
- `Step` -- a unit of work: agent reference, input wiring, config, timeouts, error handling, parallel sub-steps, and conditions

### 2. Parser (`parser.go`)

Reads YAML files into spec types and validates structural correctness. The parser enforces:

- Workflows must have a name and at least one step
- Steps must have a name and either an agent or parallel sub-steps
- Step names must be unique within a workflow
- `on_error` values must be one of `retry`, `skip`, or `abort`

The parser is intentionally separate from execution -- you can parse and validate specs without any Temporal dependency.

### 3. Engine (planned)

The engine will compile parsed specs into Temporal workflow definitions:

| flowspec concept | Temporal mapping |
|------------------|-----------------|
| Workflow | Temporal Workflow |
| Step (with agent) | Temporal Activity |
| `parallel` block | Fan-out pattern (goroutines + `workflow.Go`) |
| `$prev.*` references | Activity result wiring |
| `on_error: retry` | Temporal retry policy |
| `on_error: skip` | Error-swallowing wrapper |
| `on_error: abort` | Workflow failure |
| `timeout` | Activity timeout |
| `condition` | Conditional activity execution |
| `trigger.cron` | Temporal cron schedule |
| `trigger.event` | Signal-based workflow start |

## Step References

Steps reference outputs from previous steps using `$prev.*` syntax:

```
step-1 (agent: ingestor)
  output: { items: [...] }
       |
       v
step-2 (agent: reader, input: "$prev.items")
  receives: [...]
  output: { scored_items: [...] }
       |
       v
step-3 (agent: editor, input: "$prev.scored_items")
```

The engine will resolve `$prev.*` at runtime by extracting the named field from the previous activity's result.

## Parallel Execution

Parallel blocks compile to Temporal's fan-out pattern:

```
step: fan-out
  |
  +-- parallel[0]: task-a (agent: worker-a)
  |
  +-- parallel[1]: task-b (agent: worker-b)
  |
  v
(all complete, results collected)
  |
  v
next step
```

Each parallel sub-step runs as an independent Temporal activity. The engine waits for all to complete before advancing.

## Design Decisions

**YAML over code** -- Workflows are data, not programs. This enables non-developers to define workflows, makes them versionable and diffable, and allows tooling (linters, visualizers, policy checks) to operate on them without executing code.

**Agent abstraction** -- Steps reference agents by name, not by implementation. The engine resolves agent names to Temporal activity types at runtime, allowing the same spec to run against different agent implementations (local, remote, mock).

**$prev over explicit naming** -- The `$prev.*` syntax keeps step wiring simple for linear pipelines (the common case). For complex DAGs where steps need to reference non-adjacent outputs, explicit step-name references may be added in the future.
