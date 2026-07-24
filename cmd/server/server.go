package main

import (
	"log"
	"net/http"

	github_webhook "Pipeline-Auditor/internal/webhook/github"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	router := gin.Default()
	err := godotenv.Load()
	if err != nil {
		log.Println("Error while loading the env",err)
	}

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