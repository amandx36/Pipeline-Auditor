package github_webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	collectorpb "github.com/amandx36/Pipeline-Auditor/gen/collector"
	"github.com/amandx36/Pipeline-Auditor/internal/models"
	"github.com/gin-gonic/gin"
)

type fakeDeliveryStore struct {
	isNew bool
	err   error
	ids   []string
}

func (s *fakeDeliveryStore) TryCreate(_ context.Context, deliveryID string) (bool, error) {
	s.ids = append(s.ids, deliveryID)
	return s.isNew, s.err
}

func webhookRequest(t *testing.T, payload any, secret, deliveryID string) (*gin.Context, *httptest.ResponseRecorder) {
	t.Helper()
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	request := httptest.NewRequest(http.MethodPost, "/webhook/github", bytes.NewReader(body))
	request.Header.Set("X-Hub-Signature-256", signedHeader(secret, body))
	request.Header.Set("X-GitHub-Delivery", deliveryID)
	ctx.Request = request
	return ctx, recorder
}

// Unit Test — validates webhook filtering and dispatch behavior. [Event Filtering]
func TestHandleGitHubWebhookFiltering(t *testing.T) {
	gin.SetMode(gin.TestMode)
	const secret = "handler-webhook-secret"
	tests := []struct {
		name          string
		mutate        func(*WorkflowRunPayload)
		delivery      fakeDeliveryStore
		wantStatus    int
		wantDispatch  bool
		wantBodyField string
	}{
		{name: "[Failure Event] completed failure dispatches", delivery: fakeDeliveryStore{isNew: true}, wantStatus: http.StatusOK, wantDispatch: true},
		{name: "[Success Event] completed success ignores", mutate: func(p *WorkflowRunPayload) { value := "success"; p.WorkflowRun.Conclusion = &value }, delivery: fakeDeliveryStore{isNew: true}, wantStatus: http.StatusAccepted, wantBodyField: "workflow did not fail"},
		{name: "[Cancelled Event] completed cancelled ignores", mutate: func(p *WorkflowRunPayload) { value := "cancelled"; p.WorkflowRun.Conclusion = &value }, delivery: fakeDeliveryStore{isNew: true}, wantStatus: http.StatusAccepted, wantBodyField: "workflow did not fail"},
		{name: "[Timeout Event] completed timeout ignores", mutate: func(p *WorkflowRunPayload) { value := "timed_out"; p.WorkflowRun.Conclusion = &value }, delivery: fakeDeliveryStore{isNew: true}, wantStatus: http.StatusAccepted, wantBodyField: "workflow did not fail"},
		{name: "[Incomplete Event] running failure ignores", mutate: func(p *WorkflowRunPayload) { p.WorkflowRun.Status = "in_progress" }, delivery: fakeDeliveryStore{isNew: true}, wantStatus: http.StatusAccepted, wantBodyField: "workflow is not completed"},
		{name: "[Duplicate Delivery] stored ID ignores", delivery: fakeDeliveryStore{isNew: false}, wantStatus: http.StatusAccepted, wantBodyField: "duplicate"},
		{name: "[Idempotency Failure] storage failure returns error", delivery: fakeDeliveryStore{err: errors.New("database unavailable")}, wantStatus: http.StatusInternalServerError, wantBodyField: "failed idempotency check"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("my_git_hubSecret", secret)
			payload := failurePayload()
			if tt.mutate != nil {
				tt.mutate(&payload)
			}
			store := tt.delivery
			var dispatched []models.PipelineEvent
			dispatch := func(event models.PipelineEvent) (*collectorpb.CollectLogsResponse, error) {
				dispatched = append(dispatched, event)
				return &collectorpb.CollectLogsResponse{Accepted: true, CollectionId: "collection-1"}, nil
			}
			ctx, recorder := webhookRequest(t, payload, secret, "delivery-1")
			handleGitHubWebhook(ctx, &store, dispatch)

			if recorder.Code != tt.wantStatus {
				t.Fatalf("HTTP status = %d, want %d; body=%s", recorder.Code, tt.wantStatus, recorder.Body.String())
			}
			if got := len(dispatched) == 1; got != tt.wantDispatch {
				t.Fatalf("dispatch happened = %t, want %t", got, tt.wantDispatch)
			}
			if tt.wantBodyField != "" && !bytes.Contains(recorder.Body.Bytes(), []byte(tt.wantBodyField)) {
				t.Fatalf("response body %q does not contain %q", recorder.Body.String(), tt.wantBodyField)
			}
		})
	}
}

// Unit Test — rejects malformed post-verification payload. [Malformed Payload]
func TestHandleGitHubWebhookMalformedPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)
	const secret = "handler-webhook-secret"
	t.Setenv("my_git_hubSecret", secret)
	body := []byte(`{"workflow_run":`)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	request := httptest.NewRequest(http.MethodPost, "/webhook/github", bytes.NewReader(body))
	request.Header.Set("X-Hub-Signature-256", signedHeader(secret, body))
	request.Header.Set("X-GitHub-Delivery", "delivery-malformed")
	ctx.Request = request
	store := &fakeDeliveryStore{isNew: true}
	handleGitHubWebhook(ctx, store, func(models.PipelineEvent) (*collectorpb.CollectLogsResponse, error) { return nil, nil })

	if recorder.Code != http.StatusBadRequest || !bytes.Contains(recorder.Body.Bytes(), []byte("invalid payload")) {
		t.Fatalf("response = %d %s, want 400 invalid payload", recorder.Code, recorder.Body.String())
	}
}

// Unit Test — reports downstream collector failure. [Dispatch Failure]
func TestHandleGitHubWebhookDispatchFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	const secret = "handler-webhook-secret"
	t.Setenv("my_git_hubSecret", secret)
	ctx, recorder := webhookRequest(t, failurePayload(), secret, "delivery-dispatch-error")
	handleGitHubWebhook(ctx, &fakeDeliveryStore{isNew: true}, func(models.PipelineEvent) (*collectorpb.CollectLogsResponse, error) {
		return nil, errors.New("collector unavailable")
	})

	if recorder.Code != http.StatusFailedDependency || !bytes.Contains(recorder.Body.Bytes(), []byte("log collection request failed")) {
		t.Fatalf("response = %d %s, want 424 dispatch failure", recorder.Code, recorder.Body.String())
	}
}

// Unit Test — rejects requests without delivery identity. [Missing Delivery]
func TestHandleGitHubWebhookMissingDeliveryID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	const secret = "handler-webhook-secret"
	t.Setenv("my_git_hubSecret", secret)
	ctx, recorder := webhookRequest(t, failurePayload(), secret, "")
	handleGitHubWebhook(ctx, &fakeDeliveryStore{isNew: true}, func(models.PipelineEvent) (*collectorpb.CollectLogsResponse, error) { return nil, nil })

	if recorder.Code != http.StatusBadRequest || !bytes.Contains(recorder.Body.Bytes(), []byte("missing X-GitHub-Delivery")) {
		t.Fatalf("response = %d %s, want 400 missing delivery", recorder.Code, recorder.Body.String())
	}
}

// Unit Test — documents current event-header handling. [Event Header]
func TestHandleGitHubWebhookEventHeaderBehavior(t *testing.T) {
	gin.SetMode(gin.TestMode)
	const secret = "handler-webhook-secret"
	t.Setenv("my_git_hubSecret", secret)
	tests := []struct {
		name   string
		header string
	}{
		{name: "[Missing Header] absent event header dispatches", header: ""},
		{name: "[Unknown Header] arbitrary event header dispatches", header: "workflow_dispatch"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, recorder := webhookRequest(t, failurePayload(), secret, "delivery-event-header")
			if tt.header != "" {
				ctx.Request.Header.Set("X-GitHub-Event", tt.header)
			}
			dispatched := false
			handleGitHubWebhook(ctx, &fakeDeliveryStore{isNew: true}, func(models.PipelineEvent) (*collectorpb.CollectLogsResponse, error) {
				dispatched = true
				return &collectorpb.CollectLogsResponse{Accepted: true}, nil
			})

			if recorder.Code != http.StatusOK || !dispatched {
				t.Fatalf("response = %d, dispatched = %t; current handler does not validate X-GitHub-Event", recorder.Code, dispatched)
			}
		})
	}
}
