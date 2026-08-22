package github_webhook

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleGitHubWebHook(ctx *gin.Context) {

	// read the  request body
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to read request body",
		})
		return
	}

	// incoming raw json 
	fmt.Println(string(body))

	// github signature verification
	if !VerifySignature(ctx, body) {
		log.Println("Invalid in verifying the signature")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid signature",
		})
		return
	}

	// Unmarshal JSON
	var payload WorkflowRunPayload

	if err := json.Unmarshal(body, &payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid payload",
		})
		return
	}

	// completed work flow  send further 
	if payload.WorkflowRun.Status != "completed" {
		ctx.JSON(http.StatusAccepted, gin.H{
			"status": "ignored",
			"reason": "workflow is not completed",
		})
		return
	}

	// only send the  failure json other wise ignore it 
	if payload.WorkflowRun.Conclusion == nil ||
		*payload.WorkflowRun.Conclusion != "failure" {

		ctx.JSON(http.StatusAccepted, gin.H{
			"status": "ignored",
			"reason": "workflow did not fail",
		})
		return
	}

	//  convert the 
	pipelineEvent, err := GitHub_to_PipelineEvent(payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to convert GitHub payload",
		})
		return
	}

	fmt.Println("Pipeline Event:", pipelineEvent)

	ctx.String(http.StatusOK, "OK")
}