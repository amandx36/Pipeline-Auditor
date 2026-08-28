package github_webhook

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func signedHeader(secret string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write(body)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

func verifyRequest(t *testing.T, body []byte, signature string) bool {
	t.Helper()
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	request := httptest.NewRequest("POST", "/webhook/github", bytes.NewReader(body))
	request.Header.Set("X-Hub-Signature-256", signature)
	ctx.Request = request
	return VerifySignature(ctx, body)
}

// Unit Test — validates GitHub HMAC verification scenarios. [Signature Validation]
func TestVerifySignature(t *testing.T) {
	gin.SetMode(gin.TestMode)
	const secret = "test-webhook-secret"
	body := []byte(`{"workflow_run":{"status":"completed"}}`)
	validSignature := signedHeader(secret, body)

	tests := []struct {
		name      string
		secret    string
		body      []byte
		signature string
		want      bool
	}{
		{name: "[Valid Signature] correct secret and body", secret: secret, body: body, signature: validSignature, want: true},
		{name: "[Wrong Secret] verifier uses another secret", secret: "wrong-secret", body: body, signature: validSignature, want: false},
		{name: "[Tampered Body] signed body is modified", secret: secret, body: []byte(`{"workflow_run":{"status":"changed"}}`), signature: validSignature, want: false},
		{name: "[Modified Signature] digest is altered", secret: secret, body: body, signature: validSignature[:len(validSignature)-1] + "0", want: false},
		{name: "[Empty Signature] header is absent", secret: secret, body: body, signature: "", want: false},
		{name: "[Malformed Signature] header is not HMAC syntax", secret: secret, body: body, signature: "not-a-signature", want: false},
		{name: "[Different Body] signature is replayed", secret: secret, body: []byte(`{"workflow_run":{}}`), signature: validSignature, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("my_git_hubSecret", tt.secret)
			if got := verifyRequest(t, tt.body, tt.signature); got != tt.want {
				t.Fatalf("VerifySignature() = %t, want %t", got, tt.want)
			}
		})
	}
}

// Fuzz Test — checks verifier robustness against arbitrary inputs. [Fuzz Verification]
func FuzzVerifySignature(f *testing.F) {
	const secret = "fuzz-webhook-secret"
	f.Add([]byte(`{"workflow_run":{}}`), "sha256=invalid")
	f.Add([]byte("payload"), signedHeader(secret, []byte("payload")))

	f.Fuzz(func(t *testing.T, body []byte, signature string) {
		gin.SetMode(gin.TestMode)
		t.Setenv("my_git_hubSecret", secret)
		got := verifyRequest(t, body, signature)
		if got && signature != signedHeader(secret, body) {
			t.Fatal("VerifySignature accepted a signature not produced for the supplied body")
		}
	})
}

// Benchmark Test — measures HMAC verification throughput. [Verification Speed]
func BenchmarkVerifySignature(b *testing.B) {
	gin.SetMode(gin.TestMode)
	const secret = "benchmark-webhook-secret"
	body := []byte(`{"workflow_run":{"status":"completed","conclusion":"failure"}}`)
	signature := signedHeader(secret, body)
	b.Setenv("my_git_hubSecret", secret)
	b.ReportAllocs()

	for b.Loop() {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		request := httptest.NewRequest("POST", "/webhook/github", bytes.NewReader(body))
		request.Header.Set("X-Hub-Signature-256", signature)
		ctx.Request = request
		if !VerifySignature(ctx, body) {
			b.Fatal("valid signature rejected")
		}
	}
}
