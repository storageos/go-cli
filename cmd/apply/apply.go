package apply

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/cluster"
)

// ConfigProvider specifies the configuration
type ConfigProvider interface {
	CommandTimeout() (time.Duration, error)
}

// Client defines the functionality required by the CLI application to
// reasonably implement the "apply" verb commands.
type Client interface {
	UpdateLicence(ctx context.Context, licenceKey []byte) (*cluster.Licence, error)
}

// Displayer defines the functionality required by the CLI application to
// display the results returned by "apply" verb operations.
type Displayer interface {
	UpdateLicence(ctx context.Context, w io.Writer, licence *cluster.Licence) error
}

// NewCommand configures the set of commands which are grouped by the "apply" verb.
func NewCommand(client Client, config ConfigProvider) *cobra.Command {
	command := &cobra.Command{
		Use:   "apply",
		Short: "apply performs an update on an existing StorageOS resource, displaying the new state",
	}

	command.AddCommand(
		newLicence(os.Stdout, client, config),
	)

	return command
}
