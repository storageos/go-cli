package logout

import (
	"errors"
	"fmt"
	"github.com/dnephin/cobra"
	"net"
	"os"
	"regexp"
	"strings"

	api "github.com/storageos/go-api"
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

func formatHost(host string) (string, error) {
	validHostname := regexp.MustCompile(`^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$`)

	switch strings.Count(host, ":") {
	// Add the default port if missing then continue
	case 0:
		host += ":" + api.DefaultPort
		fallthrough

	// Validate the host section
	case 1:
		s := strings.Split(host, ":")

		// if not an ip and not a hostname, return error
		if net.ParseIP(s[0]) == nil && !validHostname.MatchString(s[0]) {
			return "", fmt.Errorf("Invalid value for host (%v)\nValue must be in the format 'HOST' or 'HOST:PORT'\n\teg. 'localhost'\n\teg. '10.1.5.249:5705'", host)
		}

		return host, nil

	// Unrecognised format
	default:
		return "", fmt.Errorf("Invalid value for host (%v)\nValue must be in the format 'HOST' or 'HOST:PORT'\n\teg. 'localhost'\n\teg. '10.1.5.249:5705'", host)
	}
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
		host = os.Getenv(config.EnvStorageOSHost)
		if host == "" {
			return "", errors.New("No setting found for host")
		}

	}

	return formatHost(host)
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
