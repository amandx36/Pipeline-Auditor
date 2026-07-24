package models


import "time"

type Pipeline struct {
	ID uint `json:"id"`

	Provider string `json:"provider"`

	Repository string `json:"repository"`

	WorkflowID int64

	RunNumber  int
	RunAttempt int

	WorkflowName string

	Branch string

	CommitSHA string

	Event string

	Status string

	Conclusion string

	JobsURL string

	LogsURL string

	HTMLURL string

	CreatedAt time.Time

	UpdatedAt time.Time
}