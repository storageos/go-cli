package pool

import (
	"fmt"

	"github.com/dnephin/cobra"
	runconfigopts "github.com/docker/docker/runconfig/opts"
	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/opts"
	"golang.org/x/net/context"
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
	opts := createOptions{
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
				if opts.name != "" {
					fmt.Fprint(storageosCli.Err(), "Conflicting options: either specify --name or provide positional arg, not both\n")
					return cli.StatusError{StatusCode: 1}
				}
				opts.name = args[0]
			}
			return runCreate(storageosCli, opts)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&opts.name, "name", "", "Specify pool name")
	flags.Lookup("name").Hidden = true
	flags.StringVarP(&opts.description, "description", "d", "", "Pool description")
	flags.BoolVar(&opts.isDefault, "default", false, "Set as default pool")
	flags.StringVar(&opts.defaultDriver, "default-driver", "", "Default capacity driver")
	flags.Var(&opts.controllers, "controllers", "Controllers that contribute capacity to the pool")
	flags.Var(&opts.drivers, "drivers", "Drivers providing capacity to the pool")
	flags.BoolVar(&opts.active, "active", true, "Enable or disable the pool")
	flags.Var(&opts.labels, "label", "Set metadata for a pool")

	return cmd
}

func runCreate(storageosCli *command.StorageOSCli, opts createOptions) error {
	client := storageosCli.Client()

	params := types.PoolCreateOptions{
		Name:            opts.name,
		Description:     opts.description,
		Default:         opts.isDefault,
		DefaultDriver:   opts.defaultDriver,
		ControllerNames: opts.controllers.GetAll(),
		DriverNames:     opts.drivers.GetAll(),
		Active:          opts.active,
		Labels:          runconfigopts.ConvertKVStringsToMap(opts.labels.GetAll()),
		Context:         context.Background(),
	}

	pool, err := client.PoolCreate(params)
	if err != nil {
		return err
	}

	fmt.Fprintf(storageosCli.Out(), "%s\n", pool.Name)
	return nil
}
