package models

import "time"

type PipelineEvent struct {
	ID           uint      `json:"id"`
	Provider     string    `json:"provider"`
	Repository   string    `json:"repository"`
	WorkflowID   int64     `json:"workflow_id"`
	RunNumber    int       `json:"run_number"`
	RunAttempt   int       `json:"run_attempt"`
	WorkflowName string    `json:"workflow_name"`
	Branch       string    `json:"branch"`
	CommitSHA    string    `json:"commit_sha"`
	Event        string    `json:"event"`
	Status       string    `json:"status"`
	Conclusion   string    `json:"conclusion"`
	JobsURL      string    `json:"jobs_url"`
	LogsURL      string    `json:"logs_url"`
	HTMLURL      string    `json:"html_url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}