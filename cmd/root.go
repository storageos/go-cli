package cmd

import (
	"strings"

	"github.com/blang/semver"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/c2-cli/apiclient/openapi"
	"code.storageos.net/storageos/c2-cli/cmd/create"
	"code.storageos.net/storageos/c2-cli/cmd/describe"
	"code.storageos.net/storageos/c2-cli/cmd/get"
	"code.storageos.net/storageos/c2-cli/config"
)

// UserAgentPrefix is used by the CLI application to identify itself to
// StorageOS.
var UserAgentPrefix string = "storageos-cli"

// InitCommand configures the CLI application's commands from the root down, using
// client as the method of communicating with the StorageOS API.
//
// The returned Command is configured with a flag set containing global configuration settings.
//
// Downstream errors are suppressed, so the caller is responsible for displaying messages.
func InitCommand(client *apiclient.Client, config config.Provider, globalFlags *pflag.FlagSet, version semver.Version) *cobra.Command {
	app := &cobra.Command{
		Use: "storageos <command>",
		Short: `Storage for Cloud Native Applications.

By using this product, you are agreeing to the terms of the the StorageOS Ltd. End
User Subscription Agreement (EUSA) found at: https://storageos.com/legal/#eusa

To be notified about stable releases and latest features, sign up at https://my.storageos.com.
`,
		Version: version.String(),

		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			userAgent := strings.Join([]string{UserAgentPrefix, version.String()}, "/")

			transport, err := openapi.NewOpenAPI(config, userAgent)
			if err != nil {
				return err
			}

			return client.ConfigureTransport(transport)
		},

		SilenceErrors: true,
	}

	app.AddCommand(
		create.NewCommand(client, config),
		get.NewCommand(client, config),
		describe.NewCommand(client, config),
	)

	app.PersistentFlags().AddFlagSet(globalFlags)

	return app
}
