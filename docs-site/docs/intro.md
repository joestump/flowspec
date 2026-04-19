---
sidebar_position: 1
title: Introduction
---

# flowspec

A YAML DSL for defining and executing Temporal workflows. Declarative workflow definitions that compile to Temporal workflows -- define pipelines as YAML, execute them with a generic Go interpreter.

Designed for AI agent orchestration but general-purpose.

## Quick Start

Install the library:

```bash
go get github.com/joestump/flowspec
```

Parse and validate a workflow:

```go
wf, err := flowspec.ParseFile("workflow.yaml")
if err != nil {
    log.Fatal(err)
}

if err := flowspec.Validate(wf); err != nil {
    log.Fatal(err)
}

fmt.Printf("Workflow: %s (%d steps)\n", wf.Name, len(wf.Steps))
```

## Documentation

- [Spec Format Reference](spec-format.md) -- full YAML spec reference with all fields documented
- [Examples](examples.md) -- annotated workflow examples
- [Architecture](architecture.md) -- how flowspec compiles to Temporal workflows

## Links

- [GitHub Repository](https://github.com/joestump/flowspec)
- [pkg.go.dev Reference](https://pkg.go.dev/github.com/joestump/flowspec)
