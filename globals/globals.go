package globals

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"io/ioutil"
	"path/filepath"

	v3 "github.com/google/go-github/v47/github"
	"gopkg.in/yaml.v2"

	"github.com/anirudhRowjee/bunsamosa-bot/database"
	"github.com/bradleyfalzon/ghinstallation/v2"
)

type App struct {

	// Initialization Information
	WebhookSecret string
	AppID         int
	OrgID         int
	CertPath      string

	// Runtime Variables and Global Dependencies
	RuntimeClient *v3.Client
	AppTransport  *ghinstallation.AppsTransport

	// TODO Add Database Dependencies
	Db_connection_string string
	Dbmanager            *database.DBManager
}

var Myapp App

func (a *App) Parse_from_YAML(path string) {

	filename, _ := filepath.Abs(path)
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Println("[ERROR] Invalid Secrets YAML Filepath")
		panic(err)
	}

	var yaml_output map[string]string

	err = yaml.Unmarshal(yamlFile, &yaml_output)
	if err != nil {
		log.Println("[ERROR] Could not Unmarshal YAML")
		panic(err)
	}

	// TODO Add error reporting here
	log.Println("[SECRETS] YAML Parsing Complete")

	a.CertPath = yaml_output["certPath"]
	a.WebhookSecret = yaml_output["webhookSecret"]

	// TODO better way to do this?
	a.AppID, err = strconv.Atoi(yaml_output["appID"])
	if err != nil {
		log.Println("[ERROR] Could not Parse AppID")
		panic(err)
	}
	a.OrgID, err = strconv.Atoi(yaml_output["orgID"])
	if err != nil {
		log.Println("[ERROR] Could not Parse OrgID")
		panic(err)
	}

	// Read in the Connection String
	a.Db_connection_string = yaml_output["dbConnectionString"]
}

func (a *App) Initialize_github_client() {
	// Initialize the Github Client and AppTransport
	log.Println("[CLIENT] Initializing Github Client")

	app_transport, err := ghinstallation.NewAppsTransportKeyFromFile(http.DefaultTransport, int64(a.AppID), a.CertPath)

	// Push to gloabl scope
	a.AppTransport = app_transport

	if err != nil {
		log.Println("[ERROR] Could not Create Github App Client")
		panic(err)
	}
	log.Println("[CLIENT] App Transport Initialized")

	// NOTE Don't forget to install the app in your repository before you do this!
	// Initialize the installation
	installation, _, err := v3.NewClient(&http.Client{Transport: app_transport}).Apps.FindOrganizationInstallation(context.TODO(), fmt.Sprint(a.OrgID))
	if err != nil {
		log.Println("[ERROR] Could not Find Organization installation")
		panic(err)
	}
	log.Println("[CLIENT] Organization Transport Initialized")

	// Initialize an authenticated transport for the installation
	installationID := installation.GetID()
	installation_transport := ghinstallation.NewFromAppsTransport(app_transport, installationID)

	a.RuntimeClient = v3.NewClient(&http.Client{Transport: installation_transport})

	log.Printf("[CLIENT] successfully initialized GitHub app client, installation-id:%s expected-events:%v\n", fmt.Sprint(installationID), installation.Events)
}

func (a *App) Initialize_database() {
	// Start the database. Panic on error.

	dbmanager := database.DBManager{}
	log.Println("[DATABASE] Initializing Database Manager")
	err := dbmanager.Init(a.Db_connection_string)
	if err != nil {
		log.Panicln("[DATABASE] DB Initialization Failed ->", err)
	} else {
		a.Dbmanager = &dbmanager
		log.Println("[DATABASE] DB Manager Initialized successfully")
	}
}

func (a *App) Leaderboard_GetAllRecords() ([]database.ContributorRecordModel, error) {

	// Get all the time series data present so far
	// from the database
	var all_records []database.ContributorRecordModel

	// Use the database method
	records, err := a.Dbmanager.Get_all_records()
	if err != nil {
		return nil, err
	} else {
		all_records = records
	}

	return all_records, nil

}

func (a *App) AssignBountyPoints() ([]database.ContributorRecordModel, error) {

	// Get all the time series data present so far
	// from the database
	var all_records []database.ContributorRecordModel

	// Use the database method
	records, err := a.Dbmanager.Get_all_records()
	if err != nil {
		return nil, err
	} else {
		all_records = records
	}

	return all_records, nil

}

func (a *App) Leaderboard_GetMaterialized() ([]database.ContributorModel, error) {

	// Get a materialized view of the leaderboard
	var leaderboard []database.ContributorModel

	records, err := a.Dbmanager.Get_leaderboard()
	if err != nil {
		return nil, err
	} else {
		leaderboard = records
	}

	return leaderboard, nil

}

func (a *App) Leaderboard_GetUserRecord(user string) {
	// Take a user's username and return their records
	// TODO
}
