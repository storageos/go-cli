package cmd

import (
	"io"

	"github.com/blang/semver"
	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/cmd/describe"
	"code.storageos.net/storageos/c2-cli/cmd/get"
	"code.storageos.net/storageos/c2-cli/pkg/cluster"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/node"
	"code.storageos.net/storageos/c2-cli/pkg/volume"
)

// Client defines the functionality required by the CLI application to
// reasonably implement the commands it provides.
type Client interface {
	GetCluster() (*cluster.Resource, error)
	GetNode(id.Node) (*node.Resource, error)
	GetListNodes(...id.Node) ([]*node.Resource, error)
	GetVolume(id.Namespace, id.Volume) (*volume.Resource, error)
	GetAllVolumes() ([]*volume.Resource, error)
	GetNamespaceVolumes(id.Namespace, ...id.Volume) ([]*volume.Resource, error)

	DescribeNode(id.Node) (*node.State, error)
	DescribeListNodes(...id.Node) ([]*node.State, error)
}

// Displayer defines the functionality required by the CLI application to
// display the results of interaction with the StorageOS API.
type Displayer interface {
	WriteGetCluster(io.Writer, *cluster.Resource) error
	WriteGetNode(io.Writer, *node.Resource) error
	WriteGetNodeList(io.Writer, []*node.Resource) error
	WriteGetVolume(io.Writer, *volume.Resource) error
	WriteGetVolumeList(io.Writer, []*volume.Resource) error

	WriteDescribeNode(io.Writer, *node.State) error
	WriteDescribeNodeList(io.Writer, []*node.State) error
}

// Init configures the CLI application's commands from the root down, using
// client as the method of communicating with the StorageOS API and display
// as the method for formatting and writing the results.
func Init(client Client, display Displayer, version semver.Version) *cobra.Command {
	app := &cobra.Command{
		Use: "storageos <command>",
		Short: `Converged storage for containers.

By using this product, you are agreeing to the terms of the the StorageOS Ltd. End
User Subscription Agreement (EUSA) found at: https://storageos.com/legal/#eusa

To be notified about stable releases and latest features, sign up at https://my.storageos.com.
`,
		Version: version.String(),
	}

	app.AddCommand(
		get.NewCommand(client, display),
		describe.NewCommand(client, display),
	)

	return app
}
