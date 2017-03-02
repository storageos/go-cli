package namespace

import (
	"sort"

	"github.com/dnephin/cobra"
	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/command/formatter"
	"github.com/storageos/go-cli/cli/opts"
)

type byNamespaceName []*types.Namespace

func (r byNamespaceName) Len() int      { return len(r) }
func (r byNamespaceName) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r byNamespaceName) Less(i, j int) bool {
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
		Short:   "List namespaces",
		Args:    cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(storageosCli, opts)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opts.quiet, "quiet", "q", false, "Only display namespace names")
	flags.StringVar(&opts.format, "format", "", "Pretty-print namespaces using a Go template")
	flags.VarP(&opts.filter, "filter", "f", "Provide filter values (e.g. 'dangling=true')")

	return cmd
}

func runList(storageosCli *command.StorageOSCli, opts listOptions) error {
	client := storageosCli.Client()

	params := types.ListOptions{
	// LabelSelector: opts.filter.Value(),
	}

	// namespaces, err := client.NamespaceList(context.Background(), opts.filter.Value())
	namespaces, err := client.NamespaceList(params)
	if err != nil {
		return err
	}

	format := opts.format
	if len(format) == 0 {
		if len(storageosCli.ConfigFile().NamespacesFormat) > 0 && !opts.quiet {
			format = storageosCli.ConfigFile().NamespacesFormat
		} else {
			format = formatter.TableFormatKey
		}
	}

	sort.Sort(byNamespaceName(namespaces))

	namespaceCtx := formatter.Context{
		Output: storageosCli.Out(),
		Format: formatter.NewNamespaceFormat(format, opts.quiet),
	}
	return formatter.NamespaceWrite(namespaceCtx, namespaces)
}
