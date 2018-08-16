package pool

import (
	"context"
	"fmt"

	"github.com/dnephin/cobra"
	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/opts"
)

type createOptions struct {
	name           string
	description    string
	isDefault      bool
	nodeSelector   string
	deviceSelector string
	labels         opts.ListOpts
}

func newCreateCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	opt := createOptions{
		labels: opts.NewListOpts(opts.ValidationPipeline(opts.ValidateEnv, opts.ValidateLabel)),
	}

	cmd := &cobra.Command{
		Use:   "create [OPTIONS] [POOL]",
		Short: "Create a capacity pool",
		Args:  cli.RequiresMaxArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				if opt.name != "" {
					fmt.Fprint(storageosCli.Err(), "Conflicting options: either specify --name or provide positional arg, not both\n")
					return cli.StatusError{StatusCode: 1}
				}
				opt.name = args[0]
			}
			return runCreate(storageosCli, opt)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&opt.name, "name", "", "Pool name")
	flags.Lookup("name").Hidden = true
	flags.StringVar(&opt.nodeSelector, flagNodeSelector, "", "Node selector")
	flags.StringVar(&opt.deviceSelector, flagDeviceSelector, "", "Device selector (on filtered nodes)")
	flags.BoolVar(&opt.isDefault, "default", false, "Set as default pool")
	flags.StringVarP(&opt.description, "description", "d", "", "Pool description")
	flags.Var(&opt.labels, "label", "Set metadata (key=value pairs) on the pool")

	return cmd
}

func runCreate(storageosCli *command.StorageOSCli, opt createOptions) error {
	client := storageosCli.Client()

	params := types.PoolOptions{
		Name:           opt.name,
		Description:    opt.description,
		Default:        opt.isDefault,
		NodeSelector:   opt.nodeSelector,
		DeviceSelector: opt.deviceSelector,
		Labels:         opts.ConvertKVStringsToMap(opt.labels.GetAll()),
		Context:        context.Background(),
	}

	pool, err := client.PoolCreate(params)
	if err != nil {
		return err
	}

	fmt.Fprintf(storageosCli.Out(), "%s\n", pool.Name)
	return nil
}
