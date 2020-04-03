package apply

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/cluster"
	"code.storageos.net/storageos/c2-cli/output"
	"code.storageos.net/storageos/c2-cli/output/jsonformat"
	"code.storageos.net/storageos/c2-cli/output/textformat"
	"code.storageos.net/storageos/c2-cli/output/yamlformat"
	"code.storageos.net/storageos/c2-cli/user"
)

// ConfigProvider specifies the configuration
type ConfigProvider interface {
	Username() (string, error)
	Password() (string, error)

	CommandTimeout() (time.Duration, error)
	OutputFormat() (output.Format, error)
}

// Client defines the functionality required by the CLI application to
// reasonably implement the "apply" verb commands.
type Client interface {
	Authenticate(ctx context.Context, username, password string) (*user.Resource, error)

	UpdateLicence(ctx context.Context, licenceKey []byte) (*cluster.Licence, error)
}

// Displayer defines the functionality required by the CLI application to
// display the results returned by "apply" verb operations.
type Displayer interface {
	UpdateLicence(ctx context.Context, w io.Writer, licence *output.Licence) error
}

// NewCommand configures the set of commands which are grouped by the "apply" verb.
func NewCommand(client Client, config ConfigProvider) *cobra.Command {
	command := &cobra.Command{
		Use:   "apply",
		Short: "Make changes to existing resources",
	}

	command.AddCommand(
		newLicence(os.Stdout, client, config),
	)

	return command
}

// SelectDisplayer returns the right command displayer specified in the
// config provider.
func SelectDisplayer(cp ConfigProvider) Displayer {
	out, err := cp.OutputFormat()
	if err != nil {
		return textformat.NewDisplayer(textformat.NewTimeFormatter())
	}

	switch out {
	case output.JSON:
		return jsonformat.NewDisplayer(jsonformat.DefaultEncodingIndent)
	case output.YAML:
		return yamlformat.NewDisplayer("")
	case output.Text:
		fallthrough
	default:
		return textformat.NewDisplayer(textformat.NewTimeFormatter())
	}
}
