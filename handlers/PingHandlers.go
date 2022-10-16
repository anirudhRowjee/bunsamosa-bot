package handlers

import(
	"net/http"
	"log"
)

func PingHandler(response http.ResponseWriter , request *http.Request){
	log.Println("Received Pin request!")
	response.Write([]byte("Pong UwU"))
	response.WriteHeader(http.StatusOK)
}

