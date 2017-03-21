package namespace

import (
	"fmt"

	"github.com/dnephin/cobra"
	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/opts"
	"context"
)

type createOptions struct {
	name        string
	displayName string
	description string
	labels      opts.ListOpts
}

func newCreateCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	opt := createOptions{
		labels: opts.NewListOpts(opts.ValidateEnv),
	}

	cmd := &cobra.Command{
		Use:   "create [OPTIONS] [NAMESPACE]",
		Short: "Create a capacity namespace",
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
	flags.StringVar(&opt.name, "name", "", "Specify namespace name")
	flags.Lookup("name").Hidden = true
	flags.StringVar(&opt.displayName, "display-name", "", "Display name of the namespace")
	flags.StringVarP(&opt.description, "description", "d", "", "Namespace description")
	flags.Var(&opt.labels, "label", "Set key:value metadata on the namespace")

	return cmd
}

func runCreate(storageosCli *command.StorageOSCli, opt createOptions) error {
	client := storageosCli.Client()

	params := types.NamespaceCreateOptions{
		Name:        opt.name,
		DisplayName: opt.displayName,
		Description: opt.description,
		Labels:      opts.ConvertKVStringsToMap(opt.labels.GetAll()),
		Context:     context.Background(),
	}

	namespace, err := client.NamespaceCreate(params)
	if err != nil {
		return err
	}

	fmt.Fprintf(storageosCli.Out(), "%s\n", namespace.Name)
	return nil
}
