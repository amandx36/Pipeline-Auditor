package client

import (
	collectorpb "github.com/amandx36/Pipeline-Auditor/gen/collector"
	"github.com/amandx36/Pipeline-Auditor/internal/models"
)

// converts the pipelineEvent to collectorpb.pipelineEvent 

func toProtoPipelineEvent(
	event models.PipelineEvent,
) *collectorpb.PipelineEvent {

	if event.Provider == "" &&
		event.Repository == "" &&
		event.PipelineID == "" {

		return nil
	}

	// Create the protobuf object.
	return &collectorpb.PipelineEvent{


		Provider: event.Provider,

		Repository: event.Repository,


		PipelineId: event.PipelineID,

		RunNumber: uint64(event.RunNumber),

		RunAttempt: uint32(event.RunAttempt),

		PipelineName: event.PipelineName,

		Branch: event.Branch,

		CommitSha: event.CommitSHA,

		Event: event.Event,

		Status: event.Status,

		Conclusion: event.Conclusion,

		LogsUrl: event.LogsURL,

		HtmlUrl: event.HTMLURL,
	}
}