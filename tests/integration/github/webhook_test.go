package github_integration_test

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"

	github_webhook "github.com/amandx36/Pipeline-Auditor/internal/webhook/github"
	"github.com/gin-gonic/gin"
)

func integrationSignature(secret string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write(body)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

// Integration Test — rejects an invalid signed HTTP webhook. [HTTP 401]
func TestGitHubWebhookRejectsInvalidSignature(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv("my_git_hubSecret", "integration-secret")
	router := gin.New()
	router.POST("/webhook/github", func(ctx *gin.Context) { github_webhook.HandleGitHubWebHook(ctx, nil) })
	body := []byte(`{"workflow_run":{}}`)
	request := httptest.NewRequest(http.MethodPost, "/webhook/github", bytes.NewReader(body))
	request.Header.Set("X-Hub-Signature-256", integrationSignature("wrong-secret", body))
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized || !bytes.Contains(recorder.Body.Bytes(), []byte("invalid signature")) {
		t.Fatalf("response = %d %s, want 401 invalid signature", recorder.Code, recorder.Body.String())
	}
}

// Integration Test — validates delivery identity before storage access. [HTTP 400]
func TestGitHubWebhookRejectsMissingDeliveryID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	const secret = "integration-secret"
	t.Setenv("my_git_hubSecret", secret)
	router := gin.New()
	router.POST("/webhook/github", func(ctx *gin.Context) { github_webhook.HandleGitHubWebHook(ctx, nil) })
	body := []byte(`{"workflow_run":{}}`)
	request := httptest.NewRequest(http.MethodPost, "/webhook/github", bytes.NewReader(body))
	request.Header.Set("X-Hub-Signature-256", integrationSignature(secret, body))
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest || !bytes.Contains(recorder.Body.Bytes(), []byte("missing X-GitHub-Delivery")) {
		t.Fatalf("response = %d %s, want 400 missing delivery", recorder.Code, recorder.Body.String())
	}
}
