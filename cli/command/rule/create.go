package rule

import (
	"fmt"

	"context"

	"github.com/dnephin/cobra"
	"github.com/storageos/go-cli/api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/opts"
)

type createOptions struct {
	name        string
	namespace   string
	description string
	active      bool
	weight      int
	ruleAction  string
	selector    string
	labels      opts.ListOpts
}

func newCreateCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	opt := createOptions{
		labels: opts.NewListOpts(opts.ValidateEnv),
	}

	cmd := &cobra.Command{
		Use: "create [OPTIONS] [RULE]",
		Short: `Creates a rule. To create a rule that configures 2 replicas for volumes with the label env=prod, run:
storageos rule create --selector env==prod --action add --label storageos.feature.replicas=2 replicator
		`,
		Args: cli.RequiresMaxArgs(1),
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
	flags.StringVarP(&opt.selector, "selector", "s", "", "selectors to trigger rule, i.e. 'environment = production' (operators !|=|==|in|!=|notin|exists|<|>")
	flags.IntVarP(&opt.weight, "weight", "w", 5, "Rule weight determines processing order (0-10)")
	flags.StringVarP(&opt.namespace, "namespace", "n", "default", "Rule namespace")
	flags.BoolVar(&opt.active, "active", true, "Enable or disable the rule")

	flags.Var(&opt.labels, "label", "Labels to apply when rule is triggered")

	return cmd
}

func runCreate(storageosCli *command.StorageOSCli, opt createOptions) error {
	if _, err := opts.ValidateRuleAction(opt.ruleAction); err != nil {
		return err
	}

	client := storageosCli.Client()

	params := types.RuleCreateOptions{
		Name:        opt.name,
		Namespace:   opt.namespace,
		Description: opt.description,
		RuleAction:  opt.ruleAction,
		Selector:    opt.selector,
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
