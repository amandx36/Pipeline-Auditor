package github_webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
)

func VerifySignature(ctx *gin.Context) bool {

	
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Error while reading the body:", err)
		return false
	}

	fmt.Println(string(body))

	// reading github signature 
	GitHUbSignature := ctx.GetHeader("X-Hub-Signature-256")
	fmt.Println("GitHub Signature :", GitHUbSignature)

	
	gitHubSecret := os.Getenv("my_git_hubSecret")
	if gitHubSecret == "" {
		log.Println("Webhook secret is missing")
		return false
	}

	// mine hmac signature machine 
	mineNewSignature := hmac.New(sha256.New, []byte(gitHubSecret))
		// filling the data 
	mineNewSignature.Write(body)

	// converting the signature according to the github 
	expectedSignature := "sha256=" + hex.EncodeToString(mineNewSignature.Sum(nil))

	fmt.Println("My Signature      :", expectedSignature)

	
	if !hmac.Equal([]byte(expectedSignature), []byte(GitHUbSignature)) {

		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid signature",
		})

		return false
	}

	log.Println("Signature Verified Successfully")

	return true
}