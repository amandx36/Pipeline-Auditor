github gives this responses in the webhook/github

https://docs.github.com/en/rest/actions/workflow-runs?apiVersion=2026-03-10
    

https://docs.github.com/en/webhooks/webhook-events-and-payloads#workflow_run

like json 

// internal/models/github_event.go
package models

type WorkflowRunEvent struct {
    Action     string      `json:"action"`
    Workflow   *Workflow   `json:"workflow"`
    WorkflowRun WorkflowRun `json:"workflow_run"`
    Repository Repository  `json:"repository"`
    Sender     Sender      `json:"sender"`
}

type Workflow struct {
    ID   int64  `json:"id"`
    Name string `json:"name"`
    Path string `json:"path"`
}

type WorkflowRun struct {
    ID          int64   `json:"id"`
    Name        string  `json:"name"`
    HeadBranch  string  `json:"head_branch"`
    HeadSHA     string  `json:"head_sha"`
    Status      string  `json:"status"`      // requested, in_progress, completed, queued, pending, waiting
    Conclusion  *string `json:"conclusion"`   // success, failure, cancelled, timed_out, etc. (null if not done)
    WorkflowID  int64   `json:"workflow_id"`
    RunNumber   int     `json:"run_number"`
    RunAttempt  int     `json:"run_attempt"`
    Event       string  `json:"event"`
    URL         string  `json:"url"`
    HTMLURL     string  `json:"html_url"`
    JobsURL     string  `json:"jobs_url"`
    LogsURL     string  `json:"logs_url"`
    CreatedAt   string  `json:"created_at"`
    UpdatedAt   string  `json:"updated_at"`
    RunStartedAt string `json:"run_started_at"`
}

type Repository struct {
    ID       int64  `json:"id"`
    Name     string `json:"name"`
    FullName string `json:"full_name"`
}

type Sender struct {
    Login string `json:"login"`
    ID    int64  `json:"id"`
}