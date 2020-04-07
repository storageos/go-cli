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

type policyGroupCommand struct {
	config  ConfigProvider
	client  Client
	display Displayer

	// useCAS determines whether the command makes the delete request
	// constrained by the provided casVersion.
	useCAS     func() bool
	casVersion string

	writer io.Writer
}

func (c *policyGroupCommand) runWithCtx(ctx context.Context, cmd *cobra.Command, args []string) error {

	var policyGroupID id.PolicyGroup

	useIDs, err := c.config.UseIDs()
	if err != nil {
		return err
	}

	if useIDs {
		policyGroupID = id.PolicyGroup(args[0])
	} else {
		pg, err := c.client.GetPolicyGroupByName(ctx, args[0])
		if err != nil {
			return err
		}
		policyGroupID = pg.ID
	}

	params := &apiclient.DeletePolicyGroupRequestParams{}
	if c.useCAS() {
		params.CASVersion = version.FromString(c.casVersion)
	}

	err = c.client.DeletePolicyGroup(
		ctx,
		policyGroupID,
		params,
	)
	if err != nil {
		return err
	}

	return c.display.DeletePolicyGroup(ctx, c.writer, output.PolicyGroupDeletion{ID: policyGroupID})
}

func newPolicyGroup(w io.Writer, client Client, config ConfigProvider) *cobra.Command {
	c := &policyGroupCommand{
		config: config,
		client: client,
		writer: w,
	}

	cobraCommand := &cobra.Command{
		Use:   "policy-group [policy group name]",
		Short: "Delete a policy group",
		Example: `
$ storageos delete policy-group my-unneeded-policy-group
$ storageos delete policy-group --use-ids my-policy-group-id
`,

		Args: argwrappers.WrapInvalidArgsError(func(_ *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("must specify exactly one policy group for deletion")
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

	return cobraCommand
}
