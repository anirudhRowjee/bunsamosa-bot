package handlers

import (
	"context"
	"log"
	"net/http"

	// "github.com/go-playground/webhooks/v6"
	"github.com/anirudhRowjee/bunsamosa-bot/globals"
	ghwebhooks "github.com/go-playground/webhooks/v6/github"
	v3 "github.com/google/go-github/v47/github"
)

func newIssueHandler(parsed_hook *ghwebhooks.IssuesPayload) {

	// Generate a New Comment - Text is Customizable

	// TODO Refactor: Add these responses to the App Struct
	response := "Thank you from Opening this issue! A Maintainer will review it soon!"
	comment := v3.IssueComment{Body: &response}

	_, _, err := globals.Myapp.RuntimeClient.Issues.CreateComment(context.TODO(), parsed_hook.Repository.Owner.Login, parsed_hook.Repository.Name, int(parsed_hook.Issue.Number), &comment)

	if err != nil {
		log.Printf("[ERROR] Could not Comment on Issue -> Repository [%s] Issue (#%d)[%s]\n", parsed_hook.Repository.FullName, parsed_hook.Issue.Number, parsed_hook.Issue.Title)
	}
}

/*
func newPullRequestHandler(parsed_hook *ghwebhooks.PullRequestPayload) {
	// TODO
}

func newCommentHandler(parsed_hook *ghwebhooks.PullRequestPayload) {
	// TODO
}
*/

func WebhookHandler(response http.ResponseWriter, request *http.Request) {

	//Creating hook parsers :
	hook_secret := ghwebhooks.Options.Secret(globals.Myapp.WebhookSecret)
	hook_parser, err := ghwebhooks.New(hook_secret)
	if err != nil {
		log.Println("[ERROR] Webhook parser creation Failed")
		panic(err)
	}

	//Listing all actions/Events to be parsed :
	NeededEvents := []ghwebhooks.Event{
		ghwebhooks.IssueCommentEvent,      // STATUS: Not handled
		ghwebhooks.IssuesEvent,            // STATUS: Handled
		ghwebhooks.PullRequestEvent,       // STATUS: Not handled
		ghwebhooks.PullRequestReviewEvent, // STATUS: POTENTIALLY WILL NOT HANDLE
		ghwebhooks.PingEvent,              // STATUS: Not Handled
		ghwebhooks.PublicEvent,            // STATUS: WILL NOT HANDLE
	}

	parsed_hook, err := hook_parser.Parse(request, NeededEvents...)

	if err != nil {

		if err == ghwebhooks.ErrEventNotFound {
			log.Println("[ERROR] Undefined GitHub event received. err :", err)
			response.WriteHeader(http.StatusOK)
			return

		} else if err == ghwebhooks.ErrEventNotSpecifiedToParse {
			// FIXME Unsure about this
			log.Println("[LOG] ", err)
			response.WriteHeader(http.StatusBadRequest)
			return

		} else {
			log.Printf("[ERROR] received malformed GitHub event: %v\n", err)
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	switch parsed_hook := parsed_hook.(type) {

	// A new issue has been opened.
	case ghwebhooks.IssuesPayload:
		log.Printf("[PAYLOAD] Someone Opened an Issue -> user [%s] Opened an Issue with title [%s] on repository [%s]", parsed_hook.Sender.Login, parsed_hook.Issue.Title, parsed_hook.Repository.FullName)
		go newIssueHandler(&parsed_hook)

	// The API has been Pinged from Github
	case ghwebhooks.PingPayload:
		log.Println("[PAYLOAD] Ping ->", parsed_hook)

	// Someone has opened a new Pull Request
	case ghwebhooks.PullRequestPayload:
		// TODO Respond with a comment saying congratulations, someone will review your PR soon
		log.Println("[PAYLOAD] There's a Pull Request ->", parsed_hook)

	// Someone has commented on an Issue
	// We'll be using this webhook for the following -
	// 		- Assigning Bounty to a user
	// 		- Freezing the Leaderboard
	case ghwebhooks.IssueCommentPayload:
		log.Printf("[PAYLOAD] Someone Commented on an issue -> user [%s] commented [%s] on repository [%s]", parsed_hook.Sender.Login, parsed_hook.Comment.Body, parsed_hook.Repository.FullName)

	// The Repository has been made public
	// TODO Consider if we really need this
	case ghwebhooks.PublicPayload:
		log.Println("[PAYLOAD] Some Public Event ->", parsed_hook)

	default:
		log.Println("missing handler")

	}

	log.Println("Webhook Has been Handled!")
	response.WriteHeader(http.StatusOK)
}
