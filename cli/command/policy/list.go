package policy

import (
	"github.com/dnephin/cobra"
	"github.com/storageos/go-cli/api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/command/formatter"
)

type listOptions struct {
	format string
}

func newListCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	opt := listOptions{}

	cmd := &cobra.Command{
		Use:     "ls [OPTIONS]",
		Aliases: []string{"list"},
		Short:   "List policies",
		Args:    cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(storageosCli, opt)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opt.format, "format", "", "Pretty-print rules using a Go template")

	return cmd
}

func runList(storageosCli *command.StorageOSCli, opt listOptions) error {
	client := storageosCli.Client()

	params := types.ListOptions{}

	policies, err := client.PolicyList(params)
	if err != nil {
		return err
	}

	format := opt.format
	if len(format) == 0 {
		if len(storageosCli.ConfigFile().PoliciesFormat) > 0 {
			format = storageosCli.ConfigFile().PoliciesFormat
		} else {
			format = formatter.TableFormatKey
		}
	}

	return formatter.PolicyWrite(formatter.Context{
		Output: storageosCli.Out(),
		Format: formatter.NewPolicyFormat(format),
	}, policies.GetPoliciesWithID())
}
