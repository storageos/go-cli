package describe

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/pkg/id"
)

// ConfigProvider specifies the configuration settings which commands require
// access to.
type ConfigProvider interface {
	CommandTimeout() (time.Duration, error)
}

// Client describes the functionality required by the CLI application
// to reasonably implement the "describe" verb commands.
type Client interface {
	DescribeNode(ctx context.Context, uid id.Node) (*node.State, error)
	DescribeNodeByName(ctx context.Context, name string) (*node.State, error)
	DescribeListNodes(ctx context.Context, uids ...id.Node) ([]*node.State, error)
	DescribeListNodesByName(ctx context.Context, names ...string) ([]*node.State, error)
}

// Displayer defines the functionality required by the CLI application
// to display the results gathered by the "describe" verb commands.
type Displayer interface {
	DescribeNode(context.Context, io.Writer, *node.State) error
	DescribeNodeList(context.Context, io.Writer, []*node.State) error
}

// NewCommand configures the set of commands which are grouped by the "describe" verb.
func NewCommand(client Client, config ConfigProvider) *cobra.Command {
	command := &cobra.Command{
		Use:   "describe",
		Short: "describe retrieves a StorageOS resource, displaying detailed information about it",
	}

	command.AddCommand(
		newNode(os.Stdout, client, config),
	)

	return command
}
