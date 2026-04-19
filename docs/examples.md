---
layout: default
title: Examples
---

# Examples

Annotated examples of flowspec workflows.

## Morning Briefing

A daily cron workflow that aggregates news, scores articles, composes an editorial layout, and delivers it via Telegram and web.

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

**How it works:**

1. **ingest** -- The `ingestor` agent fetches items from Hacker News, Ars Technica RSS, and Lobsters RSS.
2. **read-and-score** -- The `reader` agent receives `$prev.items` (the ingestor's output) and scores the top 20.
3. **compose** -- The `editor` agent takes the scored items and arranges them into an editorial-hero layout.
4. **deliver** -- The `broadcaster` agent sends the composed briefing to Telegram and web channels.

## Code Review

An event-driven workflow triggered when a pull request is opened. Fetches the diff, runs an AI code review, and posts the results.

```yaml
name: code-review
trigger:
  event: gitea.pull_request.opened
steps:
  - name: fetch-diff
    agent: watcher
    config:
      action: get_pr_diff

  - name: review
    agent: reviewer
    input: "$prev.diff"
    config:
      runtime: claude-code
      model: claude-opus-4-6
      checks: [security, correctness, style]

  - name: post-review
    agent: broadcaster
    input: "$prev.review"
    config:
      channels: [gitea-comment, telegram]
```

**How it works:**

1. **fetch-diff** -- The `watcher` agent retrieves the PR diff from Gitea.
2. **review** -- The `reviewer` agent runs the diff through Claude Code with security, correctness, and style checks.
3. **post-review** -- The `broadcaster` agent posts the review as a Gitea comment and notifies via Telegram.

## Parallel Fan-Out

An example showing parallel execution of independent tasks:

```yaml
name: data-pipeline
steps:
  - name: fetch
    agent: fetcher
    config:
      url: https://api.example.com/data

  - name: process
    parallel:
      - name: analyze-sentiment
        agent: sentiment-analyzer
        input: "$prev.data"
      - name: extract-entities
        agent: entity-extractor
        input: "$prev.data"

  - name: merge-results
    agent: merger
    input: "$prev"
```

**How it works:**

1. **fetch** -- Retrieves data from an API.
2. **process** -- Fans out into two parallel sub-steps that run concurrently: sentiment analysis and entity extraction.
3. **merge-results** -- Collects outputs from both parallel sub-steps and merges them.
