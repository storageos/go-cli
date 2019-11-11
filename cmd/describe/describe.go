package describe

import (
	"context"
	"io"
	"os"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/node"
)

// DescribeClient describes the functionality required by the CLI application
// to reasonably implement the "describe" verb commands.
type DescribeClient interface {
	DescribeNode(context.Context, id.Node) (*node.State, error)
	DescribeListNodes(context.Context, ...id.Node) ([]*node.State, error)
}

// DescribeDisplayer defines the functionality required by the CLI application
// to display the results gathered by the "describe" verb commands.
type DescribeDisplayer interface {
	DescribeNode(io.Writer, *node.State) error
	DescribeNodeList(io.Writer, []*node.State) error
}

// NewCommand configures the set of commands which are grouped by the "describe" verb.
func NewCommand(client DescribeClient) *cobra.Command {
	command := &cobra.Command{
		Use:   "describe",
		Short: "describe retrieves a StorageOS resource, displaying detailed information about it",
	}

	command.AddCommand(
		newNode(os.Stdout, client),
	)

	return command
}
