package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	github_webhook "github.com/amandx36/Pipeline-Auditor/internal/webhook/github"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {

	router := gin.Default()
	err := godotenv.Load()

	if err != nil {
		log.Println("Error while loading the env", err)
	}

	router.GET("/ping", Pong)
	// connecting to he database
	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}
	// make the anonymis function
	router.POST("/webhook/github", func(ctx *gin.Context) {
		github_webhook.HandleGitHubWebHook(ctx, db)
	})

	log.Println("Server started on port :8090")

	if err := router.Run(":8090"); err != nil {
		log.Fatal(err)
	}
}

func Pong(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Pong  hello testing 03 ")
}



