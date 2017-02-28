package volume

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
	name        string
	description string
	size        int
	pool        string
	fstype      string
	namespace   string
	labels      opts.ListOpts
}

func newCreateCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	opts := createOptions{
		labels: opts.NewListOpts(opts.ValidateEnv),
	}

	cmd := &cobra.Command{
		Use:   "create [OPTIONS] [VOLUME]",
		Short: "Create a volume",
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
	flags.StringVar(&opts.name, "name", "", "Specify volume name")
	flags.Lookup("name").Hidden = true
	flags.StringVarP(&opts.description, "description", "d", "", "Volume description")
	flags.IntVarP(&opts.size, "size", "s", 5, "Volume size in GB")
	flags.StringVarP(&opts.pool, "pool", "p", "default", "Volume capacity poool")
	flags.StringVarP(&opts.fstype, "fstype", "f", "", "Requested filesystem type")
	flags.StringVarP(&opts.namespace, "namespace", "n", "", "Volume namespace")
	flags.Var(&opts.labels, "label", "Set metadata for a volume")

	return cmd
}

func runCreate(storageosCli *command.StorageOSCli, opts createOptions) error {
	client := storageosCli.Client()

	volReq := types.VolumeCreateOptions{
		Name:        opts.name,
		Description: opts.description,
		Size:        opts.size,
		Pool:        opts.pool,
		FSType:      opts.fstype,
		Namespace:   opts.namespace,
		Labels:      runconfigopts.ConvertKVStringsToMap(opts.labels.GetAll()),
		Context:     context.Background(),
	}

	vol, err := client.VolumeCreate(volReq)
	if err != nil {
		return err
	}

	fmt.Fprintf(storageosCli.Out(), "%s\n", vol.Name)
	return nil
}
