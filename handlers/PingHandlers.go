package handlers

import (
	"log"
	"net/http"
)

func PingHandler(response http.ResponseWriter, request *http.Request) {
	log.Println("[PING] Received Ping request!")
	response.Write([]byte("Pong UwU"))
	// response.WriteHeader(http.StatusOK)
}
