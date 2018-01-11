package logs

import (
	"github.com/dnephin/cobra"

	"github.com/storageos/go-cli/cli/command"
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
		Short: "View and manage node logs",
		RunE: func(cmd *cobra.Command, args []string) error {
			opt.nodes = args
			if opt.follow && opt.level == "" && opt.filter == "" {
				return runFollow(storageosCli, opt)
			}
			if !opt.follow && (opt.level != "" || opt.filter != "") {
				return runUpdate(storageosCli, opt)
			}
			return storageosCli.ShowHelp(cmd, args)
		},
	}
	cmd.AddCommand(
		newViewCommand(storageosCli),
	)

	flags := cmd.Flags()
	flags.StringVar(&opt.level, "verbosity", "", "Set the logging verbosity")
	flags.StringVar(&opt.filter, "filter", "", "Set the logging filter")
	flags.BoolVarP(&opt.clearFilter, "clear-filter", "", false, "Clears the filter")
	flags.BoolVarP(&opt.follow, "follow", "f", false, "Tail the logs for the given node, or all nodes if not specified")
	flags.IntVarP(&opt.timeout, "timeout", "t", 1, "Timeout in seconds.")

	return cmd
}
