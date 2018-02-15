package node

import (
	"context"
	"time"

	"github.com/dnephin/cobra"

	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/command/formatter"
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
	flags.IntVarP(&opt.timeout, "timeout", "t", 1, "Timeout in seconds.")

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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(opt.timeout))
	defer cancel()

	// Ignore errors and carry on
	cpHealth, _ := storageosCli.Client().CPHealth(ctx, node.AdvertiseAddress)
	node.Health.CP = cpHealth

	dpHealth, err := storageosCli.Client().DPHealth(ctx, node.AdvertiseAddress)
	node.Health.DP = dpHealth

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
