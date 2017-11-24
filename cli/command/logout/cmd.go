package logout

import (
	"errors"
	"github.com/dnephin/cobra"
	"os"

	api "github.com/storageos/go-api"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/config"
	"github.com/storageos/go-cli/pkg/validation"
)

type logoutOptions struct {
	host string
}

func NewLogoutCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	opt := logoutOptions{}

	cmd := &cobra.Command{
		Use:   "logout [HOST]",
		Short: "Delete stored login credentials for a given storageos host",
		Args:  cli.RequiresMaxArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDelete(storageosCli, opt, args)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opt.host, "host", "", "The host to remove the credentials for")
	flags.Lookup("host").Hidden = true

	return cmd
}

func getHost(opt logoutOptions, args []string) (string, error) {
	var host string

	switch {
	case len(args) == 1:
		if opt.host != "" {
			return "", errors.New("Conflicting options: either specify --host or provide positional arg, not both")
		}
		host = args[0]

	case opt.host != "":
		host = opt.host

	default:
		return validation.ParseHostPort(api.DefaultHost, api.DefaultPort)

	}

	return validation.ParseHostPort(host, api.DefaultPort)
}

func runDelete(storageosCli *command.StorageOSCli, opt logoutOptions, args []string) error {
	host, err := getHost(opt, args)
	if err != nil {
		return err
	}

	conf := storageosCli.ConfigFile()

	conf.CredentialsStore.DeleteCredentials(host)
	return conf.Save()

}
