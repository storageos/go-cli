package get

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/cmd/flagutil"
	"code.storageos.net/storageos/c2-cli/cmd/runwrappers"
	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/output"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/selectors"
)

type namespaceCommand struct {
	config  ConfigProvider
	client  Client
	display Displayer

	selectors []string

	writer io.Writer
}

func (c *namespaceCommand) runWithCtx(ctx context.Context, cmd *cobra.Command, args []string) error {
	switch len(args) {
	case 1:
		ns, err := c.getNamespace(ctx, args[0])
		if err != nil {
			return err
		}

		return c.display.GetNamespace(ctx, c.writer, output.NewNamespace(ns))
	default:
		set, err := selectors.NewSetFromStrings(c.selectors...)
		if err != nil {
			return err
		}

		namespaces, err := c.listNamespaces(ctx, args)
		if err != nil {
			return err
		}

		namespaces = set.FilterNamespaces(namespaces)

		return c.display.GetListNamespaces(
			ctx,
			c.writer,
			output.NewNamespaces(namespaces),
		)
	}
}

// getNamespace retrieves a single namespace resource using the API client,
// determining whether to retrieve the namespace by name or ID based on config
// settings.
func (c *namespaceCommand) getNamespace(ctx context.Context, ref string) (*namespace.Resource, error) {
	useIDs, err := c.config.UseIDs()
	if err != nil {
		return nil, err
	}

	if !useIDs {
		return c.client.GetNamespaceByName(ctx, ref)
	}

	uid := id.Namespace(ref)
	return c.client.GetNamespace(ctx, uid)
}

// listNamespaces retrieves a list of namespace resources using the API client,
// determining whether to retrieve namespaces by names or IDs based on the
// current config settings.
func (c *namespaceCommand) listNamespaces(ctx context.Context, refs []string) ([]*namespace.Resource, error) {
	useIDs, err := c.config.UseIDs()
	if err != nil {
		return nil, err
	}

	if !useIDs {
		return c.client.GetListNamespacesByName(ctx, refs...)
	}

	uids := make([]id.Namespace, len(refs))
	for i, ref := range refs {
		uids[i] = id.Namespace(ref)
	}

	return c.client.GetListNamespacesByUID(ctx, uids...)
}

func newNamespace(w io.Writer, client Client, config ConfigProvider) *cobra.Command {
	c := &namespaceCommand{
		config: config,
		client: client,
		writer: w,
	}

	cobraCommand := &cobra.Command{
		Aliases: []string{"namespaces"},
		Use:     "namespace [namespace names...]",
		Short:   "Retrieve basic details of cluster namespaces",
		Example: `
$ storageos get namespaces

$ storageos get namespace my-namespace-name
`,
		PreRun: func(_ *cobra.Command, _ []string) {
			c.display = SelectDisplayer(c.config)
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			run := runwrappers.Chain(
				runwrappers.RunWithTimeout(c.config),
				runwrappers.EnsureTargetOrSelectors(&c.selectors),
				runwrappers.AuthenticateClient(c.config, c.client),
			)(c.runWithCtx)
			return run(context.Background(), cmd, args)
		},

		SilenceUsage: true,
	}

	flagutil.SupportSelectors(cobraCommand.Flags(), &c.selectors)

	return cobraCommand
}
