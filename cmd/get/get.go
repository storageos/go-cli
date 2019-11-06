package get

import (
	"io"
	"os"

	"github.com/blang/semver"
	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/node"
)

// GetClient defines the functionality required by the CLI application to
// reasonably implement the "get" verb commands.
type GetClient interface {
	GetNode(uid id.Node) (*node.Resource, error)
	GetListNodes(uids ...id.Node) ([]*node.Resource, error)
}

// GetDisplayer defines the functionality required by the CLI application to
// display the results gathered by the "get" verb commands.
type GetDisplayer interface {
	WriteGetNode(io.Writer, *node.Resource) error
	WriteGetNodeList(io.Writer, []*node.Resource) error
}

// NewCommand configures the set of commands which are grouped by the "get" verb.
func NewCommand(client GetClient, display GetDisplayer, version semver.Version) *cobra.Command {
	command := &cobra.Command{
		Use:   "get",
		Short: "get retrieves a StorageOS resource, displaying basic information about it",
	}

	command.AddCommand(
		newNode(os.Stdout, client, display),
	)

	return command
}
