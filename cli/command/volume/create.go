package volume

import (
	"fmt"

	"context"

	"github.com/dnephin/cobra"
	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/opts"
)

type createOptions struct {
	name         string
	description  string
	size         int
	pool         string
	fsType       string
	namespace    string
	nodeSelector string
	labels       opts.ListOpts
}

func newCreateCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	opt := createOptions{
		labels: opts.NewListOpts(opts.ValidateEnv),
	}

	cmd := &cobra.Command{
		Use:   "create [OPTIONS] [VOLUME]",
		Short: "Create a volume",
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
	flags.StringVar(&opt.name, "name", "", "Volume name")
	flags.Lookup("name").Hidden = true
	flags.StringVarP(&opt.description, "description", "d", "", "Volume description")
	flags.IntVarP(&opt.size, "size", "s", 5, "Volume size in GB")
	flags.StringVarP(&opt.pool, "pool", "p", "default", "Volume capacity pool")
	flags.StringVarP(&opt.fsType, "fstype", "f", "", "Requested filesystem type")
	flags.StringVarP(&opt.namespace, "namespace", "n", "default", "Volume namespace")
	flags.StringVar(&opt.nodeSelector, "nodeSelector", "", "Node selector")
	flags.Var(&opt.labels, "label", "Set metadata (key=value pairs) on the volume")

	return cmd
}

func runCreate(storageosCli *command.StorageOSCli, opt createOptions) error {
	client := storageosCli.Client()

	params := types.VolumeCreateOptions{
		Name:         opt.name,
		Description:  opt.description,
		Size:         opt.size,
		Pool:         opt.pool,
		FSType:       opt.fsType,
		Namespace:    opt.namespace,
		NodeSelector: opt.nodeSelector,
		Labels:       opts.ConvertKVStringsToMap(opt.labels.GetAll()),
		Context:      context.Background(),
	}

	vol, err := client.VolumeCreate(params)
	if err != nil {
		return err
	}

	fmt.Fprintf(storageosCli.Out(), "%s/%s\n", vol.Namespace, vol.Name)
	return nil
}
