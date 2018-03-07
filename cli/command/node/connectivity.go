package node

import (
	"fmt"
	"github.com/dnephin/cobra"
	"io"
	// storageos "github.com/storageos/go-api"
	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/command/formatter"
)

type connectivityOptions struct {
	format string
	names  []string
}

func newConnectivityCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	var opt connectivityOptions

	cmd := &cobra.Command{
		Use:   "connectivity [OPTIONS] NODE [NODE...]",
		Short: "Display detailed connectivity information on one or more nodes",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opt.names = args
			return runConnectivity(storageosCli, opt)
		},
	}

	cmd.Flags().StringVarP(&opt.format, "format", "f", "", "Format the output using the given Go template")
	return cmd
}

func runConnectivity(storageosCli *command.StorageOSCli, opt connectivityOptions) error {
	client := storageosCli.Client()

	result := make(map[string]*types.NodeConnectivity)
	for _, ref := range opt.names {
		c, err := client.Connectivity(ref)
		if err != nil {
			return err
		}

		result[ref] = c
	}

	return printConnectivityResult(storageosCli.Out(), result)
}

func printConnectivityResult(out io.Writer, result map[string]*types.NodeConnectivity) error {
	var printRefHeader = len(result) > 1

	for node, result := range result {
		if printRefHeader {
			fmt.Fprintf(out, "%s:\n", node)
		}

		fmtCtx := formatter.Context{
			Output: out,
			Format: formatter.NewConnectivityFormat(formatter.TableFormatKey),
		}

		if err := formatter.ConnectivityWrite(fmtCtx, result); err != nil {
			return err
		}
	}
	return nil
}
