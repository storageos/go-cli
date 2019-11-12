package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/blang/semver"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/c2-cli/apiclient/openapi"
	"code.storageos.net/storageos/c2-cli/cmd"
	"code.storageos.net/storageos/c2-cli/config"
	"code.storageos.net/storageos/c2-cli/config/environment"
	"code.storageos.net/storageos/c2-cli/config/flags"
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

	// Initialise the configuration provider stack.
	// flags → env → TODO(CP-3918) config file → default
	globalFlags := cmd.InitPersistentFlags()
	configProvider := flags.NewProvider(
		globalFlags,
		environment.NewProvider(
			config.NewDefaulter(),
		),
	)

	// Construct the API client with OpenAPI "transport".
	transport, err := openapi.NewOpenAPI(configProvider, userAgent)
	if err != nil {
		fmt.Printf("failure occurred during initialisation of api client transport: %v\n", err)
		os.Exit(1)
	}

	client := apiclient.New(
		transport,
		configProvider,
	)
	if err != nil {
		fmt.Printf("failure occurred during initialisation of api client: %v\n", err)
		os.Exit(1)
	}

	app := cmd.InitCommand(
		client,
		globalFlags,
		version,
	)

	if err := app.Execute(); err != nil {
		// TODO: Map err to useful exit code
		os.Exit(1)
	}
}
