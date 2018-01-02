package volume

import (
	"errors"
	"github.com/dnephin/cobra"
	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/command/inspect"
	"github.com/storageos/go-cli/pkg/validation"
	"strconv"
)

type inspectOptions struct {
	format string
	names  []string
}

func newInspectCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	var opt inspectOptions

	cmd := &cobra.Command{
		Use:   "inspect [OPTIONS] VOLUME [VOLUME...]",
		Short: "Display detailed information on one or more volumes",
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

	getAll := func() (refs []string, getter func(ref string) (interface{}, []byte, error)) {
		vols, err := client.VolumeList(types.ListOptions{})

		for i := range vols {
			refs = append(refs, strconv.Itoa(i))
		}

		return refs, func(ref string) (interface{}, []byte, error) {
			if err != nil {
				return nil, nil, err
			}

			i, err := strconv.Atoi(ref)
			if err != nil {
				return nil, nil, errors.New("iteration error in volume getter function")
			}

			return vols[i], nil, nil
		}

	}

	getFunc := func(ref string) (interface{}, []byte, error) {
		namespace, name, err := validation.ParseRefWithDefault(ref)
		if err != nil {
			return nil, nil, err
		}
		i, err := client.Volume(namespace, name)
		return i, nil, err
	}

	if len(opt.names) == 0 {
		refs, getter := getAll()
		return inspect.Inspect(storageosCli.Out(), refs, opt.format, getter)
	}

	return inspect.Inspect(storageosCli.Out(), opt.names, opt.format, getFunc)
}
