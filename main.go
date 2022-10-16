package main

import (
	"context"

	"log"
	"net/http"
	"strconv"

	// "github.com/go-playground/webhooks/v6"
	ghwebhooks "github.com/go-playground/webhooks/v6/github"
	v3 "github.com/google/go-github/v47/github"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"gopkg.in/yaml.v2"

	"fmt"
	"io/ioutil"
	"path/filepath"
)

type App struct {
	webhookSecret string
	appID         int
	orgID         int
	certPath      string
}

// TODO Write YAML Parsing for environment variables
func (a *App) parse_from_YAML(path string) {

	filename, _ := filepath.Abs(path)
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	var yaml_output map[string]string

	err = yaml.Unmarshal(yamlFile, &yaml_output)
	if err != nil {
		panic(err)
	}

	// TODO Add error reporting here
	log.Println("Parsed", yaml_output)

	a.certPath = yaml_output["certPath"]
	a.webhookSecret = yaml_output["webhookSecret"]

	// TODO better way to do this?
	a.appID, err = strconv.Atoi(yaml_output["appID"])
	if err != nil {
		log.Println("Could not Parse AppID")
		panic(err)
	}
	a.orgID, err = strconv.Atoi(yaml_output["orgID"])
	if err != nil {
		log.Println("Could not Parse orgID")
		panic(err)
	}
}

func main() {

	// parse YAML File to read in secrets
	myapp := App{}
	// Initialize state
	// TODO Separate the YAML Loading from the value setting
	myapp.parse_from_YAML("/root/bunsamosa-bot/secrets.yaml")
	log.Println("YAML Parsed successfully")

	// Initialize the App State
	app_transport, err := ghinstallation.NewAppsTransportKeyFromFile(http.DefaultTransport, int64(myapp.appID), myapp.certPath)
	if err != nil {
		log.Println("Could not Create Github App Client")
		panic(err)
	}
	log.Println("App Transport Initialized")

	// NOTE Don't forget to install the app in your repository before you do this!
	// Initialize the installation
	installation, _, err := v3.NewClient(&http.Client{Transport: app_transport}).Apps.FindOrganizationInstallation(context.TODO(), fmt.Sprint(myapp.orgID))
	if err != nil {
		log.Println("Could not Find Organization installation")
		panic(err)
	}
	log.Println("Organization Transport Initialized")

	// Initialize an authenticated transport for the installation
	installationID := installation.GetID()
	installation_transport := ghinstallation.NewFromAppsTransport(app_transport, installationID)

	log.Printf("successfully initialized GitHub app client, installation-id:%s expected-events:%v\n", fmt.Sprint(installationID), installation.Events)

	log.Println("Installation transport ->", installation_transport)
	log.Println("Hello, world!")

	// Serve!
	// TODO use Higher-Order Functions to generate this response function
	// with the webhook secret from the YAML Parsed into the app in scope
	http.HandleFunc("/Github", func(response http.ResponseWriter, request *http.Request) {

		log.Println("Handling Payload...")

		// Initialize webhook parser
		hook, err := ghwebhooks.New(ghwebhooks.Options.Secret(myapp.webhookSecret))
		if err != nil {
			log.Println("Error Initializing Webhook Parser ->", err)
			return
		}
		log.Println("Webhook Parser Initialized")

		// Parse the incoming request for payload Information, Specifically to check if it's an issue comment
		payload, err := hook.Parse(request, []ghwebhooks.Event{ghwebhooks.IssueCommentEvent, ghwebhooks.IssuesEvent, ghwebhooks.PingEvent, ghwebhooks.PullRequestEvent, ghwebhooks.PublicEvent}...)
		if err != nil {
			if err == ghwebhooks.ErrEventNotFound {
				log.Printf("received unregistered GitHub event: %v\n", err)
				response.WriteHeader(http.StatusOK)
			} else {
				log.Printf("received malformed GitHub event: %v\n", err)
				response.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		log.Println("Payload Parsed")

		// TODO Launch Goroutines with/without pool to handle incoming requests
		switch payload := payload.(type) {

		case ghwebhooks.IssueCommentPayload:
			// TODO Parse for Bounty Assignment
			log.Printf("[PAYLOAD] Someone Commented on an issue -> user %s commented %s on repository %s", payload.Sender.Login, payload.Comment.Body, payload.Repository.FullName)

		case ghwebhooks.IssuesPayload:
			// TODO Respond with a comment saying someone will check it out soon
			log.Printf("[PAYLOAD] Someone Commented on an issue -> user %s Opened an Issue with title %s on repository %s", payload.Sender.Login, payload.Issue.Title, payload.Repository.FullName)

		case ghwebhooks.PingPayload:
			log.Println("[PAYLOAD] Ping ->", payload)
		case ghwebhooks.PullRequestPayload:
			// TODO Respond with a comment saying congratulations, someone will review your PR soon
			log.Println("[PAYLOAD] There's a Pull Request ->", payload)
		case ghwebhooks.PublicPayload:
			log.Println("[PAYLOAD] Some Public Event ->", payload)
		default:
			log.Println("missing handler")
		}
		log.Println("Payload Handling Complete")

		response.WriteHeader(http.StatusOK)
		log.Println("Response Written")

	})

	http.HandleFunc("/ping", func(response http.ResponseWriter, request *http.Request) {
		response.Write([]byte("pong"))
	})

	log.Println("Starting Web Server...")
	err = http.ListenAndServe("0.0.0.0:3000", nil)
	log.Println("Started  Web Server...")
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

}
