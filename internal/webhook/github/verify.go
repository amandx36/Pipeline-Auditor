package github_webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

func VerifySignature(ctx *gin.Context, body []byte) bool {

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
	if _, err := mineNewSignature.Write(body); err != nil {
		log.Println("Error while writing body to HMAC:", err)
		return false
	}
	// converting the signature according to the github
	expectedSignature := "sha256=" + hex.EncodeToString(mineNewSignature.Sum(nil))

	fmt.Println("My Signature      :", expectedSignature)

	if !hmac.Equal([]byte(expectedSignature), []byte(GitHUbSignature)) {

		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid signature",
		})

		return false
	}
	fmt.Println("===================================")
	fmt.Println("GitHub Signature :", GitHUbSignature)
	fmt.Println("My Signature     :", expectedSignature)
	fmt.Println("Signature Match  :", hmac.Equal(
		[]byte(expectedSignature),
		[]byte(GitHUbSignature),
	))
	fmt.Println("===================================")
	log.Println("Signature Verified Successfully  ")

	return true
}
