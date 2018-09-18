package user

import (
	"errors"
	"fmt"
	"sort"

	"github.com/dnephin/cobra"
	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/command/formatter"
)

type byUserName []*types.User

func (r byUserName) Len() int      { return len(r) }
func (r byUserName) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r byUserName) Less(i, j int) bool {
	return r[i].Username < r[j].Username
}

func filter(users []*types.User, by func(*types.User) bool) []*types.User {
	rtn := make([]*types.User, 0)

	for _, v := range users {
		if by(v) {
			rtn = append(rtn, v)
		}
	}

	return rtn
}

type listOptions struct {
	quiet  bool
	format string
	admins bool
	users  bool
}

func newListCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	opt := listOptions{}

	cmd := &cobra.Command{
		Use:     "ls [OPTIONS]",
		Aliases: []string{"list"},
		Short:   "List users",
		Args:    cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if opt.admins && opt.users {
				fmt.Fprintln(storageosCli.Err(), "cannot return only admins and only users")
				return errors.New("cannot return only admins and only users")
			}

			return runList(storageosCli, opt)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opt.quiet, "quiet", "q", false, "Only display usernames")
	flags.StringVar(&opt.format, "format", "", "Pretty-print rules using a Go template")
	flags.BoolVar(&opt.admins, "admin-only", false, "Only return the admin users")
	flags.BoolVar(&opt.users, "user-only", false, "Only return the non-admin users")

	return cmd
}

func runList(storageosCli *command.StorageOSCli, opt listOptions) error {
	client := storageosCli.Client()

	params := types.ListOptions{}

	users, err := client.UserList(params)
	if err != nil {
		return err
	}

	format := opt.format
	if len(format) == 0 {
		if len(storageosCli.ConfigFile().UsersFormat) > 0 && !opt.quiet {
			format = storageosCli.ConfigFile().UsersFormat
		} else {
			format = formatter.TableFormatKey
		}
	}

	sort.Sort(byUserName(users))
	if opt.admins {
		users = filter(users, func(u *types.User) bool { return u.Role == "admin" })
	}
	if opt.users {
		users = filter(users, func(u *types.User) bool { return u.Role == "user" })
	}

	userCtx := formatter.Context{
		Output: storageosCli.Out(),
		Format: formatter.NewUserFormat(format, opt.quiet),
	}
	return formatter.UserWrite(userCtx, users)
}
