package github_webhook

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"fmt"

	"github.com/gin-gonic/gin"
)


func HandleGitHubWebHook(ctx *gin.Context) {
	dump, err := httputil.DumpRequest(ctx.Request, true)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// fmt.Println("GIT HUB WEBHOOK")
	fmt.Println(string(dump))
	

	// now Unmarshal 
	var payload WorkflowRunPayload ;

	// unmarshal 

if err := json.Unmarshal(dump, &payload); err != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{
        "error": "invalid payload",
    })
    return
}

fmt.Println("\n\n\n\n\n\n\n\n")
fmt.Println(payload.Action)
fmt.Println(payload.WorkflowRun.ID)
fmt.Println(payload.WorkflowRun.Status)
fmt.Println(payload.WorkflowRun.Conclusion)
fmt.Println(payload.WorkflowRun.JobsURL)
fmt.Println(payload.WorkflowRun.LogsURL)
fmt.Println(payload.Repository.FullName)



	ctx.String(http.StatusOK, "OK")
}
