package node

import (
	"context"

	"github.com/dnephin/cobra"

	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/command/formatter"
	"github.com/storageos/go-cli/pkg/constants"
	cliTypes "github.com/storageos/go-cli/types"
)

type healthOptions struct {
	name    string
	quiet   bool
	format  string
	timeout int
}

func newHealthCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	opt := &healthOptions{}

	cmd := &cobra.Command{
		Use:   "health [OPTIONS] NODE",
		Short: "Display detailed information on a given node",
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opt.name = args[0]
			return runHealth(storageosCli, opt)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opt.quiet, "quiet", "q", false, "Display minimal node health info.  Can be used with format.")
	flags.StringVar(&opt.format, "format", "", "Pretty-print health with formats: table (default), cp, dp or raw.")
	flags.IntVarP(&opt.timeout, "timeout", "t", constants.DefaultAPITimeout, "Timeout in seconds.")

	return cmd
}

func runHealth(storageosCli *command.StorageOSCli, opt *healthOptions) error {

	c, err := storageosCli.Client().Node(opt.name)
	if err != nil {
		return err
	}

	node := &cliTypes.Node{
		ID:               c.ID,
		Name:             c.Name,
		AdvertiseAddress: c.Address,
	}

	if err := UpdateNodeHealth(storageosCli, node); err != nil {
		return err
	}

	format := opt.format
	if len(format) == 0 {
		if len(storageosCli.ConfigFile().NodeHealthFormat) > 0 && !opt.quiet {
			format = storageosCli.ConfigFile().NodeHealthFormat
		} else {
			format = formatter.TableFormatKey
		}
	}

	nodeHealthCtx := formatter.Context{
		Output: storageosCli.Out(),
		Format: formatter.NewNodeHealthFormat(format, opt.quiet),
	}
	return formatter.NodeHealthWrite(nodeHealthCtx, node)
}

// UpdateNodeHealth updates the health status of a given node by querying the
// node endpoints.
func UpdateNodeHealth(storageosCli *command.StorageOSCli, node *cliTypes.Node) error {
	cpHealth, err := storageosCli.Client().CPHealth(context.Background(), node.AdvertiseAddress)
	if err != nil {
		return err
	}
	dpHealth, err := storageosCli.Client().DPHealth(context.Background(), node.AdvertiseAddress)
	if err != nil {
		return err
	}
	node.Health.CP = cpHealth
	node.Health.DP = dpHealth
	return nil
}
