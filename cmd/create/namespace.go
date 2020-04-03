package create

import (
	"context"
	"errors"
	"io"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/cmd/argwrappers"
	"code.storageos.net/storageos/c2-cli/cmd/runwrappers"
	"code.storageos.net/storageos/c2-cli/output"
	"code.storageos.net/storageos/c2-cli/pkg/labels"
)

var (
	errNamespaceNameSpecifiedWrong = errors.New("must specify exactly one name for the namespace")
)

type namespaceCommand struct {
	config  ConfigProvider
	client  Client
	display Displayer

	// Core namespace configuration settings
	labelPairs []string

	writer io.Writer
}

func (c *namespaceCommand) runWithCtx(ctx context.Context, cmd *cobra.Command, args []string) error {

	// Convert the flag values to the desired types/units
	labelSet, err := labels.NewSetFromPairs(c.labelPairs)
	if err != nil {
		return err
	}

	name := args[0]

	ns, err := c.client.CreateNamespace(ctx, name, labelSet)
	if err != nil {
		return err
	}

	return c.display.CreateNamespace(ctx, c.writer, output.NewNamespace(ns))
}

func newNamespace(w io.Writer, client Client, config ConfigProvider) *cobra.Command {
	c := &namespaceCommand{
		config: config,
		client: client,
		writer: w,
	}

	cobraCommand := &cobra.Command{
		Use:   "namespace",
		Short: "Provision a new namespace",
		Example: `
$ storageos create namespace --labels env=prod,rack=db-1 my-namespace-name
`,

		Args: argwrappers.WrapInvalidArgsError(func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errNamespaceNameSpecifiedWrong
			}
			return nil
		}),
		PreRunE: argwrappers.WrapInvalidArgsError(func(_ *cobra.Command, _ []string) error {
			c.display = SelectDisplayer(c.config)

			return nil
		}),

		RunE: func(cmd *cobra.Command, args []string) error {
			run := runwrappers.Chain(
				runwrappers.RunWithTimeout(c.config),
				runwrappers.EnsureNamespaceSetWhenUseIDs(c.config),
				runwrappers.AuthenticateClient(c.config, c.client),
				runwrappers.HandleLicenceError(client),
			)(c.runWithCtx)

			return run(context.Background(), cmd, args)
		},

		SilenceUsage: true,
	}

	cobraCommand.Flags().StringSliceVarP(&c.labelPairs, "labels", "l", []string{}, "an optional set of labels to assign to the new namespace, provided as a comma-separated list of key=value pairs")

	return cobraCommand
}
