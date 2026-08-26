package github_webhook

import (
	"strconv"
	"time"

	"github.com/amandx36/Pipeline-Auditor/internal/models"
)

func GitHub_to_PipelineEvent(
	payload WorkflowRunPayload,
) (models.PipelineEvent, error) {

	createdAt, err := time.Parse(
		time.RFC3339,
		payload.WorkflowRun.CreatedAt,
	)
	if err != nil {
		return models.PipelineEvent{}, err
	}

	updatedAt, err := time.Parse(
		time.RFC3339,
		payload.WorkflowRun.UpdatedAt,
	)
	if err != nil {
		return models.PipelineEvent{}, err
	}

	conclusion := ""

	if payload.WorkflowRun.Conclusion != nil {
		conclusion = *payload.WorkflowRun.Conclusion
	}

	return models.PipelineEvent{
		Provider:     "github",
		Repository:   payload.Repository.FullName,
		PipelineID:   strconv.FormatInt(payload.WorkflowRun.WorkflowID, 10),
		RunNumber:    payload.WorkflowRun.RunNumber,
		RunAttempt:   payload.WorkflowRun.RunAttempt,
		PipelineName: payload.WorkflowRun.Name,
		Branch:       payload.WorkflowRun.HeadBranch,
		CommitSHA:    payload.WorkflowRun.HeadSHA,
		Event:        payload.WorkflowRun.Event,
		Status:       payload.WorkflowRun.Status,
		Conclusion:   conclusion,
		LogsURL:      payload.WorkflowRun.LogsURL,
		HTMLURL:      payload.WorkflowRun.HTMLURL,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}, nil
}
