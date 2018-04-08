package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/dnephin/cobra"

	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
)

type updateOptions struct {
	sourceAccount string
	password      bool
	groups        stringSlice
	addGroups     stringSlice
	removeGroups  stringSlice
	role          string
}

// processGroups is destructive to given slice
func (u updateOptions) processGroups(current []string) []string {
	if u.groups != nil {
		return u.groups
	}

	needsRemoval := func(s string) bool {
		for _, v := range u.removeGroups {
			if s == v {
				return true
			}
		}
		return false
	}

	newGroups := current[:0]

	// remove groups
	for _, v := range current {
		if !needsRemoval(v) && v != "" {
			newGroups = append(newGroups, v)
		}
	}

	// Checks if a given group exists in the newGroups.
	containsGroup := func(s string) bool {
		for _, v := range newGroups {
			if s == v {
				return true
			}
		}
		return false
	}

	// add groups
	for _, v := range u.addGroups {
		if !containsGroup(v) {
			newGroups = append(newGroups, v)
		}
	}

	return newGroups
}

func newUpdateCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	opt := updateOptions{}

	cmd := &cobra.Command{
		Use:   "update [OPTIONS] USER",
		Short: "Update select fields in a user account",
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opt.sourceAccount = args[0]
			return runUpdate(storageosCli, opt)
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&opt.password, "password", false, "Prompt for new password (interactive)")
	flags.StringVar(&opt.role, "role", "", "Provide a new role")
	flags.Var(&opt.groups, "groups", "Provide a new set of groups (replacing old set)")
	flags.Var(&opt.addGroups, "add-groups", "Add the user to the following groups")
	flags.Var(&opt.removeGroups, "remove-groups", "Remove the user from the following groups")
	return cmd
}

func verifyGroupLogic(opt updateOptions) error {
	if (len(opt.groups) > 0) && (len(opt.addGroups)+len(opt.removeGroups)) > 0 {
		return errors.New("Cannot set both groups and add/remove groups")
	}

	// Check if a group is in both add and remove.
	if (len(opt.addGroups) > 0) && (len(opt.removeGroups) > 0) {
		for _, i := range opt.addGroups {
			for _, j := range opt.removeGroups {
				if i == j {
					return errors.New("Cannot add and remove the same group at a time")
				}
			}
		}
	}
	return nil
}

func verifyUpdate(opt updateOptions) error {
	if i, pass := verifyGroups(opt.groups); !pass {
		return fmt.Errorf(`Group element %d doesn't follow format "[a-zA-Z0-9]+"`, i)
	}

	if i, pass := verifyGroups(opt.addGroups); !pass {
		return fmt.Errorf(`add-group element %d doesn't follow format "[a-zA-Z0-9]+"`, i)
	}

	if i, pass := verifyGroups(opt.removeGroups); !pass {
		return fmt.Errorf(`remove-group element %d doesn't follow format "[a-zA-Z0-9]+"`, i)
	}

	if !(opt.role == "" || verifyRole(opt.role)) {
		return fmt.Errorf(`Role must be either "user" or "admin", not %q`, opt.role)
	}

	return nil
}

func runUpdate(storageosCli *command.StorageOSCli, opt updateOptions) error {
	var password string

	if opt.password {
		var err error

		password, err = getPassword(storageosCli)
		if err != nil {
			return err
		}
	}

	if err := verifyGroupLogic(opt); err != nil {
		return err
	}

	if err := verifyUpdate(opt); err != nil {
		return err
	}

	client := storageosCli.Client()

	currentState, err := client.User(opt.sourceAccount)
	if err != nil {
		return fmt.Errorf("Failed to get user (%s): %s", opt.sourceAccount, err)
	}
	currentState.Groups = opt.processGroups(currentState.Groups)

	if opt.password {
		currentState.Password = password
	}

	if opt.role != "" {
		currentState.Role = opt.role
	}

	return client.UserUpdate(currentState, context.Background())
}
