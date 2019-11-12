package cmd

import (
	"context"

	"github.com/blang/semver"
	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/cluster"
	"code.storageos.net/storageos/c2-cli/cmd/describe"
	"code.storageos.net/storageos/c2-cli/cmd/get"
	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/volume"
)

// Client defines the functionality required by the CLI application to
// reasonably implement the commands it provides.
type Client interface {
	GetCluster(context.Context) (*cluster.Resource, error)
	GetNode(context.Context, id.Node) (*node.Resource, error)
	GetListNodes(context.Context, ...id.Node) ([]*node.Resource, error)
	GetVolume(context.Context, id.Namespace, id.Volume) (*volume.Resource, error)
	GetAllVolumes(context.Context) ([]*volume.Resource, error)
	GetNamespaceVolumes(context.Context, id.Namespace, ...id.Volume) ([]*volume.Resource, error)

	DescribeNode(context.Context, id.Node) (*node.State, error)
	DescribeListNodes(context.Context, ...id.Node) ([]*node.State, error)
}

// Init configures the CLI application's commands from the root down, using
// client as the method of communicating with the StorageOS API and display
// as the method for formatting and writing the results.
func Init(client Client, version semver.Version) *cobra.Command {
	app := &cobra.Command{
		Use: "storageos <command>",
		Short: `Storage for Cloud Native Applications.

By using this product, you are agreeing to the terms of the the StorageOS Ltd. End
User Subscription Agreement (EUSA) found at: https://storageos.com/legal/#eusa

To be notified about stable releases and latest features, sign up at https://my.storageos.com.
`,
		Version: version.String(),
	}

	app.AddCommand(
		get.NewCommand(client),
		describe.NewCommand(client),
	)

	return app
}
