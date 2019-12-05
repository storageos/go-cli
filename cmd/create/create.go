package create

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/user"
)

// ConfigProvider specifies the configuration settings which commands require
// access to.
type ConfigProvider interface {
	CommandTimeout() (time.Duration, error)
}

// CreateClient describes the functionality required by the CLI application
// to reasonably implement the "create" verb commands.
type CreateClient interface {
	CreateUser(ctx context.Context, username, password string, withAdmin bool, groups ...id.PolicyGroup) (*user.Resource, error)
}

// CreateDisplayer describes the functionality required by the CLI application
// to display the resources produced by the "create" verb commands.
type CreateDisplayer interface {
	CreateUser(context.Context, io.Writer, *user.Resource) error
}

func NewCommand(client CreateClient, config ConfigProvider) *cobra.Command {
	command := &cobra.Command{
		Use:   "create",
		Short: "create requests the creation of a new StorageOS resource",
	}

	command.AddCommand(
		newUser(os.Stdout, client, config),
	)

	return command
}
