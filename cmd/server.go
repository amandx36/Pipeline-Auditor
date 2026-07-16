package main

import (
	"log"
	"net/http"
)


func main 	(){
	// new http request and map url with the function 
	server := 	http.NewServeMux();
	
	// map the function with  the url 
	server.HandleFunc("/webhook/github",handleGitHubWebHook);
	// start  the server 
	err := http.ListenAndServe(":8090",server);
	if err != nil {
		log.Println(("Error while Running the server "))
		log.Fatal(err)
	}

	log.Println("Server listening on 8080");
}

// http methods has 3 things 

// Request line 
// header 
// Body 

func handleGitHubWebHook(resp http.ResponseWriter, clientReq  *http.Request){
	// send ok status code 
	resp.WriteHeader(http.StatusOK)
	// write in body 
	resp.Write([]byte("Request received"))
}