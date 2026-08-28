package github_webhook

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func failurePayload() WorkflowRunPayload {
	conclusion := "failure"
	return WorkflowRunPayload{
		Action: "completed",
		WorkflowRun: WorkflowRun{
			ID: 30433642, WorkflowID: 159038, RunNumber: 42, RunAttempt: 1,
			Name: "CI", Event: "push", HeadBranch: "main", HeadSHA: "acb5820",
			Status: "completed", Conclusion: &conclusion, LogsURL: "https://example.test/logs",
			HTMLURL: "https://example.test/run", CreatedAt: "2026-07-16T10:30:00Z", UpdatedAt: "2026-07-16T10:32:15Z",
		},
		Repository: Repository{FullName: "amandx36/pipelineguard", Owner: Owner{Login: "amandx36"}},
	}
}

// Unit Test — verifies GitHub field normalization. [Field Mapping]
func TestGitHubToPipelineEvent(t *testing.T) {
	event, err := GitHub_to_PipelineEvent(failurePayload())
	if err != nil {
		t.Fatalf("GitHub_to_PipelineEvent() error = %v", err)
	}

	createdAt := time.Date(2026, time.July, 16, 10, 30, 0, 0, time.UTC)
	updatedAt := time.Date(2026, time.July, 16, 10, 32, 15, 0, time.UTC)
	if event.Provider != "github" || event.Repository != "amandx36/pipelineguard" ||
		event.PipelineID != "159038" || event.RunNumber != 42 || event.RunAttempt != 1 ||
		event.PipelineName != "CI" || event.Branch != "main" || event.CommitSHA != "acb5820" ||
		event.Event != "push" || event.Status != "completed" || event.Conclusion != "failure" ||
		event.LogsURL != "https://example.test/logs" || event.HTMLURL != "https://example.test/run" ||
		!event.CreatedAt.Equal(createdAt) || !event.UpdatedAt.Equal(updatedAt) {
		t.Fatalf("unexpected normalized event: %+v", event)
	}
}

// Unit Test — normalizes the supplied GitHub failure fixture. [Failure Fixture]
func TestGitHubToPipelineEventFailureFixture(t *testing.T) {
	body, err := os.ReadFile("../../../assets/workflow_run_failure_sample.json")
	if err != nil {
		t.Fatal(err)
	}
	var payload WorkflowRunPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatal(err)
	}

	event, err := GitHub_to_PipelineEvent(payload)
	if err != nil {
		t.Fatal(err)
	}
	if event.Repository != "amandx36/pipelineguard" || event.PipelineID != "159038" || event.Status != "completed" || event.Conclusion != "failure" || event.RunNumber != 42 {
		t.Fatalf("unexpected fixture event: %+v", event)
	}
}

// Unit Test — rejects invalid workflow timestamps. [Malformed Timestamp]
func TestGitHubToPipelineEventInvalidTimestamp(t *testing.T) {
	payload := failurePayload()
	payload.WorkflowRun.CreatedAt = "not-a-timestamp"
	if _, err := GitHub_to_PipelineEvent(payload); err == nil {
		t.Fatal("GitHub_to_PipelineEvent() error = nil, want timestamp parsing error")
	}
}

// Benchmark Test — measures payload normalization performance. [Mapping Speed]
func BenchmarkGitHubToPipelineEvent(b *testing.B) {
	payload := failurePayload()
	b.ReportAllocs()
	for b.Loop() {
		if _, err := GitHub_to_PipelineEvent(payload); err != nil {
			b.Fatal(err)
		}
	}
}
