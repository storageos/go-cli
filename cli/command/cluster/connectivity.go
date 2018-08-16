package cluster

import (
	"github.com/dnephin/cobra"

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
	flags.BoolVarP(&opt.quiet, "quiet", "q", false, "Only display test and status")
	flags.StringVarP(&opt.format, "format", "f", "table", "Format the output using the given Go template. \"summary\", \"table\" and \"raw\" also supported.")
	return cmd
}

func runConnectivity(storageosCli *command.StorageOSCli, opt connectivityOptions) error {
	client := storageosCli.Client()

	results, err := client.NetworkDiagnostics("")
	if err != nil {
		return err
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
