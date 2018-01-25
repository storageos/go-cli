package pool

import (
	"fmt"

	"github.com/dnephin/cobra"
	"github.com/storageos/go-cli/api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/opts"
	"context"
)

type createOptions struct {
	name          string
	description   string
	isDefault     bool
	defaultDriver string
	controllers   opts.ListOpts
	drivers       opts.ListOpts
	active        bool
	labels        opts.ListOpts
}

func newCreateCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	opt := createOptions{
		controllers: opts.NewListOpts(opts.ValidateEnv),
		drivers:     opts.NewListOpts(opts.ValidateEnv),
		labels:      opts.NewListOpts(opts.ValidateEnv),
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
	flags.StringVarP(&opt.description, "description", "d", "", "Pool description")
	flags.BoolVar(&opt.isDefault, "default", false, "Set as default pool")
	flags.StringVar(&opt.defaultDriver, "default-driver", "", "Default capacity driver")
	flags.Var(&opt.controllers, "controllers", "Controllers that contribute capacity to the pool")
	flags.Var(&opt.drivers, "drivers", "Drivers providing capacity to the pool")
	flags.BoolVar(&opt.active, "active", true, "Enable or disable the pool")
	flags.Var(&opt.labels, "label", "Set metadata (key=value pairs) on the pool")

	return cmd
}

func runCreate(storageosCli *command.StorageOSCli, opt createOptions) error {
	client := storageosCli.Client()

	params := types.PoolCreateOptions{
		Name:            opt.name,
		Description:     opt.description,
		Default:         opt.isDefault,
		DefaultDriver:   opt.defaultDriver,
		ControllerNames: opt.controllers.GetAll(),
		DriverNames:     opt.drivers.GetAll(),
		Active:          opt.active,
		Labels:          opts.ConvertKVStringsToMap(opt.labels.GetAll()),
		Context:         context.Background(),
	}

	pool, err := client.PoolCreate(params)
	if err != nil {
		return err
	}

	fmt.Fprintf(storageosCli.Out(), "%s\n", pool.Name)
	return nil
}
