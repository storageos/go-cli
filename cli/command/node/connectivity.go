package node

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
	names  []string
}

func newConnectivityCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	var opt connectivityOptions

	cmd := &cobra.Command{
		Use:   "connectivity [OPTIONS] NODE [NODE...]",
		Short: "Display detailed connectivity information on one or more nodes",
		Args:  cli.RequiresMinArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			opt.names = args
			return runConnectivity(storageosCli, opt)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opt.quiet, "quiet", "q", false, "Only display node names")
	flags.StringVarP(&opt.format, "format", "f", "table", "Format the output using the given Go template")
	return cmd
}

func runConnectivity(storageosCli *command.StorageOSCli, opt connectivityOptions) (err error) {

	client := storageosCli.Client()

	var results types.ConnectivityResults
	switch {
	case len(opt.names) == 0:
		results, err = client.NetworkDiagnostics("")
		if err != nil {
			return err
		}
	default:
		for _, ref := range opt.names {
			nodeResults, err := client.NetworkDiagnostics(ref)
			if err != nil {
				return err
			}
			results = append(results, nodeResults...)
		}
	}

	summary := false
	if opt.format == "summary" {
		summary = true
	}

	fmtCtx := formatter.Context{
		Output: storageosCli.Out(),
		Format: formatter.NewConnectivityFormat(opt.format, opt.quiet),
		Trunc:  summary, // Use Trunc to flag that we should summarize results
	}
	return formatter.ConnectivityWrite(fmtCtx, results)
}
