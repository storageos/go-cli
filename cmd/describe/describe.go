package describe

import (
	"io"
	"os"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/node"
	"code.storageos.net/storageos/c2-cli/pkg/volume"
)

// DescribeClient describes the functionality required by the CLI application
// to reasonably implement the "describe" verb commands.
type DescribeClient interface {
	DescribeNode(id.Node) (*node.Resource, error)
	DescribeListNodes(...id.Node) ([]*node.Resource, error)
	DescribeVolume(id.Namespace, id.Volume) (*volume.Resource, error)
}

// DescribeDisplayer defines the functionality required by the CLI application
// to display the results gathered by the "describe" verb commands.
type DescribeDisplayer interface {
	WriteDescribeNode(io.Writer, *node.Resource) error
	WriteDescribeNodeList(io.Writer, []*node.Resource) error
	WriteDescribeVolume(io.Writer, *volume.Resource) error
}

// NewCommand configures the set of commands which are grouped by the "describe" verb.
func NewCommand(client DescribeClient, display DescribeDisplayer) *cobra.Command {
	command := &cobra.Command{
		Use:   "describe",
		Short: "describe retrieves a StorageOS resource, displaying detailed information about it",
	}

	command.AddCommand(
		newNode(os.Stdout, client, display),
		newVolume(os.Stdout, client, display),
	)

	return command
}
