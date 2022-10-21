package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"

	// "github.com/go-playground/webhooks/v6"
	"github.com/anirudhRowjee/bunsamosa-bot/globals"
	ghwebhooks "github.com/go-playground/webhooks/v6/github"
	v3 "github.com/google/go-github/v47/github"
)

// Function to check if a string is in an array
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

// Function to check if a URL is a Pull Request URL
func is_pull_request(url string) bool {
	// Github Pull Request URLs are of the form
	// https://github.com/<org>/<repo>/pull/<number>
	// If we can verify that the second-last element is a string
	// Then we can verify that the given URL is a pull request URL
	parts := strings.Split(url, "/")
	if contains(parts, "pulls") {
		log.Println("[PR_URLVALID] This is a Pull Request.", parts)
		return true
	} else {
		log.Println("[PR_URLVALID] This is not a Pull Request.", parts)
		return false
	}

}

func newIssueHandler(parsed_hook *ghwebhooks.IssuesPayload) {

	// Generate a New Comment - Text is Customizable

	// TODO Refactor: Add these responses to the App Struct
	response := "Thank you for opening this issue! A Maintainer will review it soon!"
	comment := v3.IssueComment{Body: &response}

	_, _, err := globals.Myapp.RuntimeClient.Issues.CreateComment(context.TODO(), parsed_hook.Repository.Owner.Login, parsed_hook.Repository.Name, int(parsed_hook.Issue.Number), &comment)

	if err != nil {
		log.Printf("[ERROR] Could not Comment on Issue -> Repository [%s] Issue (#%d)[%s]\n", parsed_hook.Repository.FullName, parsed_hook.Issue.Number, parsed_hook.Issue.Title)
	} else {
		log.Printf("[ISSUEHANDLER] Successfully Commented on Issue -> Repository [%s] Issue (#%d)[%s]\n", parsed_hook.Repository.FullName, parsed_hook.Issue.Number, parsed_hook.Issue.Title)
	}
}

func newPRHandler(parsed_hook *ghwebhooks.PullRequestPayload) {

	// Generate a New Comment - Text is Customizable

	// TODO Refactor: Add these responses to the App Struct
	response := "Thank you from Opening this Pull Request, @" + parsed_hook.Sender.Login + " ! A Maintainer will review it soon!"
	comment := v3.IssueComment{Body: &response}

	_, _, err := globals.Myapp.RuntimeClient.Issues.CreateComment(context.TODO(), parsed_hook.Repository.Owner.Login, parsed_hook.Repository.Name, int(parsed_hook.PullRequest.Number), &comment)

	if err != nil {
		log.Printf("[ERROR] Could not Comment on Pull Request -> Repository [%s] PR (#%d)[%s]\n", parsed_hook.Repository.FullName, parsed_hook.PullRequest.Number, parsed_hook.PullRequest.Title)
		log.Println("Error ->", err)
	} else {
		log.Printf("[PRHANDLER] Successfully Commented on Pull Request -> Repository [%s] PR (#%d)[%s]\n", parsed_hook.Repository.FullName, parsed_hook.PullRequest.Number, parsed_hook.PullRequest.Title)
	}
}

func newPRCommentHandler(parsed_hook *ghwebhooks.IssueCommentPayload) {
	// Parse the current webhook

	// List of maintainers
	// TODO Move this to the app
	maintainers := []string{
		"anirudhrowjee",
		"charan2308",
		"suhaskv1",
		"mohamed-ayaan358",
		"typeaway14",
		"navinshrinivas",
		"alfadelta10010",
		"anuragrao04",
		"adarsh-liju",
		"nishtha981",
		"phoenixflame101",
		"crypto-vbg",
		"nigeldias27",
		"bluishhh",
		"himank101",
		"theyashwanthsai",
		"noel-saju",
		"joyenbenitto",
		"mukunddeepak",
		"razerads",
	}

	// Step 1 -> Validate, make sure the issuecomment is on a PR and not on an issue,
	if (parsed_hook.Issue.PullRequest != nil) && is_pull_request(parsed_hook.Issue.PullRequest.URL) && parsed_hook.Action == "created" && contains(maintainers, strings.ToLower(parsed_hook.Sender.Login)) {

		log.Println("A Maintainer Commented -> ")
		log.Printf("[PR_COMMENTHANDLER] Successfully Commented on Pull Request -> Repository [%s] PR (#%d)[%s]\n", parsed_hook.Repository.FullName, parsed_hook.Issue.Number, parsed_hook.Issue.Title)

		// parse the comment here to give a bounty
		comment_text_parts := strings.Split(parsed_hook.Comment.Body, " ")
		if comment_text_parts[0] == "!bounty" {

			// Convert the points
			points, err := strconv.Atoi(comment_text_parts[1])
			if err != nil {
				log.Println("[ERROR][BOUNTY] Invalid Points Assigned -> ", comment_text_parts[1])
			}

			// Assign the bounty points
			err = globals.Myapp.Dbmanager.AssignBounty(
				parsed_hook.Sender.Login,
				parsed_hook.Issue.User.Login,
				parsed_hook.Issue.PullRequest.HTMLURL,
				points,
			)
			if err != nil {
				log.Println("[ERROR][BOUNTY] Could not assign bounty points ->", err)
				return
			}

			log.Printf("[PR_COMMENTHANDLER] Successfully Assigned Bounty on Pull Request -> Repository [%s] PR (#%d)[%s] to user %s for %s points\n", parsed_hook.Repository.FullName, parsed_hook.Issue.Number, parsed_hook.Issue.Title, parsed_hook.Issue.User.Login, comment_text_parts[1])

			response := "Assigned " + comment_text_parts[1] + " Bounty points to user @" + parsed_hook.Issue.User.Login + " !"
			comment := v3.IssueComment{Body: &response}

			_, _, new_err := globals.Myapp.RuntimeClient.Issues.CreateComment(context.TODO(), parsed_hook.Repository.Owner.Login, parsed_hook.Repository.Name, int(parsed_hook.Issue.Number), &comment)
			if new_err != nil {
				log.Printf("[ERROR] Could not Comment on Pull Request -> Repository [%s] PR (#%d)[%s]\n", parsed_hook.Repository.FullName, parsed_hook.Issue.Number, parsed_hook.Issue.Title)
				log.Println("Error ->", new_err)
			} else {
				log.Printf("[PRHANDLER] Successfully Commented on Pull Request -> Repository [%s] PR (#%d)[%s]\n", parsed_hook.Repository.FullName, parsed_hook.Issue.Number, parsed_hook.Issue.Title)
			}

		}
	} else {
		log.Printf("[WARN] Someone else commented on Issue -> Repository [%s] Issue (#%d)[%s]\n", parsed_hook.Repository.FullName, parsed_hook.Issue.Number, parsed_hook.Issue.Title)
	}
	// Return error

}

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
			log.Println("[WARN] Undefined GitHub event received. err :", err)
			response.WriteHeader(http.StatusOK)
			return

		} else if err == ghwebhooks.ErrEventNotSpecifiedToParse {
			// FIXME Unsure about this
			log.Println("[WARN] This event hasn't been specified to parse", err)
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
		if parsed_hook.Action == "opened" {
			log.Printf("[PAYLOAD] Someone Opened an Issue -> user [%s] Opened an Issue with title [%s] on repository [%s]", parsed_hook.Sender.Login, parsed_hook.Issue.Title, parsed_hook.Repository.FullName)
			go newIssueHandler(&parsed_hook)
		} else {
			log.Printf("[PAYLOAD] Non-Open Issue Event -> user [%s] Did something [%s] On an Issue with title [%s] on repository [%s]", parsed_hook.Sender.Login, parsed_hook.Action, parsed_hook.Issue.Title, parsed_hook.Repository.FullName)
		}

	// The API has been Pinged from Github
	case ghwebhooks.PingPayload:
		log.Println("[PAYLOAD] Ping ->", parsed_hook)

	// Someone has opened a new Pull Request
	case ghwebhooks.PullRequestPayload:

		// TODO Respond with a comment saying congratulations, someone will review your PR soon
		if parsed_hook.Action == "opened" {
			log.Printf("[PAYLOAD] Someone Opened an PR -> user [%s] Opened an Issue with title [%s] on repository [%s]", parsed_hook.Sender.Login, parsed_hook.PullRequest.Title, parsed_hook.Repository.FullName)
			go newPRHandler(&parsed_hook)
			// TODO Add handler to assign bounty points
		} else {
			log.Printf("[PAYLOAD] Non-Open PR Event -> user [%s] Did something [%s] On an PR with title [%s] on repository [%s]", parsed_hook.Sender.Login, parsed_hook.Action, parsed_hook.PullRequest.Title, parsed_hook.Repository.FullName)
		}

	// Someone has commented on an Issue
	// We'll be using this webhook for the following -
	// 		- Assigning Bounty to a user
	// 		- Freezing the Leaderboard
	case ghwebhooks.IssueCommentPayload:

		log.Printf("[PAYLOAD] Someone Commented on an issue -> user [%s] commented [%s] on repository [%s]", parsed_hook.Sender.Login, parsed_hook.Comment.Body, parsed_hook.Repository.FullName)

		go newPRCommentHandler(&parsed_hook)

	// The Repository has been made public
	// TODO Consider if we really need this
	case ghwebhooks.PublicPayload:
		log.Println("[PAYLOAD] Some Public Event ->", parsed_hook)

	default:
		log.Println("[WARN] missing handler")

	}

	log.Println("[PAYLOAD] Webhook Has been Handled!")
	response.WriteHeader(http.StatusOK)
}
