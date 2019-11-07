package main

import (
	"fmt"
	"os"
	"time"

	"github.com/blang/semver"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/c2-cli/apiclient/openapi"
	"code.storageos.net/storageos/c2-cli/cmd"
	"code.storageos.net/storageos/c2-cli/config"
)

var (
	// Version is the semantic version string which has been assigned to the
	// cli application.
	//
	// TODO: Don't think that ldflags work for variables in main - may have to
	// try something else here.
	Version string
	// UserAgent is used by the CLI application to identify itself to
	// StorageOS.
	UserAgent string = "storageos-cli-unknown"
)

// defaultTimeout is the standard timeout for a single request to the CLI's API
// client.
const defaultTimeout = 5 * time.Second

func main() {
	// Determine the cli app version from the embedded semver.
	version, err := semver.Make(Version)

	// Initialise the configuration providers
	cfg := &config.Environment{}

	// Construct the API client.
	apiEndpoint, err := cfg.APIEndpoint()
	if err != nil {
		fmt.Printf("failure occurred during initialisation: %v", err)
		os.Exit(1)
	}

	client, err := apiclient.New(
		openapi.NewOpenAPI(apiEndpoint, UserAgent),
		cfg,
	)
	if err != nil {
		fmt.Printf("failure occurred during initialisation: %v", err)
	}

	app := cmd.Init(
		client,
		version,
	)

	err = app.Execute()

	if err != nil {
		os.Exit(1)
	}
}
