package github_webhook

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/amandx36/Pipeline-Auditor/internal/client"
	// "github.com/amandx36/Pipeline-Auditor/internal/storage/postgres"
	"github.com/gin-gonic/gin"
)

func HandleGitHubWebHook(ctx *gin.Context, db *sql.DB) {

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
	// if !VerifySignature(ctx, body) {
	// 	log.Println("Invalid in verifying the signature")
	// 	ctx.JSON(http.StatusUnauthorized, gin.H{
	// 		"error": "invalid signature",
	// 	})
	// 	return
	// }
	// achieve the idempotency
	deliveryId := ctx.GetHeader("X-GitHub-Delivery")
	if deliveryId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "missing X-GitHub-Delivery",
		})
		return
	}
	// deliveryStore := postgres.NewDeliveryStore(db)
	// // isNew, err := deliveryStore.TryCreate(ctx.Request.Context(), deliveryId)
	// // if err != nil {
	// // 	log.Printf("idempotency check failed for delivery %q: %v", deliveryId, err)
	// // 	ctx.JSON(http.StatusInternalServerError, gin.H{
	// // 		"status": "failed idempotency check",
	// // 	})
	// // 	return
	// // }
	// if !isNew {
	// 	log.Printf("Webhook ignored: duplicate delivery_id=%q", deliveryId)
	// 	ctx.JSON(http.StatusAccepted, gin.H{
	// 		"status": "duplicate",
	// 	})
	// 	return
	// }

	// Unmarshal JSON
	var payload WorkflowRunPayload

	if err := json.Unmarshal(body, &payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid payload",
		})
		return
	}

	// completed work flow  send further
	// if payload.WorkflowRun.Status != "completed" {
	// 	log.Printf("[PIPELINE-AUDITOR] Webhook ignored: workflow status=%q (requires completed)", payload.WorkflowRun.Status)
	// 	ctx.JSON(http.StatusAccepted, gin.H{
	// 		"status": "ignored",
	// 		"reason": "workflow is not completed",
	// 	})
	// 	return
	// }

	// only send the  failure json other wise ignore it
	// if payload.WorkflowRun.Conclusion == nil ||
	// 	*payload.WorkflowRun.Conclusion != "failure" {
	// 	conclusion := ""
	// 	if payload.WorkflowRun.Conclusion != nil {
	// 		conclusion = *payload.WorkflowRun.Conclusion
	// 	}
	// 	log.Printf("[PIPELINE-AUDITOR] Webhook ignored: workflow conclusion=%q (requires failure)", conclusion)

	// 	ctx.JSON(http.StatusAccepted, gin.H{
	// 		"status": "ignored",
	// 		"reason": "workflow did not fail",
	// 	})
	// 	return
	// }

	//  convert the
	pipelineEvent, err := GitHub_to_PipelineEvent(payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to convert GitHub payload",
		})
		return
	}

	fmt.Println("Pipeline Event:", pipelineEvent)

	log.Println("[PIPELINE-AUDITOR] Webhook reached CollectLogs; sending PipelineEvent to CI-LogCollector")
	response, err := client.CollectLogWrapper(pipelineEvent)
	if err != nil {
		log.Printf("[PIPELINE-AUDITOR] gRPC CollectLogs failed: %v", err)
		ctx.JSON(http.StatusFailedDependency, gin.H{
			"error": "log collection request failed",
		})
		return
	}
	if response != nil {
		log.Printf(
			"[PIPELINE-AUDITOR] CI-LogCollector response: accepted=%t collection_id=%s message=%s",
			response.GetAccepted(),
			response.GetCollectionId(),
			response.GetMessage(),
		)
	}

	ctx.String(http.StatusOK, "OK")
}
