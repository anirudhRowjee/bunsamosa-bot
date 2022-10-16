package globals



import (
	"log"
	"strconv"
	// "github.com/go-playground/webhooks/v6"
	v3 "github.com/google/go-github/v47/github"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)


type App struct {
	WebhookSecret string
	AppID         int
	OrgID         int
	CertPath      string
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
