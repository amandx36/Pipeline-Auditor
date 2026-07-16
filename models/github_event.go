package models

type WorkflowRunEvent struct {
    Action      string      `json:"action"`
    WorkflowRun WorkflowRun `json:"workflow_run"`
    Repository  Repository  `json:"repository"`
}

type WorkflowRun struct {
    ID         int64   `json:"id"`          // unique  identifier 
    Name       string  `json:"name"`        // workflow,s name 
    HeadBranch string  `json:"head_branch"` // kaunsi branch pe fail hua 
    HeadSHA    string  `json:"head_sha"`    // commit hash
    Status     string  `json:"status"`      // for status check 
    Conclusion *string `json:"conclusion"`  // failure/success — 
    JobsURL    string  `json:"jobs_url"`    // GitHub API se jobs fetch karne ka seedha URL
    LogsURL    string  `json:"logs_url"`    // for downlaod logs 
    HTMLURL    string  `json:"html_url"`    // for clickable link 
}

type Repository struct {
    FullName string `json:"full_name"` // "owner/repo"
}