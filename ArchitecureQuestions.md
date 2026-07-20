# Design Notes

These notes explain the reasoning behind PipelineGuard's architecture and answer common interview questions.

---

# Q1. How does PipelineGuard update metrics when a pipeline succeeds?

## Flow

```text
Developer Push
        │
        ▼
GitHub Actions
        │
        ▼
Workflow Completed
        │
        ▼
GitHub sends Webhook
        │
        ▼
PipelineGuard receives Webhook
```

GitHub sends a webhook for both successful and failed workflow runs.

Example webhook payload:

```json
{
  "workflow_run": {
    "id": 123,
    "status": "completed",
    "conclusion": "success",
    "run_started_at": "2026-07-20T10:00:00Z",
    "updated_at": "2026-07-20T10:05:32Z"
  }
}
```

The webhook contains all the metadata required to update application metrics.

PipelineGuard reads the event and updates the metrics.

Example:

```go
type MetricsService interface {
    RecordWorkflow(event WorkflowRunEvent)
}
```

Pseudo Logic

```text
if event.Conclusion == "success"

    Total Pipelines++

    Successful Pipelines++

    Update Average Build Duration

else

    Total Pipelines++

    Failed Pipelines++

    Update Failure Count
```

Example

Before

```text
Total Pipelines : 100

Success : 90

Failure : 10
```

A successful webhook arrives.

After

```text
Total Pipelines : 101

Success : 91

Failure : 10
```

A failed webhook arrives.

```text
Total Pipelines : 102

Success : 91

Failure : 11
```

### Interview Note

PipelineGuard does **not** periodically poll GitHub for metrics.

Instead, it uses GitHub Webhooks.

Each webhook contains enough metadata to update the application's internal metrics.

---

# Q2. How do workflow logs reach PipelineGuard?

One common misconception is that GitHub sends logs inside the webhook.

**This is incorrect.**

GitHub only sends metadata.

## Actual Flow

```text
Developer

↓

Git Push

↓

GitHub

↓

GitHub Actions

↓

Workflow Completed

↓

GitHub sends Webhook

↓

PipelineGuard receives metadata

↓

Extract logs_url

↓

Call GitHub REST API

↓

Download Logs (ZIP)

↓

Extract ZIP

↓

Analyze Logs
```

Webhook Example

```json
{
  "workflow_run": {
    "id": 30433642,
    "jobs_url": "...",
    "logs_url": "..."
  }
}
```

Notice that there are **URLs**, not log contents.

PipelineGuard downloads logs itself.

### Why?

Because logs may contain thousands of lines.

Sending them inside every webhook would be inefficient.

### Interview Note

GitHub Webhooks notify PipelineGuard that something happened.

GitHub APIs provide the detailed information needed for analysis.

---

# Q3. How does PipelineGuard receive pipeline events?

PipelineGuard registers webhook endpoints with GitHub.

Example

```text
POST /webhook/github
```

Repository Settings

↓

Webhooks

↓

Payload URL

↓

https://pipelineguard.example.com/webhook/github

Whenever GitHub detects a subscribed event, it sends an HTTP POST request to this endpoint.

Typical events include

- workflow_run
- push
- pull_request
- release

Example Flow

```text
Developer creates Pull Request

↓

GitHub

↓

GitHub Actions starts

↓

Workflow completes

↓

GitHub sends Webhook

↓

PipelineGuard receives event

↓

Downloads workflow information

↓

Downloads logs

↓

Analyzes failure
```

### Interview Note

PipelineGuard is event-driven.

It never continuously polls GitHub for workflow status.

GitHub pushes events to PipelineGuard using Webhooks.

---

# Q4. Every CI/CD provider has different log formats. How does PipelineGuard analyze them?

Every provider generates logs differently.

Example

## GitHub Actions

```text
Run go build

go: module github.com/gin-gonic/gin not found

Error: Process completed with exit code 1
```

---

## GitLab CI

```text
Executing "build"

fatal: module missing

Job failed
```

---

## Jenkins

```text
Building...

BUILD FAILURE

Dependency missing
```

If the analyzer tries to understand all these formats directly, the code becomes difficult to maintain.

Instead, PipelineGuard first normalizes all logs.

## Normalization Flow

```text
GitHub Logs

↓

GitHub Parser

↓

Common Format

------------------------

GitLab Logs

↓

GitLab Parser

↓

Common Format

------------------------

Jenkins Logs

↓

Jenkins Parser

↓

Common Format
```

Common internal model

```go
type LogEntry struct {
    Timestamp time.Time
    Stage     string
    Level     string
    Message   string
}
```

After normalization

GitHub

```go
{
    Stage: "Build",
    Level: "ERROR",
    Message: "module not found"
}
```

GitLab

```go
{
    Stage: "Build",
    Level: "ERROR",
    Message: "module not found"
}
```

Jenkins

```go
{
    Stage: "Build",
    Level: "ERROR",
    Message: "module not found"
}
```

The analyzer only understands **LogEntry**.

It does not care which provider generated the logs.

### Interview Note

This design makes PipelineGuard provider-independent.

Adding GitLab or Jenkins only requires implementing another parser.

The analyzer never changes.

---

# Q5. Why are there separate webhook endpoints?

Each CI/CD provider sends different payloads.

Example

```text
POST /webhook/github

POST /webhook/gitlab

POST /webhook/jenkins

POST /webhook/azure
```

Each endpoint

- verifies provider-specific signatures
- validates provider-specific requests
- parses provider-specific payloads

After parsing, every provider is converted into the same internal event model.

```go
type PipelineEvent struct {
    Provider   string
    WorkflowID string
    Repository string
    Branch     string
    Status     string
    Conclusion string
}
```

After normalization, the rest of the application becomes provider-independent.

### Interview Note

The webhook layer is provider-specific.

The business logic is provider-agnostic.

This separation makes the system easy to extend.

---

# Key Takeaways

- GitHub Webhooks send metadata, not workflow logs.
- PipelineGuard downloads logs using the GitHub REST API.
- Every provider has a dedicated webhook endpoint.
- Every provider has its own parser.
- All provider-specific data is normalized into common internal models.
- The analyzer works only with normalized models, making it independent of GitHub, GitLab, Jenkins, or Azure DevOps.
- Metrics are updated using webhook events rather than polling external APIs.