package logs

import (
	"github.com/dnephin/cobra"

	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/pkg/constants"
)

type logOptions struct {
	nodes       []string
	level       string
	filter      string
	clearFilter bool
	follow      bool
	timeout     int
	quiet       bool
	format      string
	selector    string
}

// NewLogsCommand returns a cobra command for `logs` subcommands
func NewLogsCommand(storageosCli *command.StorageOSCli) *cobra.Command {

	opt := logOptions{}

	cmd := &cobra.Command{
		Use:   "logs",
		Short: "View and manage node logs on the active cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			opt.nodes = args
			if opt.follow && opt.level == "" && opt.filter == "" && !opt.clearFilter {
				return runFollow(storageosCli, opt)
			}
			if !opt.follow && (opt.level != "" || opt.filter != "" || opt.clearFilter) {
				return runUpdate(storageosCli, opt)
			}
			return storageosCli.ShowHelp(cmd, args)
		},
	}
	cmd.AddCommand(
		newViewCommand(storageosCli),
	)

	flags := cmd.Flags()
	flags.StringVarP(&opt.level, "log-level", "l", "", "Set the logging level (\"debug\"|\"info\"|\"warn\"|\"error\"|\"fatal\")")
	flags.StringVar(&opt.filter, "filter", "", "Set the logging filter")
	flags.BoolVarP(&opt.clearFilter, "clear-filter", "", false, "Clears the filter")
	flags.BoolVarP(&opt.follow, "follow", "f", false, "Tail the logs for the given node, or all nodes if not specified")
	flags.IntVarP(&opt.timeout, "timeout", "t", constants.DefaultAPITimeout, "Timeout in seconds.")
	flags.BoolVarP(&opt.quiet, "quiet", "q", false, "Only display volume names")
	flags.StringVar(&opt.format, "format", "raw", "Output format (raw or table) or a Go template (type --format -h or --help for a detail usage)")

	return cmd
}
