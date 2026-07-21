package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	
)


func main 	(){
	// gin router 
	router := gin.Default()
	

	router.GET("/ping",Pong)
	// map the url 
	router.POST("webhook/github",handleGitHubWebHook)
	log.Println("Server Started in port 8090")

	err := router.Run(":8090")
	if err != nil{
		log.Println("Error while running the server",err)
	}


}

// http methods has 3 things 

// Request line 
// header 
// Body 

func handleGitHubWebHook(ctxWrapper* gin.Context){
	ctxWrapper.String(http.StatusOK,"Request received")
	
}
func Pong(ctxWrapper*gin.Context){
	ctxWrapper.String(http.StatusOK,"Pong")
}