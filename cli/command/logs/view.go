package logs

import (
	"context"
	"time"

	"github.com/dnephin/cobra"
	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/command/formatter"
	"github.com/storageos/go-cli/pkg/constants"
)

func newViewCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	opt := logOptions{}

	cmd := &cobra.Command{
		Use:   "view [OPTIONS]",
		Short: "Show logging configuration",
		Args:  cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runView(storageosCli, opt)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opt.quiet, "quiet", "q", false, "Only display log level")
	flags.StringVar(&opt.format, "format", "", "Pretty-print config using a Go template")
	flags.StringVarP(&opt.selector, "selector", "s", "", "Provide selector (e.g. to list all nodes with label disk=ssd' --selector=disk=ssd')")
	flags.IntVarP(&opt.timeout, "timeout", "t", constants.DefaultAPITimeout, "Timeout in seconds.")

	return cmd
}

func runView(storageosCli *command.StorageOSCli, opt logOptions) error {
	client := storageosCli.Client()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(opt.timeout))
	defer cancel()

	params := types.ListOptions{
		LabelSelector: opt.selector,
		Context:       ctx,
	}

	configs, err := client.LoggerConfig(params)
	if err != nil {
		return err
	}

	format := opt.format
	if len(format) == 0 {
		format = formatter.TableFormatKey
	}

	fmtCtx := formatter.Context{
		Output: storageosCli.Out(),
		Format: formatter.NewLoggerFormat(format, opt.quiet),
	}

	return formatter.LoggerWrite(fmtCtx, configs)
}
