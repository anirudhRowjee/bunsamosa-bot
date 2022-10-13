package main

import (
	"github.com/bradleyfalzon/ghinstallation/v2"
	// v3 "github.com/google/go-github/v47/github"
	// ghwebhooks "github.com/go-playground/webhooks/v6"
	"gopkg.in/yaml.v2"

	"fmt"
	"io/ioutil"
	"path/filepath"
)

type App struct {
	webhookSecret  string
	appID          int64
	orgID          string
	certPath       string
	installationID int64
	itr            *ghinstallation.Transport
	// `itr` looks to be InstallationTransport
}

func (a *App) parse_from_YAML(path string) {

	// TODO Parse these credentials from YAML
	filename, _ := filepath.Abs(path)
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, a)
	if err != nil {
		panic(err)
	}

	fmt.Println("Parsed", a)
}

func main() {

	// parse YAML File to read in secrets

	myapp := App{}
	myapp.parse_from_YAML("./secrets.yaml")
	fmt.Println("Hello, world!")
}
