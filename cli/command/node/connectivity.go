package node

import (
	"context"
	"fmt"
	"io"

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

			if len(opt.names) == 0 {
				names, err := allNodeNames(storageosCli)
				if err != nil {
					return err
				}
				opt.names = names
			}
			return runConnectivity(storageosCli, opt)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opt.quiet, "quiet", "q", false, "Only display node names")
	flags.StringVarP(&opt.format, "format", "f", "table", "Format the output using the given Go template")
	return cmd
}

func allNodeNames(storageosCli *command.StorageOSCli) ([]string, error) {
	nodes, err := storageosCli.Client().NodeList(types.ListOptions{
		Context: context.Background(),
	})

	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(nodes))
	for _, n := range nodes {
		names = append(names, n.Name)
	}

	return names, nil
}

type result struct {
	initiator string
	result    []types.ConnectivityResult
}

func runConnectivity(storageosCli *command.StorageOSCli, opt connectivityOptions) error {
	client := storageosCli.Client()

	res := make([]result, 0, len(opt.names))
	for _, ref := range opt.names {
		c, err := client.Connectivity(ref)
		if err != nil {
			return err
		}

		res = append(res, result{ref, c})
	}

	return printConnectivityResult(storageosCli.Out(), res, opt)
}

func printConnectivityResult(out io.Writer, results []result, opt connectivityOptions) error {
	for i, result := range results {
		if len(results) > 1 {
			if i > 0 {
				fmt.Fprintf(out, "\n\n")
			}
			fmt.Fprintf(out, "Connectivity of %s:\n", result.initiator)
		}

		fmtCtx := formatter.Context{
			Output: out,
			Format: formatter.NewConnectivityFormat(opt.format, opt.quiet),
		}

		if result.result != nil {
			if err := formatter.ConnectivityWrite(fmtCtx, result.result); err != nil {
				return err
			}
		}
	}
	return nil
}
