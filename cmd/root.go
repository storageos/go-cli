package cmd

import (
	"context"
	"time"

	"github.com/blang/semver"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"code.storageos.net/storageos/c2-cli/cluster"
	"code.storageos.net/storageos/c2-cli/cmd/create"
	"code.storageos.net/storageos/c2-cli/cmd/describe"
	"code.storageos.net/storageos/c2-cli/cmd/get"
	"code.storageos.net/storageos/c2-cli/config"
	"code.storageos.net/storageos/c2-cli/config/flags"
	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/user"
	"code.storageos.net/storageos/c2-cli/volume"
)

// ConfigProvider specifies the configuration settings which commands require
// access to.
type ConfigProvider interface {
	CommandTimeout() (time.Duration, error)
}

// Client defines the functionality required by the CLI application to
// reasonably implement the commands it provides.
type Client interface {
	CreateUser(ctx context.Context, username, password string, withAdmin bool, groups ...id.PolicyGroup) (*user.Resource, error)

	GetCluster(context.Context) (*cluster.Resource, error)
	GetNode(context.Context, id.Node) (*node.Resource, error)
	GetListNodes(context.Context, ...id.Node) ([]*node.Resource, error)
	GetVolume(context.Context, id.Namespace, id.Volume) (*volume.Resource, error)
	GetAllVolumes(context.Context) ([]*volume.Resource, error)
	GetNamespaceVolumes(context.Context, id.Namespace, ...id.Volume) ([]*volume.Resource, error)

	DescribeNode(context.Context, id.Node) (*node.State, error)
	DescribeListNodes(context.Context, ...id.Node) ([]*node.State, error)
}

// InitPersistentFlags builds a flag set containing the global flag set for the
// configuration settings.
func InitPersistentFlags() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("storageos", pflag.ContinueOnError)

	flagSet.StringArray(
		flags.APIEndpointsFlag,
		[]string{config.DefaultAPIEndpoint},
		"set the list of endpoints which are used when connecting to the StorageOS API",
	)
	flagSet.Duration(
		flags.CommandTimeoutFlag,
		config.DefaultCommandTimeout,
		"set the timeout duration to use for execution of the command",
	)
	flagSet.String(
		flags.UsernameFlag,
		config.DefaultUsername,
		"set the StorageOS account username to authenticate as",
	)
	flagSet.String(
		flags.PasswordFlag,
		config.DefaultPassword,
		"set the StorageOS account password to authenticate with",
	)

	return flagSet
}

// InitCommand configures the CLI application's commands from the root down, using
// client as the method of communicating with the StorageOS API.
//
// The returned Command is configured with a flag set containing global configuration settings.
//
// Downstream errors are suppressed, so the caller is responsible for displaying messages.
func InitCommand(client Client, config ConfigProvider, globalFlags *pflag.FlagSet, version semver.Version) *cobra.Command {
	app := &cobra.Command{
		Use: "storageos <command>",
		Short: `Storage for Cloud Native Applications.

By using this product, you are agreeing to the terms of the the StorageOS Ltd. End
User Subscription Agreement (EUSA) found at: https://storageos.com/legal/#eusa

To be notified about stable releases and latest features, sign up at https://my.storageos.com.
`,
		Version: version.String(),

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
