package rule

import (
	"sort"

	"github.com/spf13/cobra"
	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/command/formatter"
	"github.com/storageos/go-cli/pkg/constants"
)

type byRuleName []*types.Rule

func (r byRuleName) Len() int      { return len(r) }
func (r byRuleName) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r byRuleName) Less(i, j int) bool {
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
		Short:   "List rules",
		Args:    cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(storageosCli, opt)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opt.quiet, "quiet", "q", false, "Only display rule names")
	flags.StringVar(&opt.format, "format", "", "Pretty-print rules using a Go template"+constants.RuleContextTemplate)
	flags.StringVarP(&opt.selector, "selector", "s", "", "Provide selector (e.g. to list all rules with label app=cassandra ' --selector=app=cassandra')")
	flags.StringVarP(&opt.namespace, "namespace", "n", "", "Namespace scope")

	return cmd
}

func runList(storageosCli *command.StorageOSCli, opt listOptions) error {
	client := storageosCli.Client()

	params := types.ListOptions{
		LabelSelector: opt.selector,
		Namespace:     opt.namespace,
	}

	// rules, err := client.RuleList(context.Background(), opt.filter.Value())
	rules, err := client.RuleList(params)
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

	sort.Sort(byRuleName(rules))

	ruleCtx := formatter.Context{
		Output: storageosCli.Out(),
		Format: formatter.NewRuleFormat(format, opt.quiet),
	}
	return formatter.RuleWrite(ruleCtx, rules)
}
