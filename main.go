package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/blang/semver"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/c2-cli/apiclient/openapi"
	"code.storageos.net/storageos/c2-cli/cmd"
	"code.storageos.net/storageos/c2-cli/config/environment"
)

var (
	// Version is the semantic version string which has been assigned to the
	// cli application.
	Version string
	// UserAgentPrefix is used by the CLI application to identify itself to
	// StorageOS.
	UserAgentPrefix string = "storageos-cli"
)

func main() {
	// Determine the cli app version from the embedded semver.
	// If it errors 0.0.0 is returned
	version, _ := semver.Make(Version)

	userAgent := strings.Join([]string{UserAgentPrefix, version.String()}, "/")

	// Initialise the configuration providers
	cfg := environment.NewProvider()

	// Construct the API client.
	transport, err := openapi.NewOpenAPI(cfg, userAgent)
	if err != nil {
		fmt.Printf("failure occurred during initialisation of transport: %v\n", err)
		os.Exit(1)
	}

	client := apiclient.New(
		transport,
		cfg,
	)
	if err != nil {
		fmt.Printf("failure occurred during initialisation of api client: %v\n", err)
		os.Exit(1)
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
