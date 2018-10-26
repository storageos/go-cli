package node

import (
	"github.com/dnephin/cobra"
	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/command/inspect"
)

type inspectOptions struct {
	format string
	names  []string
}

func newInspectCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	var opt inspectOptions

	cmd := &cobra.Command{
		Use:   "inspect [OPTIONS] NODE [NODE...]",
		Short: "Display detailed information on one or more nodes",
		Args:  cli.RequiresMinArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			opt.names = args
			return runInspect(storageosCli, opt)
		},
	}

	cmd.Flags().StringVarP(&opt.format, "format", "f", "", "Format the output using the given Go template")

	return cmd
}

func runInspect(storageosCli *command.StorageOSCli, opt inspectOptions) error {
	client := storageosCli.Client()

	getFunc := func(ref string) (interface{}, []byte, error) {
		i, err := client.Node(ref)
		return i, nil, err
	}

	if len(opt.names) == 0 {
		getAll := func() ([]interface{}, error) {
			nodes, err := client.NodeList(types.ListOptions{})
			if err != nil {
				return nil, err
			}

			res := make([]interface{}, 0, len(nodes))
			for _, node := range nodes {
				res = append(res, node)
			}
			return res, nil
		}
		return inspect.All(storageosCli.Out(), opt.format, getAll)
	}

	return inspect.Inspect(storageosCli.Out(), opt.names, opt.format, getFunc)
}
