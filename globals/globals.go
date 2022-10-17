package globals

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	// "github.com/go-playground/webhooks/v6"
	"io/ioutil"
	"path/filepath"

	v3 "github.com/google/go-github/v47/github"
	"gopkg.in/yaml.v2"

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
}

var Myapp App
var MainClient *v3.Client

func (a *App) Parse_from_YAML(path string) {

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

	a.CertPath = yaml_output["certPath"]
	a.WebhookSecret = yaml_output["webhookSecret"]

	// TODO better way to do this?
	a.AppID, err = strconv.Atoi(yaml_output["appID"])
	if err != nil {
		log.Println("Could not Parse AppID")
		panic(err)
	}
	a.OrgID, err = strconv.Atoi(yaml_output["orgID"])
	if err != nil {
		log.Println("Could not Parse OrgID")
		panic(err)
	}
}

func (a *App) Initialize_client() {
	// Initialize the Github Client and AppTransport

	app_transport, err := ghinstallation.NewAppsTransportKeyFromFile(http.DefaultTransport, int64(a.AppID), a.CertPath)

	// Push to gloabl scope
	a.AppTransport = app_transport

	if err != nil {
		log.Println("Could not Create Github App Client")
		panic(err)
	}
	log.Println("App Transport Initialized")

	// NOTE Don't forget to install the app in your repository before you do this!
	// Initialize the installation
	installation, _, err := v3.NewClient(&http.Client{Transport: app_transport}).Apps.FindOrganizationInstallation(context.TODO(), fmt.Sprint(a.OrgID))
	if err != nil {
		log.Println("Could not Find Organization installation")
		panic(err)
	}
	log.Println("Organization Transport Initialized")

	// Initialize an authenticated transport for the installation
	installationID := installation.GetID()
	installation_transport := ghinstallation.NewFromAppsTransport(app_transport, installationID)

	a.RuntimeClient = v3.NewClient(&http.Client{Transport: installation_transport})

	log.Printf("successfully initialized GitHub app client, installation-id:%s expected-events:%v\n", fmt.Sprint(installationID), installation.Events)

	log.Println("Installation transport ->", installation_transport)
}
