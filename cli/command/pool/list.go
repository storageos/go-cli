package pool

import (
	"sort"

	"github.com/dnephin/cobra"
	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/command/formatter"
	"github.com/storageos/go-cli/cli/opts"
)

type byPoolName []*types.Pool

func (r byPoolName) Len() int      { return len(r) }
func (r byPoolName) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r byPoolName) Less(i, j int) bool {
	return r[i].Name < r[j].Name
}

type listOptions struct {
	quiet  bool
	format string
	filter opts.FilterOpt
}

func newListCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	opts := listOptions{filter: opts.NewFilterOpt()}

	cmd := &cobra.Command{
		Use:     "ls [OPTIONS]",
		Aliases: []string{"list"},
		Short:   "List capacity pools",
		Args:    cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(storageosCli, opts)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opts.quiet, "quiet", "q", false, "Only display pool names")
	flags.StringVar(&opts.format, "format", "", "Pretty-print pools using a Go template")
	flags.VarP(&opts.filter, "filter", "f", "Provide filter values (e.g. 'dangling=true')")

	return cmd
}

func runList(storageosCli *command.StorageOSCli, opts listOptions) error {
	client := storageosCli.Client()

	params := types.ListOptions{
	// LabelSelector: opts.filter.Value(),
	}

	// pools, err := client.PoolList(context.Background(), opts.filter.Value())
	pools, err := client.PoolList(params)
	if err != nil {
		return err
	}

	format := opts.format
	if len(format) == 0 {
		if len(storageosCli.ConfigFile().PoolsFormat) > 0 && !opts.quiet {
			format = storageosCli.ConfigFile().PoolsFormat
		} else {
			format = formatter.TableFormatKey
		}
	}

	sort.Sort(byPoolName(pools))

	poolCtx := formatter.Context{
		Output: storageosCli.Out(),
		Format: formatter.NewPoolFormat(format, opts.quiet),
	}
	return formatter.PoolWrite(poolCtx, pools)
}
