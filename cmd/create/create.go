package create

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/cluster"
	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/labels"
	"code.storageos.net/storageos/c2-cli/user"
	"code.storageos.net/storageos/c2-cli/volume"
)

// ConfigProvider specifies the configuration settings which commands require
// access to.
type ConfigProvider interface {
	CommandTimeout() (time.Duration, error)
	UseIDs() (bool, error)
	Namespace() (string, error)
}

// Client describes the functionality required by the CLI application
// to reasonably implement the "create" verb commands.
type Client interface {
	CreateUser(ctx context.Context, username, password string, withAdmin bool, groups ...id.PolicyGroup) (*user.Resource, error)
	CreateVolume(ctx context.Context, namespace id.Namespace, name, description string, fs volume.FsType, sizeBytes uint64, labelSet labels.Set) (*volume.Resource, error)

	GetCluster(ctx context.Context) (*cluster.Resource, error)
	GetNamespaceByName(ctx context.Context, name string) (*namespace.Resource, error)
}

// Displayer describes the functionality required by the CLI application
// to display the resources produced by the "create" verb commands.
type Displayer interface {
	CreateUser(ctx context.Context, w io.Writer, resource *user.Resource) error
	CreateVolume(ctx context.Context, w io.Writer, resource *volume.Resource) error
}

// NewCommand configures the set of commands which are grouped by the "create"
// verb.
func NewCommand(client Client, config ConfigProvider) *cobra.Command {
	command := &cobra.Command{
		Use:   "create",
		Short: "create new resources",
	}

	command.AddCommand(
		newUser(os.Stdout, client, config),
		newVolume(os.Stdout, client, config),
	)

	return command
}
