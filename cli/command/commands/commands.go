package commands

import (
	"fmt"
	"github.com/dnephin/cobra"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/command/cluster"
	"github.com/storageos/go-cli/cli/command/login"
	"github.com/storageos/go-cli/cli/command/logout"
	"github.com/storageos/go-cli/cli/command/namespace"
	"github.com/storageos/go-cli/cli/command/node"
	"github.com/storageos/go-cli/cli/command/policy"
	"github.com/storageos/go-cli/cli/command/pool"
	"github.com/storageos/go-cli/cli/command/rule"
	"github.com/storageos/go-cli/cli/command/system"
	"github.com/storageos/go-cli/cli/command/user"
	"github.com/storageos/go-cli/cli/command/volume"
)

// AddCommands adds all the commands from cli/command to the root command
func AddCommands(cmd *cobra.Command, storageosCli *command.StorageOSCli) {
	cmd.AddCommand(
		command.WithAlias(namespace.NewNamespaceCommand(storageosCli), "ns"),
		pool.NewPoolCommand(storageosCli),
		rule.NewRuleCommand(storageosCli),
		command.WithAlias(user.NewUserCommand(storageosCli), "u"),
		command.WithAlias(policy.NewPolicyCommand(storageosCli), "pol"),
		command.WithAlias(volume.NewVolumeCommand(storageosCli), "v", "vol"),
		command.WithAlias(node.NewNodeCommand(storageosCli), "n"),
		login.NewLoginCommand(storageosCli),
		logout.NewLogoutCommand(storageosCli),

		// system
		// system.NewSystemCommand(storageosCli),
		system.NewVersionCommand(storageosCli),

		// clustering
		command.WithAlias(cluster.NewClusterCommand(storageosCli), "c"),

		NewBashGenerationFunction(storageosCli),
	)
}

func NewBashGenerationFunction(storageosCli *command.StorageOSCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "install-bash-completion",
		Short:  "Install the bash completion for the storageos command",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			const bashdir = "/etc/bash_completion.d/storageos"

			// ensure user wants to perform this action
			buf := make([]byte, 1024)
			fmt.Fprintf(storageosCli.Out(), "will write bash completion to %s, continue? [y/N] ", bashdir)
			i, err := storageosCli.In().Read(buf)
			if err != nil {
				return err
			}

			switch string(buf[:i-1]) {
			case "y":
				break // just continue

			case "", "n", "N":
				return nil

			default:
				return fmt.Errorf("unknown response (%s) aborting", string(buf[:i-1]))
			}

			if err := cmd.Parent().GenBashCompletionFile(bashdir); err != nil {
				return err
			}

			fmt.Fprintln(storageosCli.Out(), "saved bash completions, please reload your terminal")
			return nil
		},
	}
	return cmd
}
