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
	if !VerifySignature(ctx, body) {
		log.Println("Invalid in verifying the signature")
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

	// turn into the json payload 
	pipelineEvent , err := GitHub_to_PipelineEvent(payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest,gin.H{
			"Error":"failed to convert GitHub payload",
		})
		return 
	}
	fmt.Println("Pipeline Event: ", pipelineEvent)

	
	ctx.String(http.StatusOK, "OK")
}