package github_webhook

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleGitHubWebHook(ctx *gin.Context) {

	// Read request body
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to read request body",
		})
		return
	}

	// Print raw JSON (for debugging)
	fmt.Println(string(body))

	// Verify GitHub Signature
	if VerifySignature(ctx, body) {
		fmt.Println("\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\nPassedddddd \n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n ")
		
	}

	// Unmarshal JSON
	var payload WorkflowRunPayload

	if err := json.Unmarshal(body, &payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid payload",
		})
		return
	}

	fmt.Println("GitHub Workflow ")
	fmt.Println("Action      :", payload.Action)
	fmt.Println("Run ID      :", payload.WorkflowRun.ID)
	fmt.Println("Status      :", payload.WorkflowRun.Status)
	fmt.Println("Conclusion  :", payload.WorkflowRun.Conclusion)
	fmt.Println("Jobs URL    :", payload.WorkflowRun.JobsURL)
	fmt.Println("Logs URL    :", payload.WorkflowRun.LogsURL)
	fmt.Println("Repository  :", payload.Repository.FullName)

	ctx.String(http.StatusOK, "OK")
}