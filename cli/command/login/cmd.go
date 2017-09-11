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
		Use:   "login",
		Short: "Store login credentials for a given storageos host",
		Args:  cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogin(storageosCli, opt)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opt.host, "host", "", "The host to store the credentials for")
	flags.StringVar(&opt.username, "username", "", "The username to use for this host")
	flags.StringVar(&opt.password, "password", "", "The password to use for this host")

	return cmd
}

func getHost(opt loginOptions) (string, error) {
	if opt.host != "" {
		return opt.host, nil
	}

	if host := os.Getenv(config.EnvStorageOSHost); host != "" {
		return host, nil
	}

	return "", errors.New("No setting found for host")
}

func runLogin(storageosCli *command.StorageOSCli, opt loginOptions) error {
	host, err := getHost(opt)
	if err != nil {
		return err
	}

	switch {
	case opt.username == "":
		return errors.New("Please provide a --username")

	case opt.password == "":
		return errors.New("Please provide a --password")

	default:
		conf := storageosCli.ConfigFile()

		err := conf.CredentialsStore.SetCredentials(host, opt.username, opt.password)
		if err != nil {
			return err
		}

		return conf.Save()
	}
}
