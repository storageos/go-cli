package volume

import (
	"sort"

	"github.com/dnephin/cobra"
	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/command/formatter"
)

type byVolumeName []*types.Volume

func (r byVolumeName) Len() int      { return len(r) }
func (r byVolumeName) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r byVolumeName) Less(i, j int) bool {
	return r[i].Name < r[j].Name
}

type listOptions struct {
	quiet     bool
	format    string
	selector  string
	namespace string
}

func newListCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	opt := listOptions{}

	cmd := &cobra.Command{
		Use:     "ls [OPTIONS]",
		Aliases: []string{"list"},
		Short:   "List volumes",
		Args:    cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(storageosCli, opt)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opt.quiet, "quiet", "q", false, "Only display volume names")
	flags.StringVar(&opt.format, "format", "", "Pretty-print volumes using a Go template")
	flags.StringVarP(&opt.selector, "selector", "s", "", "Provide selector (e.g. to list all volumes with label app=cassandra ' --selector=app=cassandra')")
	flags.StringVarP(&opt.namespace, "namespace", "n", "", "Namespace scope")

	return cmd
}

func runList(storageosCli *command.StorageOSCli, opt listOptions) error {
	client := storageosCli.Client()

	params := types.ListOptions{
		LabelSelector: opt.selector,
		Namespace:     opt.namespace,
	}

	// volumes, err := client.VolumeList(context.Background(), opts.filter.Value())
	volumes, err := client.VolumeList(params)
	if err != nil {
		return err
	}

	nodes, err := client.ControllerList(params)
	if err != nil {
		return err
	}

	format := opt.format
	if len(format) == 0 {
		if len(storageosCli.ConfigFile().VolumesFormat) > 0 && !opt.quiet {
			format = storageosCli.ConfigFile().VolumesFormat
		} else {
			format = formatter.TableFormatKey
		}
	}

	sort.Sort(byVolumeName(volumes))

	volumeCtx := formatter.Context{
		Output: storageosCli.Out(),
		Format: formatter.NewVolumeFormat(format, opt.quiet),
	}
	return formatter.VolumeWrite(volumeCtx, volumes, nodes)
}
