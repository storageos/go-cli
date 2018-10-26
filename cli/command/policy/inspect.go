package policy

import (
	"github.com/dnephin/cobra"
	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/command/inspect"
)

type inspectOptions struct {
	format   string
	policies []string
}

func newInspectCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	var opt inspectOptions

	cmd := &cobra.Command{
		Use:   "inspect [OPTIONS] POLICY [POLICY...]",
		Short: "Display detailed information on one or more polic(y|ies)",
		Args:  cli.RequiresMinArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			opt.policies = args
			return runInspect(storageosCli, opt)
		},
	}

	cmd.Flags().StringVarP(&opt.format, "format", "f", "", "Format the output using the given Go template")

	return cmd
}

func runInspect(storageosCli *command.StorageOSCli, opt inspectOptions) error {
	client := storageosCli.Client()

	if len(opt.policies) == 0 {
		getAll := func() ([]interface{}, error) {
			policies, err := client.PolicyList(types.ListOptions{})
			if err != nil {
				return nil, err
			}

			res := make([]interface{}, 0, len(policies))
			for _, policy := range policies {
				res = append(res, policy)
			}
			return res, nil
		}
		return inspect.All(storageosCli.Out(), opt.format, getAll)

	}

	getFunc := func(ref string) (interface{}, []byte, error) {
		i, err := client.Policy(ref)
		return i, nil, err
	}

	return inspect.Inspect(storageosCli.Out(), opt.policies, opt.format, getFunc)
}
