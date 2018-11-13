package pool

import (
	"sort"

	"github.com/spf13/cobra"
	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/command/formatter"
	"github.com/storageos/go-cli/pkg/constants"
)

type byPoolName []*types.Pool

func (r byPoolName) Len() int      { return len(r) }
func (r byPoolName) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r byPoolName) Less(i, j int) bool {
	return r[i].Name < r[j].Name
}

type listOptions struct {
	quiet    bool
	format   string
	selector string
}

func newListCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	opt := listOptions{}

	cmd := &cobra.Command{
		Use:     "ls [OPTIONS]",
		Aliases: []string{"list"},
		Short:   "List capacity pools",
		Args:    cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(storageosCli, opt)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opt.quiet, "quiet", "q", false, "Only display pool names")
	flags.StringVar(&opt.format, "format", "", "Pretty-print pools using a Go template"+constants.PoolContextTemplate)
	flags.StringVarP(&opt.selector, "selector", "s", "", "Provide selector (e.g. to list all pools with label app=cassandra ' --selector=app=cassandra')")

	return cmd
}

func runList(storageosCli *command.StorageOSCli, opt listOptions) error {
	client := storageosCli.Client()

	params := types.ListOptions{
		LabelSelector: opt.selector,
	}

	// pools, err := client.PoolList(context.Background(), opt.filter.Value())
	pools, err := client.PoolList(params)
	if err != nil {
		return err
	}

	format := opt.format
	if len(format) == 0 {
		if len(storageosCli.ConfigFile().PoolsFormat) > 0 && !opt.quiet {
			format = storageosCli.ConfigFile().PoolsFormat
		} else {
			format = formatter.TableFormatKey
		}
	}

	sort.Sort(byPoolName(pools))

	poolCtx := formatter.Context{
		Output: storageosCli.Out(),
		Format: formatter.NewPoolFormat(format, opt.quiet),
	}
	return formatter.PoolWrite(poolCtx, pools)
}
