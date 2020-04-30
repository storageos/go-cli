package get

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/cmd/runwrappers"
	"code.storageos.net/storageos/c2-cli/output"
)

type licenceCommand struct {
	config  ConfigProvider
	client  Client
	display Displayer

	writer io.Writer
}

func (c *licenceCommand) runWithCtx(ctx context.Context, cmd *cobra.Command, _ []string) error {
	lic, err := c.client.GetLicence(ctx)
	if err != nil {
		return err
	}

	return c.display.GetLicence(ctx, c.writer, output.NewLicence(lic))
}

func newLicence(w io.Writer, client Client, config ConfigProvider) *cobra.Command {
	c := &licenceCommand{
		config: config,
		client: client,
		writer: w,
	}
	cobraCommand := &cobra.Command{
		Use:   "licence",
		Short: "Fetch current licence configuration details",
		Example: `
$ storageos get licence
`,

		PreRun: func(_ *cobra.Command, _ []string) {
			c.display = SelectDisplayer(c.config)
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			run := runwrappers.Chain(
				runwrappers.RunWithTimeout(c.config),
				runwrappers.AuthenticateClient(c.config, c.client),
			)(c.runWithCtx)
			return run(context.Background(), cmd, args)
		},

		// If a legitimate error occurs as part of the VERB cluster command
		// we don't need to barf the usage template.
		SilenceUsage: true,
	}

	return cobraCommand
}
