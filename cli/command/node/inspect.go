package node

import (
	"errors"
	"strconv"

	"github.com/dnephin/cobra"
	// storageos "github.com/storageos/go-api"
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

	cmd.Flags().StringVarP(&opt.format, "format", "f", "", "Format the output using the given Go template (type --format -h or --help for a detail usage)")

	return cmd
}

func runInspect(storageosCli *command.StorageOSCli, opt inspectOptions) error {
	client := storageosCli.Client()

	getAll := func() (refs []string, getter func(ref string) (interface{}, []byte, error)) {
		nodes, err := client.NodeList(types.ListOptions{})

		for i := range nodes {
			refs = append(refs, strconv.Itoa(i))
		}

		return refs, func(ref string) (interface{}, []byte, error) {
			if err != nil {
				return nil, nil, err
			}

			i, err := strconv.Atoi(ref)
			if err != nil {
				return nil, nil, errors.New("iteration error in node getter function")
			}

			return nodes[i], nil, nil
		}
	}

	getFunc := func(ref string) (interface{}, []byte, error) {
		i, err := client.Node(ref)
		return i, nil, err
	}

	if len(opt.names) == 0 {
		refs, getter := getAll()
		return inspect.Inspect(storageosCli.Out(), refs, opt.format, getter)
	}

	return inspect.Inspect(storageosCli.Out(), opt.names, opt.format, getFunc)
}
