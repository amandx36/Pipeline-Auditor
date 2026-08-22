package github_webhook 
import (
	"time"
	"Pipeline-Auditor/internal/models"
)
func GitHub_to_PipelineEvent(
	payload WorkflowRunPayload ,
)(models.PipelineEvent, error){
	createdAt , err := time.Parse(time.RFC3339,payload.WorkflowRun.CreatedAt);
	if err != nil {
		return models.PipelineEvent{}, err ;
	}
	updatedAt , err := time.Parse(time.RFC3339,payload.WorkflowRun.UpdatedAt);
	if err != nil {
		return models.PipelineEvent{}, err ;
	}
	conclusion := ""
	if payload.WorkflowRun.Conclusion != nil {
		conclusion = *payload.WorkflowRun.Conclusion
	}
	// mapping  it dude 
	return models.PipelineEvent{
		Provider:     "github",
		Repository:   payload.Repository.FullName,
		WorkflowID:   payload.WorkflowRun.WorkflowID,
		RunNumber:    payload.WorkflowRun.RunNumber,
		RunAttempt:   payload.WorkflowRun.RunAttempt,
		WorkflowName: payload.WorkflowRun.Name,
		Branch:       payload.WorkflowRun.HeadBranch,
		CommitSHA:    payload.WorkflowRun.HeadSHA,
		Event:        payload.WorkflowRun.Event,
		Status:       payload.WorkflowRun.Status,
		Conclusion:   conclusion,
		JobsURL:      payload.WorkflowRun.JobsURL,
		LogsURL:      payload.WorkflowRun.LogsURL,
		HTMLURL:      payload.WorkflowRun.HTMLURL,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}, nil

}