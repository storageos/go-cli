package login

import (
	"errors"
	"github.com/dnephin/cobra"
	"os"

	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/config"
)

type loginOptions struct {
	host     string
	username string
	password string
}

func NewLoginCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	opt := loginOptions{}

	cmd := &cobra.Command{
		Use:   "login [HOST]",
		Short: "Store login credentials for a given storageos host",
		Args:  cli.RequiresMaxArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogin(storageosCli, opt, args)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opt.host, "host", "", "The host to store the credentials for")
	flags.Lookup("host").Hidden = true
	flags.StringVar(&opt.username, "username", "", "The username to use for this host")
	flags.StringVar(&opt.password, "password", "", "The password to use for this host")

	return cmd
}

func getHost(opt loginOptions, args []string) (string, error) {
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
		host = os.Getenv(config.EnvStorageOSHost)
		if host == "" {
			return "", errors.New("No setting found for host")
		}

	}

	return host, nil
}

func runLogin(storageosCli *command.StorageOSCli, opt loginOptions, args []string) error {
	host, err := getHost(opt, args)
	if err != nil {
		return err
	}

	switch {
	case opt.username == "":
		return errors.New("Please provide a --username")

	case opt.password == "":
		return errors.New("Please provide a --password")

	default:
		err := storageosCli.ConfigFile().CredentialsStore.SetCredentials(host, opt.username, opt.password)
		if err != nil {
			return err
		}

		return storageosCli.ConfigFile().Save()
	}
}
