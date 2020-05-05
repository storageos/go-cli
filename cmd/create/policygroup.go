package create

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/cmd/argwrappers"
	"code.storageos.net/storageos/c2-cli/cmd/runwrappers"
	"code.storageos.net/storageos/c2-cli/output"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/policygroup"
)

var (
	errPolicyGroupNameWrong         = errors.New("must specify exactly one name for the policy group")
	errPolicyGroupSpecWrong         = errors.New("rule must be a triple separated by colon")
	errPolicyGroupSpecResourceWrong = errors.New("second string in the rule triple must be one in [* volume policy user namespace node cluster]")
	errPolicyGroupSpecReadOnlyWrong = errors.New("third string in the rule triple must be 'r' (read only), 'rw' or 'w' (read-write)")
)

func newErrMissingNamespace(ns string) error {
	return fmt.Errorf(`referenced namespace "%s" does not exist`, ns)
}

type policyGroupCommand struct {
	config  ConfigProvider
	client  Client
	display Displayer

	rules []string

	writer io.Writer
}

func (c *policyGroupCommand) runWithCtx(ctx context.Context, cmd *cobra.Command, args []string) error {
	useIDs, err := c.config.UseIDs()
	if err != nil {
		return err
	}

	// need namespaces to retrieve the ID of the namespace from the name
	// and to fill the names of namespaces in the output type
	namespaces, err := c.client.ListNamespaces(ctx)
	if err != nil {
		return err
	}

	// mapping dictionary to retrieve many IDs by name
	name2ID := make(map[string]id.Namespace)
	for _, ns := range namespaces {
		name2ID[ns.Name] = ns.ID
	}

	// Convert the string triples to policygroup.Specs
	specs := make([]*policygroup.Spec, 0, len(c.rules))
	for _, r := range c.rules {
		spec, err := c.stringToSpec(r, useIDs, name2ID)
		if err != nil {
			return err
		}
		specs = append(specs, spec)
	}

	pg, err := c.client.CreatePolicyGroup(ctx, args[0], specs)
	if err != nil {
		return err
	}

	return c.display.CreatePolicyGroup(ctx, c.writer, output.NewPolicyGroup(pg, namespaces))
}

func (c *policyGroupCommand) stringToSpec(s string, useIDs bool, mapping map[string]id.Namespace) (*policygroup.Spec, error) {
	// input string has already been validated

	triple := strings.Split(s, ":")
	ns, res, rw := triple[0], triple[1], triple[2]

	var namespaceID id.Namespace

	switch {
	case ns == "*":
		namespaceID = "*"
	case useIDs:
		namespaceID = id.Namespace(ns)
	default:
		nsID, ok := mapping[ns]
		if !ok {
			return nil, newErrMissingNamespace(ns)
		}
		namespaceID = nsID
	}

	return &policygroup.Spec{
		NamespaceID:  namespaceID,
		ResourceType: res,
		ReadOnly:     rw == "r", // w or rw means ReadOnly: false
	}, nil
}

func newPolicyGroup(w io.Writer, client Client, config ConfigProvider) *cobra.Command {
	c := &policyGroupCommand{
		config: config,
		client: client,
		writer: w,
	}

	cobraCommand := &cobra.Command{
		Use:   "policy-group",
		Short: "Provision a new policy group",
		Example: `
$ storageos create policy-group -r 'namespace-name:*:r' -r 'namespace-name-2:volume:w'  my-policy-group-name
$ storageos create policy-group -r 'namespace-name:*:r,namespace-name-2:volume:w'  my-policy-group-name
`,

		Args: argwrappers.WrapInvalidArgsError(func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errPolicyGroupNameWrong
			}

			for _, r := range c.rules {
				triple := strings.Split(r, ":")
				if len(triple) != 3 {
					return errPolicyGroupSpecWrong
				}

				// check resource validity
				switch triple[1] {
				case "*":
				case "volume":
				case "policy":
				case "user":
				case "namespace":
				case "node":
				case "cluster":
				default:
					return errPolicyGroupSpecResourceWrong
				}

				// check readOnly validity
				switch triple[2] {
				case "r":
				case "w":
				case "rw":
				default:
					return errPolicyGroupSpecReadOnlyWrong
				}
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
				runwrappers.AuthenticateClient(c.config, c.client),
			)(c.runWithCtx)

			return run(context.Background(), cmd, args)
		},

		SilenceUsage: true,
	}

	cobraCommand.Flags().StringSliceVarP(&c.rules, "rules", "r", []string{}, "set of rules to assign to the new policy group, provided as a comma-separated list of namespace:resource:rw triples.")

	return cobraCommand
}
