package handlers

import (
	"context"
	"log"
	"net/http"

	// "github.com/go-playground/webhooks/v6"
	ghwebhooks "github.com/go-playground/webhooks/v6/github"
	v3 "github.com/google/go-github/v47/github"
	"github.com/anirudhRowjee/bunsamosa-bot/globals"

)
func newIssueHandler(parsed_hook *ghwebhooks.IssuesPayload){
	response := "Thank you from Opening this issue! A globalstainer will review it soon!"
	comment := v3.IssueComment{Body: &response}
	_,_, err := globals.MainClient.Issues.CreateComment(context.TODO(), parsed_hook.Repository.Owner.Login,parsed_hook.Repository.Name, int(parsed_hook.Issue.Number), &comment)

	if err!=nil{
		log.Println("Uh oh, a Issues was created and error occured commnenting on it. Repository and Issues number : ", parsed_hook.Repository, parsed_hook.Issue.Number)
	}
}

func WebhookHandler(response http.ResponseWriter, request *http.Request){
	//Creating hook parsers : 
	hook_secret := ghwebhooks.Options.Secret(globals.Myapp.WebhookSecret)
	hook_parser, err := ghwebhooks.New(hook_secret);

	if err!=nil{
		log.Println("Error creating web hook parser, possibly wrong webhook secret in YAML.")
		panic(err)
	}

	//Listing all actions/Events to be parsed : 
	NeededEvents := []ghwebhooks.Event{
		ghwebhooks.IssueCommentEvent, //Not handled
		ghwebhooks.IssuesEvent, //Handled Partially
		ghwebhooks.PullRequestEvent, //Not handled
		ghwebhooks.PullRequestReviewEvent, //Not handles yet
		ghwebhooks.PingEvent, //Not handled
		ghwebhooks.PublicEvent, //Not handled yet
	}

	parsed_hook, err := hook_parser.Parse(request, NeededEvents...)

	if err!=nil{
		if err == ghwebhooks.ErrEventNotFound{
			log.Println("Undefined GitHub event received. err :", err)
			response.WriteHeader(http.StatusOK)
		}else if err == ghwebhooks.ErrEventNotSpecifiedToParse{
			log.Println("Omitting received even, not part of parsing requirements. err", err)
			response.WriteHeader(http.StatusBadRequest)
		}else{
			log.Println("Received malformed event, err :", err)
		}
	}

	switch parsed_hook :=  parsed_hook.(type){
	case ghwebhooks.IssuesPayload:
		go newIssueHandler(&parsed_hook) 
	case ghwebhooks.PingPayload:
		log.Println("[PAYLOAD] Ping ->", parsed_hook)
	case ghwebhooks.PullRequestPayload:
		// TODO Respond with a comment saying congratulations, someone will review your PR soon
		log.Println("[PAYLOAD] There's a Pull Request ->", parsed_hook)
	case ghwebhooks.PublicPayload:
		log.Println("[PAYLOAD] Some Public Event ->", parsed_hook)
	default:
		log.Println("missing handler")

	}
	log.Println("Received Web hook serviced successfully!")
	response.WriteHeader(http.StatusOK)
}
