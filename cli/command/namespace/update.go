package namespace

import (
	"fmt"

	"github.com/dnephin/cobra"
	"github.com/spf13/pflag"
	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/opts"
	"context"
)

const (
	flagDisplayName = "display-name"
	flagDescription = "description"
	flagLabelAdd    = "label-add"
	flagLabelRemove = "label-rm"
)

type updateOptions struct {
	name        string
	displayName string
	description string
	labels      opts.ListOpts
}

func newUpdateCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	opt := updateOptions{
		labels: opts.NewListOpts(opts.ValidateEnv),
	}

	cmd := &cobra.Command{
		Use:   "update [OPTIONS] NAMESPACE",
		Short: "Update a namespace",
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpdate(storageosCli, cmd.Flags(), args[0])
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opt.displayName, flagDisplayName, "", `Display name of the namespace`)
	flags.StringVarP(&opt.description, flagDescription, "d", "", `Namespace description`)
	flags.Var(&opt.labels, flagLabelAdd, "Add or update a namespace label (key=value)")
	labelKeys := opts.NewListOpts(nil)
	flags.Var(&labelKeys, flagLabelRemove, "Remove a node label if exists")
	return cmd
}

func runUpdate(storageosCli *command.StorageOSCli, flags *pflag.FlagSet, name string) error {
	success := func(_ string) {
		fmt.Fprintln(storageosCli.Out(), name)
	}
	return updateNamespaces(storageosCli, []string{name}, mergeNamespaceUpdate(flags), success)
}

func updateNamespaces(storageosCli *command.StorageOSCli, names []string, mergeNamespace func(namespace *types.Namespace) error, success func(name string)) error {
	client := storageosCli.Client()
	ctx := context.Background()

	for _, name := range names {
		namespace, err := client.Namespace(name)
		if err != nil {
			return err
		}

		err = mergeNamespace(namespace)
		if err != nil {
			return err
		}
		params := types.NamespaceCreateOptions{
			Name:        namespace.Name,
			DisplayName: namespace.DisplayName,
			Description: namespace.Description,
			Labels:      namespace.Labels,
			Context:     ctx,
		}
		_, err = client.NamespaceUpdate(params)
		if err != nil {
			return err
		}
		success(name)
	}
	return nil
}

func mergeNamespaceUpdate(flags *pflag.FlagSet) func(*types.Namespace) error {
	return func(namespace *types.Namespace) error {
		if flags.Changed(flagDisplayName) {
			str, err := flags.GetString(flagDisplayName)
			if err != nil {
				return err
			}
			namespace.DisplayName = str
		}
		if flags.Changed(flagDescription) {
			str, err := flags.GetString(flagDescription)
			if err != nil {
				return err
			}
			namespace.Description = str
		}
		if namespace.Labels == nil {
			namespace.Labels = make(map[string]string)
		}
		if flags.Changed(flagLabelAdd) {
			labels := flags.Lookup(flagLabelAdd).Value.(*opts.ListOpts).GetAll()
			for k, v := range opts.ConvertKVStringsToMap(labels) {
				namespace.Labels[k] = v
			}
		}
		if flags.Changed(flagLabelRemove) {
			keys := flags.Lookup(flagLabelRemove).Value.(*opts.ListOpts).GetAll()
			for _, k := range keys {
				// if a key doesn't exist, fail the command explicitly
				if _, exists := namespace.Labels[k]; !exists {
					return fmt.Errorf("key %s doesn't exist in namespace's labels", k)
				}
				delete(namespace.Labels, k)
			}
		}
		return nil
	}
}
