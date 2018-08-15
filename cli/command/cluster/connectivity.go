package cluster

import (
	"github.com/dnephin/cobra"

	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/command/formatter"
)

type connectivityOptions struct {
	quiet  bool
	format string
}

func newConnectivityCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	var opt connectivityOptions

	cmd := &cobra.Command{
		Use:   "connectivity [OPTIONS]",
		Short: "Display connectivity diagnostics for the cluster",
		Args:  cli.RequiresMinArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConnectivity(storageosCli, opt)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opt.quiet, "quiet", "q", false, "Only display node names")
	flags.StringVarP(&opt.format, "format", "f", "table", "Format the output using the given Go template")
	return cmd
}

func runConnectivity(storageosCli *command.StorageOSCli, opt connectivityOptions) error {
	client := storageosCli.Client()

	results, err := client.Connectivity("")
	if err != nil {
		return err
	}

	switch opt.quiet {
	case true:
		return formatter.ConnectivityWriteSummary(formatter.Context{
			Output: storageosCli.Out(),
			Format: formatter.NewConnectivityFormat(opt.format, opt.quiet),
		}, isOK(results))
	default:
		fmtCtx := formatter.Context{
			Output: storageosCli.Out(),
			Format: formatter.NewConnectivityFormat(opt.format, opt.quiet),
		}
		return formatter.ConnectivityWrite(fmtCtx, results)
	}

}

// isOK returns false if any connectivity result was not ok.
func isOK(results []types.ConnectivityResult) bool {
	for _, result := range results {
		if !result.Passes() {
			return false
		}
	}
	return true
}
