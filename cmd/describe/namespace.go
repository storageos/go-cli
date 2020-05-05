package describe

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/cmd/argwrappers"
	"code.storageos.net/storageos/c2-cli/cmd/runwrappers"
	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/output"
	"code.storageos.net/storageos/c2-cli/pkg/id"
)

type namespaceCommand struct {
	config  ConfigProvider
	client  Client
	display Displayer

	writer io.Writer
}

func (c *namespaceCommand) runWithCtx(ctx context.Context, cmd *cobra.Command, args []string) error {
	useIDs, err := c.config.UseIDs()
	if err != nil {
		return err
	}

	switch len(args) {

	case 1:
		var ns *namespace.Resource
		var err error

		if useIDs {
			ns, err = c.client.GetNamespace(ctx, id.Namespace(args[0]))
		} else {
			ns, err = c.client.GetNamespaceByName(ctx, args[0])
		}
		if err != nil {
			return err
		}

		return c.display.DescribeNamespace(ctx, c.writer, output.NewNamespace(ns))

	default:
		if useIDs {
			listIDs := make([]id.Namespace, 0, len(args))
			for _, s := range args {
				listIDs = append(listIDs, id.Namespace(s))
			}

			namespaces, err := c.client.GetListNamespacesByUID(ctx, listIDs...)
			if err != nil {
				return err
			}

			return c.display.DescribeListNamespaces(ctx, c.writer, output.NewNamespaces(namespaces))
		}

		namespaces, err := c.client.GetListNamespacesByName(ctx, args...)
		if err != nil {
			return err
		}

		return c.display.DescribeListNamespaces(ctx, c.writer, output.NewNamespaces(namespaces))
	}
}

func newNamespace(w io.Writer, client Client, config ConfigProvider) *cobra.Command {
	c := &namespaceCommand{
		config: config,
		client: client,
		writer: w,
	}

	cobraCommand := &cobra.Command{
		Aliases: []string{"namespaces"},
		Use:     "namespace",
		Short:   "Retrieve detailed information for one or many namespaces",
		Example: `
$ storageos describe namespace my-namespace-name
$ storageos describe namespace --use-ids my-namespace-id
$ storageos describe namespaces
`,

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

	return cobraCommand
}
