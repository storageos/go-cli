package user

import (
	"errors"
	"fmt"
	"regexp"

	"context"

	"github.com/dnephin/cobra"
	api "github.com/storageos/go-cli/api"
	"github.com/storageos/go-cli/api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
)

func newRemoveCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rm [OPTIONS] USER [USER...]",
		Aliases: []string{"remove"},
		Short:   "Remove one or more user(s)",
		Args:    cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !verifyUsers(args) {
				fmt.Fprintln(storageosCli.Err(), "Invalid username")
				return errors.New("Invalid username")
			}

			return runRemove(storageosCli, args)
		},
	}

	return cmd
}

func verifyUsers(users []string) bool {
	unameReg := regexp.MustCompile(`[a-zA-Z0-9]+`)

	for _, v := range users {
		if !(api.IsUUID(v) || unameReg.MatchString(v)) {
			return false
		}
	}
	return true
}

func runRemove(storageosCli *command.StorageOSCli, users []string) error {
	client := storageosCli.Client()
	status := 0

	for _, user := range users {
		params := types.DeleteOptions{
			Name:    user,
			Context: context.Background(),
		}

		if err := client.UserDelete(params); err != nil {
			fmt.Fprintf(storageosCli.Err(), "%s\n", err)
			status = 1
			continue
		}
		fmt.Fprintf(storageosCli.Out(), "%s\n", user)
	}

	if status != 0 {
		return cli.StatusError{StatusCode: status}
	}
	return nil
}
