package main

import (
	"Pipeline-Auditor/internal/webhook/github"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/ping", Pong)
	router.POST("/webhook/github", handleGitHubWebHook)

	log.Println("Server started on port :8090")

	if err := router.Run(":8090"); err != nil {
		log.Fatal(err)
	}
}
func handleGitHubWebHook(ctx *gin.Context) {
	dump, err := httputil.DumpRequest(ctx.Request, true)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// fmt.Println("GIT HUB WEBHOOK")
	fmt.Println(string(dump))
	

	// now Unmarshal 
	var payload github.WorkflowRunPayload ;

	// unmarshal 

if err := json.Unmarshal(dump, &payload); err != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{
        "error": "invalid payload",
    })
    return
}
fmt.Println(payload.Action)
fmt.Println(payload.WorkflowRun.ID)
fmt.Println(payload.WorkflowRun.Status)
fmt.Println(payload.WorkflowRun.Conclusion)
fmt.Println(payload.WorkflowRun.JobsURL)
fmt.Println(payload.WorkflowRun.LogsURL)
fmt.Println(payload.Repository.FullName)



	ctx.String(http.StatusOK, "OK")
}

func Pong(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Pong")
}