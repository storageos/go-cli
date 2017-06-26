package cluster

import (
	"fmt"

	"github.com/dnephin/cobra"

	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"

	"github.com/storageos/go-cli/discovery"
)

type createOptions struct {
	name string
	size int
}

func newCreateCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	opt := createOptions{}

	cmd := &cobra.Command{
		Use:   "create [OPTIONS] [CLUSTER]",
		Short: `Creates a cluster initialization token.`,
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
	flags.StringVar(&opt.name, "name", "", "Cluster name")
	flags.Lookup("name").Hidden = true

	flags.IntVarP(&opt.size, "size", "s", 3, "Cluster size (3-7)")
	return cmd
}

func runCreate(storageosCli *command.StorageOSCli, opt createOptions) error {

	client, err := discovery.NewClient("", "", "")
	if err != nil {
		return err
	}

	token, err := client.ClusterCreate(opt.name, opt.size)
	if err != nil {
		return err
	}

	fmt.Fprintf(storageosCli.Out(), "cluster token: %s\n", token)
	return nil
}
