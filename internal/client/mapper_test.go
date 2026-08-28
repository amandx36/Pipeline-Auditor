package client

import (
	"testing"

	"github.com/amandx36/Pipeline-Auditor/internal/models"
)

// Unit Test — maps normalized events into protobuf fields. [Proto Mapping]
func TestToProtoPipelineEvent(t *testing.T) {
	event := models.PipelineEvent{Provider: "github", Repository: "owner/repository", PipelineID: "42", RunNumber: 7, RunAttempt: 2, PipelineName: "CI", Branch: "main", CommitSHA: "abc123", Event: "push", Status: "completed", Conclusion: "failure", LogsURL: "https://example.test/logs", HTMLURL: "https://example.test/run"}
	pipeline := toProtoPipelineEvent(event)
	if pipeline == nil || pipeline.GetProvider() != event.Provider || pipeline.GetRepository() != event.Repository ||
		pipeline.GetPipelineId() != event.PipelineID || pipeline.GetRunNumber() != uint64(event.RunNumber) ||
		pipeline.GetRunAttempt() != uint32(event.RunAttempt) || pipeline.GetPipelineName() != event.PipelineName ||
		pipeline.GetBranch() != event.Branch || pipeline.GetCommitSha() != event.CommitSHA ||
		pipeline.GetEvent() != event.Event || pipeline.GetStatus() != event.Status ||
		pipeline.GetConclusion() != event.Conclusion || pipeline.GetLogsUrl() != event.LogsURL || pipeline.GetHtmlUrl() != event.HTMLURL {
		t.Fatalf("unexpected proto pipeline: %+v", pipeline)
	}
}

// Unit Test — omits entirely empty pipeline events. [Empty Event]
func TestToProtoPipelineEventEmpty(t *testing.T) {
	if pipeline := toProtoPipelineEvent(models.PipelineEvent{}); pipeline != nil {
		t.Fatalf("toProtoPipelineEvent() = %+v, want nil", pipeline)
	}
}
