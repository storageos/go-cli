package rule

import (
	"fmt"

	"github.com/dnephin/cobra"
	"github.com/spf13/pflag"
	storageos "github.com/storageos/go-api"
	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/opts"
	"golang.org/x/net/context"
)

const (
	flagDescription    = "description"
	flagActive         = "active"
	flagWeight         = "weight"
	flagOperator       = "operator"
	flagRuleAction     = "action"
	flagSelectorAdd    = "selector-add"
	flagSelectorRemove = "selector-rm"
	flagLabelAdd       = "label-add"
	flagLabelRemove    = "label-rm"
)

type updateOptions struct {
	name        string
	description string
	active      bool
	weight      int
	operator    string
	ruleAction  string
	selectors   opts.ListOpts
	labels      opts.ListOpts
}

func newUpdateCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	opt := updateOptions{
		selectors: opts.NewListOpts(opts.ValidateEnv),
		labels:    opts.NewListOpts(opts.ValidateEnv),
	}

	cmd := &cobra.Command{
		Use:   "update [OPTIONS] RULE",
		Short: "Update a rule",
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpdate(storageosCli, cmd.Flags(), args[0])
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opt.description, flagDescription, "d", "", `Rule description`)
	flags.StringVarP(&opt.ruleAction, flagRuleAction, "a", "add", "Rule action (add|remove)")
	flags.StringVarP(&opt.operator, flagOperator, "o", "==", "Comparison operator (!|=|==|in|!=|notin|exists|gt|lt)")
	flags.Var(&opt.selectors, flagSelectorAdd, "Add or update a rule selector (key=value)")
	selectorKeys := opts.NewListOpts(nil)
	flags.Var(&selectorKeys, flagSelectorRemove, "Remove a rule selector if exists")
	flags.IntVarP(&opt.weight, flagWeight, "w", 5, "Rule weight determines processing order (0-10)")
	flags.BoolVar(&opt.active, flagActive, true, "Enable or disable the pool")
	flags.Var(&opt.labels, flagLabelAdd, "Add or update a label (key=value)")
	labelKeys := opts.NewListOpts(nil)
	flags.Var(&labelKeys, flagLabelRemove, "Remove a label if exists")
	return cmd
}

func runUpdate(storageosCli *command.StorageOSCli, flags *pflag.FlagSet, ref string) error {
	success := func(_ string) {
		fmt.Fprintln(storageosCli.Out(), ref)
	}
	return updateRules(storageosCli, []string{ref}, mergeRuleUpdate(flags), success)
}

func updateRules(storageosCli *command.StorageOSCli, refs []string, mergeRule func(rule *types.Rule) error, success func(name string)) error {
	client := storageosCli.Client()
	ctx := context.Background()

	for _, ref := range refs {

		namespace, name, err := storageos.ParseRef(ref)
		if err != nil {
			return err
		}

		rule, err := client.Rule(namespace, name)
		if err != nil {
			return err
		}

		err = mergeRule(rule)
		if err != nil {
			return err
		}
		params := types.RuleUpdateOptions{
			Name:        rule.Name,
			Namespace:   rule.Namespace,
			Description: rule.Description,
			RuleAction:  rule.RuleAction,
			Operator:    rule.Operator,
			Selectors:   rule.Selectors,
			Active:      rule.Active,
			Weight:      rule.Weight,
			Labels:      rule.Labels,
			Context:     ctx,
		}
		_, err = client.RuleUpdate(params)
		if err != nil {
			return err
		}
		success(name)
	}
	return nil
}

func mergeRuleUpdate(flags *pflag.FlagSet) func(*types.Rule) error {
	return func(rule *types.Rule) error {
		if flags.Changed(flagDescription) {
			str, err := flags.GetString(flagDescription)
			if err != nil {
				return err
			}
			rule.Description = str
		}
		if flags.Changed(flagRuleAction) {
			str, err := flags.GetString(flagRuleAction)
			if err != nil {
				return err
			}
			rule.RuleAction = str
		}
		if flags.Changed(flagOperator) {
			str, err := flags.GetString(flagOperator)
			if err != nil {
				return err
			}
			rule.Operator = str
		}
		if rule.Selectors == nil {
			rule.Selectors = make(map[string]string)
		}
		if flags.Changed(flagSelectorAdd) {
			selectors := flags.Lookup(flagSelectorAdd).Value.(*opts.ListOpts).GetAll()
			for k, v := range opts.ConvertKVStringsToMap(selectors) {
				rule.Selectors[k] = v
			}
		}
		if flags.Changed(flagSelectorRemove) {
			keys := flags.Lookup(flagSelectorRemove).Value.(*opts.ListOpts).GetAll()
			for _, k := range keys {
				// if a key doesn't exist, fail the command explicitly
				if _, exists := rule.Selectors[k]; !exists {
					return fmt.Errorf("key %s doesn't exist in rule's selectors", k)
				}
				delete(rule.Selectors, k)
			}
		}
		if flags.Changed(flagActive) {
			active, err := flags.GetBool(flagActive)
			if err != nil {
				return err
			}
			rule.Active = active
		}
		if flags.Changed(flagWeight) {
			weight, err := flags.GetInt(flagWeight)
			if err != nil {
				return err
			}
			rule.Weight = weight
		}
		if rule.Labels == nil {
			rule.Labels = make(map[string]string)
		}
		if flags.Changed(flagLabelAdd) {
			labels := flags.Lookup(flagLabelAdd).Value.(*opts.ListOpts).GetAll()
			for k, v := range opts.ConvertKVStringsToMap(labels) {
				rule.Labels[k] = v
			}
		}
		if flags.Changed(flagLabelRemove) {
			keys := flags.Lookup(flagLabelRemove).Value.(*opts.ListOpts).GetAll()
			for _, k := range keys {
				// if a key doesn't exist, fail the command explicitly
				if _, exists := rule.Labels[k]; !exists {
					return fmt.Errorf("key %s doesn't exist in rule's labels", k)
				}
				delete(rule.Labels, k)
			}
		}
		return nil
	}
}
