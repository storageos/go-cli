package volume

import (
	"fmt"

	"context"

	"github.com/spf13/cobra"
	api "github.com/storageos/go-api"
	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/pkg/validation"
)

type removeOptions struct {
	all       bool
	force     bool
	namespace string
	volumes   []string
}

func newRemoveCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	var opt removeOptions

	cmd := &cobra.Command{
		Use:     "rm [OPTIONS] VOLUME [VOLUME...]",
		Aliases: []string{"remove"},
		Short:   "Remove one or more volumes",
		Long:    removeDescription,
		Example: removeExample,
		Args:    cli.RequiresMinArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			opt.volumes = args

			if cmd.Flag("all").Value.String() == "false" && len(opt.volumes) == 0 {
				return fmt.Errorf(
					"\"%s\" requires at least 1 argument(s).\nSee '%s --help'.\n\nUsage:  %s\n\n%s",
					cmd.CommandPath(),
					cmd.CommandPath(),
					cmd.UseLine(),
					cmd.Short,
				)
			}

			return runRemove(storageosCli, &opt)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opt.force, "force", "f", false, "Force the removal of one or more volumes")
	flags.BoolVar(&opt.all, "all", false, "Remove all volumes")
	flags.StringVarP(&opt.namespace, "namespace", "n", "", `Volume namespace (default "default")`)

	return cmd
}

func runRemove(storageosCli *command.StorageOSCli, opt *removeOptions) error {
	client := storageosCli.Client()
	status := 0

	if opt.all {
		listOpts := types.ListOptions{}

		// Set namespace for volume list if specified.
		if opt.namespace != "" {
			listOpts.Namespace = opt.namespace
		}

		volumes, err := storageosCli.Client().VolumeList(listOpts)
		if err != nil {
			fmt.Fprintf(storageosCli.Err(), "%s\n", err)
			status = 1
		}

		if len(volumes) == 0 {
			fmt.Fprintf(storageosCli.Err(), "%s\n", "no volumes to delete")
			status = 1
		} else {
			opt.volumes = make([]string, len(volumes))
			for i, volume := range volumes {
				ref := volume.Namespace + "/" + volume.Name
				opt.volumes[i] = ref
			}
		}
	}

	for _, ref := range opt.volumes {
		namespace, name, err := validation.ParseRefWithDefault(ref)
		if err != nil {
			fmt.Fprintf(storageosCli.Err(), "%s\n", err)
			status = 1
			continue
		}

		// Override default namespace with the specified namespace.
		if opt.namespace != "" {
			namespace = opt.namespace
		}

		params := types.DeleteOptions{
			Name:      name,
			Namespace: namespace,
			Force:     opt.force,
			Context:   context.Background(),
		}

		if err := client.VolumeDelete(params); err != nil {
			fmt.Fprintf(storageosCli.Err(), "%s '%s'\n", err, ref)
			status = 1
			continue
		}

		if api.IsUUID(name) {
			fmt.Fprintf(storageosCli.Out(), "%s\n", name)
		} else {
			fmt.Fprintf(storageosCli.Out(), "%s/%s\n", namespace, name)
		}
	}

	if status != 0 {
		return cli.StatusError{StatusCode: status}
	}
	return nil
}

var removeDescription = `
Remove one or more volumes. You cannot remove a volume that is in use by a container.
`

var removeExample = `
$ storageos volume rm default/testvol
testvol
`
