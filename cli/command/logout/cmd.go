package logout

import (
	"errors"
	"github.com/dnephin/cobra"
	"os"

	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/config"
)

type logoutOptions struct {
	host string
}

func NewLogoutCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	opt := logoutOptions{}

	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Delete stored login credentials for a given storageos host",
		Args:  cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDelete(storageosCli, opt)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opt.host, "host", "", "The host to remove the credentials for")

	return cmd
}

func getHost(opt logoutOptions) (string, error) {
	if opt.host != "" {
		return opt.host, nil
	}

	if host := os.Getenv(config.EnvStorageOSHost); host != "" {
		return host, nil
	}

	return "", errors.New("No setting found for host")
}

func runDelete(storageosCli *command.StorageOSCli, opt logoutOptions) error {
	host, err := getHost(opt)
	if err != nil {
		return err
	}

	conf := storageosCli.ConfigFile()

	conf.CredentialsStore.DeleteCredentials(host)
	return conf.Save()

}
