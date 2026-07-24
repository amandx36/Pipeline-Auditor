package main

import (

	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

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
func handleGitHubWebHook(c *gin.Context) {
	dump, err := httputil.DumpRequest(c.Request, true)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	fmt.Println("========== GITHUB WEBHOOK ==========")
	fmt.Println(string(dump))
	fmt.Println("====================================")

	c.String(http.StatusOK, "OK")
}

func Pong(c *gin.Context) {
	c.String(http.StatusOK, "Pong")
}