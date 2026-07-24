package main

import (

	"log"
	"net/http"
	"Pipeline-Auditor/internal/webhook/github"
	

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/ping", Pong)
	router.POST("/webhook/github", github_webhook.HandleGitHubWebHook)

	log.Println("Server started on port :8090")

	if err := router.Run(":8090"); err != nil {
		log.Fatal(err)
	}
}

func Pong(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Pong")
}