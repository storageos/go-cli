package apply

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/apiclient"

	"code.storageos.net/storageos/c2-cli/cmd/flagutil"
	"code.storageos.net/storageos/c2-cli/pkg/version"

	"code.storageos.net/storageos/c2-cli/cmd/runwrappers"
	"code.storageos.net/storageos/c2-cli/output"
)

var (
	errConflictingLicenceSources = errors.New("must specify exactly one input source to read a product licence from")
)

type licenceCommand struct {
	config  ConfigProvider
	client  Client
	display Displayer

	fromStdin    bool
	fromFilepath string

	// useCAS determines whether the command makes the update licence request
	// constrained by the provided casVersion.
	useCAS     func() bool
	casVersion string

	writer io.Writer
}

func (c *licenceCommand) runWithCtx(ctx context.Context, cmd *cobra.Command, args []string) error {
	var licenceKey []byte
	var err error

	if c.fromStdin {
		if licenceKey, err = ioutil.ReadAll(os.Stdin); err != nil {
			return fmt.Errorf("failed to read product licence from stdin: %w", err)
		}
	} else {
		if licenceKey, err = ioutil.ReadFile(filepath.Clean(c.fromFilepath)); err != nil {
			return fmt.Errorf("failed to read product licence from file: %w", err)
		}
	}

	var params *apiclient.UpdateLicenceRequestParams
	if c.useCAS() {
		params = &apiclient.UpdateLicenceRequestParams{
			CASVersion: version.FromString(c.casVersion),
		}
	}

	updated, err := c.client.UpdateLicence(ctx, licenceKey, params)
	if err != nil {
		return err
	}

	return c.display.UpdateLicence(ctx, c.writer, output.NewLicence(updated))
}

func newLicence(w io.Writer, client Client, config ConfigProvider) *cobra.Command {
	c := &licenceCommand{
		config: config,
		client: client,
		writer: w,
	}

	cobraCommand := &cobra.Command{
		Use:   "licence",
		Short: "Apply a product licence to the cluster",
		Example: `
$ storageos apply licence --from-file <path-to-licence-file>

$ echo "<licence file contents>" | storageos apply licence --from-stdin 
		`,

		Args: cobra.NoArgs,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			c.display = SelectDisplayer(c.config)

			if c.fromStdin && c.fromFilepath != "" {
				return errConflictingLicenceSources
			}
			// TODO(CP-3908): If the automated portal licencing is the default
			// option here then the check of "has the user specified a source"
			// can be removed.
			if !c.fromStdin && c.fromFilepath == "" {
				return errors.New("did not specify any input source to read a product licence from")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			run := runwrappers.Chain(
				runwrappers.RunWithTimeout(c.config),
				runwrappers.AuthenticateClient(c.config, c.client),
			)(c.runWithCtx)
			return run(context.Background(), cmd, args)
		},

		SilenceUsage: true,
	}

	cobraCommand.Flags().StringVar(&c.fromFilepath, "from-file", "", "reads a StorageOS product licence from a specified file path")
	cobraCommand.Flags().BoolVar(&c.fromStdin, "from-stdin", false, "reads a StorageOS product licence from the standard input")

	c.useCAS = flagutil.SupportCAS(cobraCommand.Flags(), &c.casVersion)

	return cobraCommand
}
