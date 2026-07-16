package models

type WorkflowRunEvent struct {
    Action      string      `json:"action"`
    Workflow    *Workflow   `json:"workflow"`
    WorkflowRun WorkflowRun `json:"workflow_run"`
    Repository  Repository  `json:"repository"`
    Sender      Sender      `json:"sender"`
}

type WorkflowRun struct {
    ID           int64   `json:"id"`
    Name         string  `json:"name"`
    NodeID       string  `json:"node_id"`
    HeadBranch   string  `json:"head_branch"`
    HeadSHA      string  `json:"head_sha"`
    Path         string  `json:"path"`
    RunNumber    int     `json:"run_number"`
    Event        string  `json:"event"`
    DisplayTitle string  `json:"display_title"`
    Status       string  `json:"status"`       // queued, in_progress, completed
    Conclusion   *string `json:"conclusion"`    // null until completed; then success/failure/etc
    WorkflowID   int64   `json:"workflow_id"`
    URL          string  `json:"url"`
    HTMLURL      string  `json:"html_url"`
    RunAttempt   int     `json:"run_attempt"`
    RunStartedAt string  `json:"run_started_at"`
    CreatedAt    string  `json:"created_at"`
    UpdatedAt    string  `json:"updated_at"`
    JobsURL      string  `json:"jobs_url"`
    LogsURL      string  `json:"logs_url"`
    CheckSuiteURL string `json:"check_suite_url"`
    WorkflowURL  string  `json:"workflow_url"`
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

type Workflow struct {
    ID   int64  `json:"id"`
    Name string `json:"name"`
    Path string `json:"path"`
}