package cmd

import (
	"io"

	"github.com/blang/semver"
	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/cmd/describe"
	"code.storageos.net/storageos/c2-cli/cmd/get"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/node"
)

// Client defines the functionality required by the CLI application to
// reasonably implement the commands it provides.
type Client interface {
	GetNode(uid id.Node) (*node.Resource, error)
	GetListNodes(uids ...id.Node) ([]*node.Resource, error)

	DescribeNode(uid id.Node) (*node.Resource, error)
	DescribeListNodes(uids ...id.Node) ([]*node.Resource, error)
}

// Displayer defines the functionality required by the CLI application to
// display the results of interaction with the StorageOS API.
type Displayer interface {
	WriteGetNode(w io.Writer, resource *node.Resource) error
	WriteGetNodeList(w io.Writer, resources []*node.Resource) error

	WriteDescribeNode(w io.Writer, resource *node.Resource) error
	WriteDescribeNodeList(w io.Writer, resources []*node.Resource) error
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
		get.NewCommand(client, display, version),
		describe.NewCommand(client, display),
	)

	return app
}
