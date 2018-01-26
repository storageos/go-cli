package namespace

import (
	"sort"

	"github.com/dnephin/cobra"
	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/command/formatter"
)

type byNamespaceName []*types.Namespace

func (r byNamespaceName) Len() int      { return len(r) }
func (r byNamespaceName) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r byNamespaceName) Less(i, j int) bool {
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
		Short:   "List namespaces",
		Args:    cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(storageosCli, opt)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opt.quiet, "quiet", "q", false, "Only display namespace names")
	flags.StringVar(&opt.format, "format", "", "Pretty-print namespaces using a Go template")
	flags.StringVarP(&opt.selector, "selector", "s", "", "Provide selector (e.g. to list all namespaces with label app=cassandra ' --selector=app=cassandra')")

	return cmd
}

func runList(storageosCli *command.StorageOSCli, opt listOptions) error {
	client := storageosCli.Client()

	params := types.ListOptions{
		LabelSelector: opt.selector,
	}

	// namespaces, err := client.NamespaceList(context.Background(), opt.filter.Value())
	namespaces, err := client.NamespaceList(params)
	if err != nil {
		return err
	}

	format := opt.format
	if len(format) == 0 {
		if len(storageosCli.ConfigFile().NamespacesFormat) > 0 && !opt.quiet {
			format = storageosCli.ConfigFile().NamespacesFormat
		} else {
			format = formatter.TableFormatKey
		}
	}

	sort.Sort(byNamespaceName(namespaces))

	namespaceCtx := formatter.Context{
		Output: storageosCli.Out(),
		Format: formatter.NewNamespaceFormat(format, opt.quiet),
	}
	return formatter.NamespaceWrite(namespaceCtx, namespaces)
}
