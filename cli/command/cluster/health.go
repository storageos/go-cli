package cluster

import (
	"context"
	"sort"
	"time"

	"github.com/dnephin/cobra"

	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/command/formatter"
	"github.com/storageos/go-cli/pkg/constants"
	cliTypes "github.com/storageos/go-cli/types"
)

type healthOpt struct {
	quiet   bool
	format  string
	timeout int
}

func newHealthCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	opt := &healthOpt{}

	cmd := &cobra.Command{
		Use:   "health",
		Short: `Displays the cluster's health.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runHealth(storageosCli, opt)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opt.quiet, "quiet", "q", false, "Display minimal cluster health info.  Can be used with format.")
	flags.IntVarP(&opt.timeout, "timeout", "t", constants.DefaultAPITimeout, "Timeout in seconds.")
	flags.StringVar(&opt.format, "format", "", "Pretty-print health with formats: table (default), detailed, cp, dp or raw.")

	return cmd
}

func runHealth(storageosCli *command.StorageOSCli, opt *healthOpt) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(opt.timeout))
	defer cancel()

	status, err := storageosCli.Client().ClusterHealth(ctx)
	if err != nil {
		return err
	}

	sort.Slice(status, func(i, j int) bool {
		return cliTypes.HumanisedStringLess(status[i].NodeName, status[j].NodeName)
	})

	format := opt.format
	if len(format) == 0 {
		if len(storageosCli.ConfigFile().ClusterHealthFormat) > 0 && !opt.quiet {
			format = storageosCli.ConfigFile().ClusterHealthFormat
		} else {
			format = formatter.TableFormatKey
		}
	}

	clusterHealthCtx := formatter.Context{
		Output: storageosCli.Out(),
		Format: formatter.NewClusterHealthFormat(format, opt.quiet),
	}
	return formatter.ClusterHealthWrite(clusterHealthCtx, status)
}
