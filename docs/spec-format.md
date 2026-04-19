---
layout: default
title: Spec Format Reference
---

# Spec Format Reference

flowspec workflows are defined in YAML. This page documents every field.

## Top-Level: Workflow

| Field      | Type       | Required | Description                                          |
|------------|------------|----------|------------------------------------------------------|
| `name`     | `string`   | Yes      | Unique identifier for the workflow.                  |
| `trigger`  | `Trigger`  | No       | When the workflow should be executed.                |
| `steps`    | `[]Step`   | Yes      | Ordered list of steps to execute.                    |
| `on_error` | `string`   | No       | Default error strategy: `retry`, `skip`, or `abort`. |

```yaml
name: my-workflow
trigger:
  cron: "0 8 * * *"
on_error: abort
steps:
  - name: step-one
    agent: worker
```

## Trigger

Defines when a workflow should run. At most one trigger type should be set.

| Field     | Type     | Description                                      |
|-----------|----------|--------------------------------------------------|
| `cron`    | `string` | Cron expression (e.g., `"0 8 * * *"`).           |
| `event`   | `string` | Event name (e.g., `"gitea.pull_request.opened"`). |
| `webhook` | `string` | Webhook endpoint path.                           |

```yaml
trigger:
  cron: "0 8 * * *"
```

```yaml
trigger:
  event: gitea.pull_request.opened
```

## Step

A single unit of work in a workflow.

| Field       | Type             | Required | Description                                                  |
|-------------|------------------|----------|--------------------------------------------------------------|
| `name`      | `string`         | Yes      | Unique name within the workflow.                             |
| `agent`     | `string`         | *        | Agent that executes this step. Required unless `parallel` is set. |
| `input`     | `string`         | No       | Input expression, typically `$prev.*` referencing prior output. |
| `config`    | `map[string]any` | No       | Arbitrary key-value configuration passed to the agent.       |
| `timeout`   | `string`         | No       | Maximum duration (e.g., `"5m"`, `"1h"`).                     |
| `on_error`  | `string`         | No       | Error strategy: `retry`, `skip`, or `abort`.                 |
| `parallel`  | `[]Step`         | No       | Sub-steps to execute in parallel (fan-out).                  |
| `condition` | `string`         | No       | Boolean expression for conditional execution (e.g., `"$prev.urgent == true"`). |

### Agent Steps

A step with an `agent` field runs a single activity:

```yaml
- name: fetch-data
  agent: fetcher
  config:
    url: https://example.com/api
  timeout: 30s
  on_error: retry
```

### Parallel Steps

A step with a `parallel` block fans out into concurrent sub-steps:

```yaml
- name: fan-out
  parallel:
    - name: task-a
      agent: worker-a
    - name: task-b
      agent: worker-b
```

### Input References

Steps reference prior outputs with `$prev.*` syntax:

```yaml
- name: process
  agent: processor
  input: "$prev.items"
```

### Conditional Execution

Steps can be conditionally executed:

```yaml
- name: alert
  agent: notifier
  input: "$prev.result"
  condition: "$prev.urgent == true"
```

## Error Handling

The `on_error` field accepts three values:

| Value   | Behavior                              |
|---------|---------------------------------------|
| `retry` | Retry the step with backoff.          |
| `skip`  | Skip the step and continue.          |
| `abort` | Stop the entire workflow immediately. |

`on_error` can be set at the workflow level (default for all steps) or per-step (overrides the workflow default).
