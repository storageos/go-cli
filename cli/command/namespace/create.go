package namespace

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
	displayName string
	description string
	labels      opts.ListOpts
}

func newCreateCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	opts := createOptions{
		labels: opts.NewListOpts(opts.ValidateEnv),
	}

	cmd := &cobra.Command{
		Use:   "create [OPTIONS] [NAMESPACE]",
		Short: "Create a capcity namespace",
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
	flags.StringVar(&opts.name, "name", "", "Specify namespace name")
	flags.Lookup("name").Hidden = true
	flags.StringVar(&opts.displayName, "display-name", "", "Namespace display name")
	flags.StringVarP(&opts.description, "description", "d", "", "Namespace description")
	flags.Var(&opts.labels, "label", "Set metadata for a namespace")

	return cmd
}

func runCreate(storageosCli *command.StorageOSCli, opts createOptions) error {
	client := storageosCli.Client()

	params := types.NamespaceCreateOptions{
		Name:        opts.name,
		DisplayName: opts.displayName,
		Description: opts.description,
		Labels:      runconfigopts.ConvertKVStringsToMap(opts.labels.GetAll()),
		Context:     context.Background(),
	}

	namespace, err := client.NamespaceCreate(params)
	if err != nil {
		return err
	}

	fmt.Fprintf(storageosCli.Out(), "%s\n", namespace.Name)
	return nil
}
