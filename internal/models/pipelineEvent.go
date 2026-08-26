package models

import "time"

type PipelineEvent struct {
	ID           uint      `json:"id"`
	Provider     string    `json:"provider"`
	Repository   string    `json:"repository"`
	PipelineID   string    `json:"pipeline_id"`
	RunNumber    int       `json:"run_number"`
	RunAttempt   int       `json:"run_attempt"`
	PipelineName string    `json:"pipeline_name"`
	Branch       string    `json:"branch"`
	CommitSHA    string    `json:"commit_sha"`
	Event        string    `json:"event"`
	Status       string    `json:"status"`
	Conclusion   string    `json:"conclusion"`
	LogsURL      string    `json:"logs_url"`
	HTMLURL      string    `json:"html_url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
