package github_webhook
type WorkflowRunPayload struct {
	Action      string       `json:"action"`
	WorkflowRun WorkflowRun  `json:"workflow_run"`
	Repository  Repository   `json:"repository"`
}

type WorkflowRun struct {
	ID          int64   `json:"id"`
	WorkflowID  int64   `json:"workflow_id"`

	RunNumber   int     `json:"run_number"`
	RunAttempt  int     `json:"run_attempt"`

	Name        string  `json:"name"`
	Event       string  `json:"event"`

	HeadBranch  string  `json:"head_branch"`
	HeadSHA     string  `json:"head_sha"`

	Status      string  `json:"status"`
	Conclusion  *string `json:"conclusion"`

	JobsURL     string  `json:"jobs_url"`
	LogsURL     string  `json:"logs_url"`
	HTMLURL     string  `json:"html_url"`

	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type Repository struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Owner    Owner  `json:"owner"`
}

type Owner struct {
	Login string `json:"login"`
}