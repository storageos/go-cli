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

type namespaceCommand struct {
	config  ConfigProvider
	client  Client
	display Displayer

	// useCAS determines whether the command makes the delete request
	// constrained by the provided casVersion.
	useCAS     func() bool
	casVersion string

	writer io.Writer
}

func (c *namespaceCommand) runWithCtx(ctx context.Context, cmd *cobra.Command, args []string) error {

	var namespaceID id.Namespace

	useIDs, err := c.config.UseIDs()
	if err != nil {
		return err
	}

	if useIDs {
		namespaceID = id.Namespace(args[0])
	} else {
		ns, err := c.client.GetNamespaceByName(ctx, args[0])
		if err != nil {
			return err
		}
		namespaceID = ns.ID
	}

	params := &apiclient.DeleteNamespaceRequestParams{}
	if c.useCAS() {
		params.CASVersion = version.FromString(c.casVersion)
	}

	err = c.client.DeleteNamespace(
		ctx,
		namespaceID,
		params,
	)
	if err != nil {
		return err
	}

	return c.display.DeleteNamespace(ctx, c.writer, output.NamespaceDeletion{ID: namespaceID})
}

func newNamespace(w io.Writer, client Client, config ConfigProvider) *cobra.Command {
	c := &namespaceCommand{
		config: config,
		client: client,
		writer: w,
	}

	cobraCommand := &cobra.Command{
		Use:   "namespace [namespace name]",
		Short: "Delete a namespace",
		Example: `
$ storageos delete namespace my-unneeded-namespace
$ storageos delete namespace --use-ids my-namespace-id
`,

		Args: argwrappers.WrapInvalidArgsError(func(_ *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("must specify exactly one namespace for deletion")
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
				runwrappers.EnsureNamespaceSetWhenUseIDs(c.config),
				runwrappers.AuthenticateClient(c.config, c.client),
			)(c.runWithCtx)
			return run(context.Background(), cmd, args)
		},

		SilenceUsage: true,
	}

	c.useCAS = flagutil.SupportCAS(cobraCommand.Flags(), &c.casVersion)

	return cobraCommand
}
