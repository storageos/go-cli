package rule

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
	namespace   string
	description string
	active      bool
	weight      int
	operator    string
	ruleAction  string
	selectors   opts.ListOpts
	labels      opts.ListOpts
}

func newCreateCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	opt := createOptions{
		selectors: opts.NewListOpts(opts.ValidateEnv),
		labels:    opts.NewListOpts(opts.ValidateEnv),
	}

	cmd := &cobra.Command{
		Use:   "create [OPTIONS] [RULE]",
		Short: "Create a rule",
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
	flags.StringVar(&opt.name, "name", "", "Rule name")
	flags.Lookup("name").Hidden = true
	flags.StringVarP(&opt.description, "description", "d", "", "Rule description")
	flags.StringVarP(&opt.ruleAction, "action", "a", "add", "Rule action (add|remove)")
	flags.StringVarP(&opt.operator, "operator", "o", "==", "Comparison operator (!|=|==|in|!=|notin|exists|gt|lt)")
	flags.VarP(&opt.selectors, "selector", "s", "key=value selectors to trigger rule")
	flags.IntVarP(&opt.weight, "weight", "w", 5, "Rule weight determines processing order (0-10)")
	flags.StringVarP(&opt.namespace, "namespace", "n", "", "Volume namespace")
	flags.BoolVar(&opt.active, "active", true, "Enable or disable the rule")

	flags.Var(&opt.labels, "label", "Labels to apply when rule is triggered")

	return cmd
}

func runCreate(storageosCli *command.StorageOSCli, opt createOptions) error {
	if _, err := opts.ValidateOperator(opt.operator); err != nil {
		return err
	}
	if _, err := opts.ValidateRuleAction(opt.ruleAction); err != nil {
		return err
	}

	client := storageosCli.Client()

	params := types.RuleCreateOptions{
		Name:        opt.name,
		Namespace:   opt.namespace,
		Description: opt.description,
		RuleAction:  opt.ruleAction,
		Operator:    opt.operator,
		Selectors:   opts.ConvertKVStringsToMap(opt.selectors.GetAll()),
		Active:      opt.active,
		Weight:      opt.weight,
		Labels:      opts.ConvertKVStringsToMap(opt.labels.GetAll()),
		Context:     context.Background(),
	}

	rule, err := client.RuleCreate(params)
	if err != nil {
		return err
	}

	fmt.Fprintf(storageosCli.Out(), "%s/%s\n", rule.Namespace, rule.Name)
	return nil
}
