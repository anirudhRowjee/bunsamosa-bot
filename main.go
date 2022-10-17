package main

import (
	"log"
	"net/http"

	// "github.com/go-playground/webhooks/v6"

	"github.com/anirudhRowjee/bunsamosa-bot/globals"
	"github.com/anirudhRowjee/bunsamosa-bot/handlers"
)

// TODO Write YAML Parsing for environment variables

func main() {

	// parse YAML File to read in secrets
	// Initialize state
	// TODO Separate the YAML Loading from the value setting

	globals.Myapp = globals.App{}

	globals.Myapp.Parse_from_YAML("/root/bunsamosa-bot/secrets.yaml")
	log.Println("YAML Parsed successfully")

	// Initialize the Github Client
	globals.Myapp.Initialize_client()

	// Serve!
	// TODO use Higher-Order Functions to generate this response function
	// with the webhook secret from the YAML Parsed into the app in scope
	http.HandleFunc("/Github", handlers.WebhookHandler)
	http.HandleFunc("/ping", handlers.PingHandler)

	log.Println("Starting Web Server...")

	err := http.ListenAndServe("0.0.0.0:3000", nil)
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

}
