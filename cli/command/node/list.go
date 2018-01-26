package node

import (
	"sort"

	"github.com/dnephin/cobra"
	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/command/formatter"
)

type byControllerName []*types.Controller

func (r byControllerName) Len() int      { return len(r) }
func (r byControllerName) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r byControllerName) Less(i, j int) bool {
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
		Short:   "List nodes",
		Args:    cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(storageosCli, opt)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opt.quiet, "quiet", "q", false, "Only display node names")
	flags.StringVar(&opt.format, "format", "", "Pretty-print nodes using a Go template")
	flags.StringVarP(&opt.selector, "selector", "s", "", "Provide selector (e.g. to list all nodes with label disk=ssd' --selector=disk=ssd')")

	return cmd
}

func runList(storageosCli *command.StorageOSCli, opt listOptions) error {
	client := storageosCli.Client()

	params := types.ListOptions{
		LabelSelector: opt.selector,
	}

	nodes, err := client.ControllerList(params)
	if err != nil {
		return err
	}

	format := opt.format
	if len(format) == 0 {
		if len(storageosCli.ConfigFile().RulesFormat) > 0 && !opt.quiet {
			format = storageosCli.ConfigFile().RulesFormat
		} else {
			format = formatter.TableFormatKey
		}
	}

	sort.Sort(byControllerName(nodes))

	nodeCtx := formatter.Context{
		Output: storageosCli.Out(),
		Format: formatter.NewNodeFormat(format, opt.quiet),
	}

	return formatter.NodeWrite(nodeCtx, nodes)
}
