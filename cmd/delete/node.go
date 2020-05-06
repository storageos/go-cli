package delete

import (
	"context"
	"errors"
	"io"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/c2-cli/cmd/argwrappers"
	"code.storageos.net/storageos/c2-cli/cmd/flagutil"
	"code.storageos.net/storageos/c2-cli/cmd/runwrappers"
	"code.storageos.net/storageos/c2-cli/output"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/version"
)

type nodeCommand struct {
	config  ConfigProvider
	client  Client
	display Displayer

	// useCAS determines whether the command makes the delete request
	// constrained by the provided casVersion.
	useCAS     func() bool
	casVersion string

	useAsync bool

	writer io.Writer
}

func (c *nodeCommand) runWithCtx(ctx context.Context, cmd *cobra.Command, args []string) error {
	useIDs, err := c.config.UseIDs()
	if err != nil {
		return err
	}

	params := &apiclient.DeleteNodeRequestParams{}

	if c.useCAS() {
		params.CASVersion = version.FromString(c.casVersion)
	}

	// If asynchrony is specified then source the timeout and set the
	// async timeout from it.
	if c.useAsync {
		timeout, err := c.config.CommandTimeout()
		if err != nil {
			return err
		}
		params.AsyncMax = timeout
	}

	nodeID := id.Node(args[0])

	if !useIDs {
		nodeName := args[0]
		n, err := c.client.GetNodeByName(ctx, nodeName)
		if err != nil {
			return err
		}
		nodeID = n.ID
	}

	err = c.client.DeleteNode(ctx, nodeID, params)
	if err != nil {
		return err
	}

	nodeDisplay := output.NodeDeletion{ID: nodeID}

	// Display the "request submitted" message if it was async, instead of
	// the deletion confirmation below.
	if c.useAsync {
		return c.display.DeleteNodeAsync(ctx, c.writer, nodeDisplay)
	}

	return c.display.DeleteNode(ctx, c.writer, nodeDisplay)
}

func newNode(w io.Writer, client Client, config ConfigProvider) *cobra.Command {
	c := &nodeCommand{
		config: config,
		client: client,
		writer: w,
	}

	cobraCommand := &cobra.Command{
		Use:   "node [node name]",
		Short: "Delete a node",
		Example: `
$ storagoes delete node my-old-node
`,

		Args: argwrappers.WrapInvalidArgsError(func(_ *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("must specify exactly one node for deletion")
			}
			return nil
		}),

		PreRunE: argwrappers.WrapInvalidArgsError(func(_ *cobra.Command, args []string) error {
			c.display = SelectDisplayer(c.config)

			return nil
		}),

		RunE: func(cmd *cobra.Command, args []string) error {
			run := runwrappers.Chain(
				runwrappers.RunWithTimeout(c.config),
				runwrappers.AuthenticateClient(c.config, c.client),
			)(c.runWithCtx)
			return run(context.Background(), cmd, args)
		},

		SilenceUsage: true,
	}

	c.useCAS = flagutil.SupportCAS(cobraCommand.Flags(), &c.casVersion)
	flagutil.SupportAsync(cobraCommand.Flags(), &c.useAsync)

	return cobraCommand
}
