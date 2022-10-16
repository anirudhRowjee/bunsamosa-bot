package main

import (
	"context"
	"log"
	"net/http"
	// "github.com/go-playground/webhooks/v6"
	v3 "github.com/google/go-github/v47/github"
	"github.com/bradleyfalzon/ghinstallation/v2"
	"fmt"
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

	// Initialize the App State
	app_transport, err := ghinstallation.NewAppsTransportKeyFromFile(http.DefaultTransport, int64(globals.Myapp.AppID), globals.Myapp.CertPath)
	if err != nil {
		log.Println("Could not Create Github App Client")
		panic(err)
	}
	log.Println("App Transport Initialized")

	// NOTE Don't forget to install the app in your repository before you do this!
	// Initialize the installation
	installation, _, err := v3.NewClient(&http.Client{Transport: app_transport}).Apps.FindOrganizationInstallation(context.TODO(), fmt.Sprint(Myapp.OrgID))
	if err != nil {
		log.Println("Could not Find Organization installation")
		panic(err)
	}
	log.Println("Organization Transport Initialized")

	// Initialize an authenticated transport for the installation
	installationID := installation.GetID()
	installation_transport := ghinstallation.NewFromAppsTransport(app_transport, installationID)
	globals.MainClient = v3.NewClient(&http.Client{Transport: installation_transport})

	log.Printf("successfully initialized GitHub app client, installation-id:%s expected-events:%v\n", fmt.Sprint(installationID), installation.Events)

	log.Println("Installation transport ->", installation_transport)
	log.Println("Hello, world!")

	// Serve!
	// TODO use Higher-Order Functions to generate this response function
	// with the webhook secret from the YAML Parsed into the app in scope
	http.HandleFunc("/Github",handlers.WebhookHandler)
	http.HandleFunc("/ping", handlers.PingHandler)

	log.Println("Starting Web Server...")
	err = http.ListenAndServe("0.0.0.0:3000", nil)
	log.Println("Started  Web Server...")
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

}
