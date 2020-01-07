package main

import (
	"fmt"
	"os"

	"github.com/blang/semver"
	"github.com/spf13/pflag"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/c2-cli/cmd"
	"code.storageos.net/storageos/c2-cli/config"
	"code.storageos.net/storageos/c2-cli/config/environment"
	"code.storageos.net/storageos/c2-cli/config/flags"
)

var (
	// Version is the semantic version string which has been assigned to the
	// cli application.
	Version string
)

func main() {
	// Determine the cli app version from the embedded semver.
	// If it errors 0.0.0 is returned
	version, _ := semver.Make(Version)

	// Initialise the configuration provider stack:
	//
	// → flags first. note we init the flagset here because it needs to be
	// given to the InitCommand call (setting up global flags)
	globalFlags := pflag.NewFlagSet("storageos", pflag.ContinueOnError)
	configProvider := flags.NewProvider(
		globalFlags,
		// → environment next
		environment.NewProvider(
			// → TODO(CP-3918) config file next
			//
			// → default values as final fallback
			config.NewDefaulter(),
		),
	)

	client := apiclient.New(configProvider)

	app := cmd.InitCommand(
		client,
		configProvider,
		globalFlags,
		version,
	)

	if err := app.Execute(); err != nil {
		// Attempt to map err to a command error.
		err = cmd.MapCommandError(err)
		// Get the appropriate exit code for the error.
		code := cmd.ExitCodeForError(err)

		fmt.Fprintf(app.OutOrStderr(), "Error: %v\n", err)

		os.Exit(code)
	}
}
