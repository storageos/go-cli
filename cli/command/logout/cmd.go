package logout

import (
	"errors"
	"fmt"

	"github.com/dnephin/cobra"

	api "github.com/storageos/go-api"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/pkg/jointools"
)

type logoutOptions struct {
	host string
}

// NewLogoutCommand returns the Cobra command for logout
func NewLogoutCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	opt := logoutOptions{}

	cmd := &cobra.Command{
		Use:   "logout HOST",
		Short: "Delete stored login credentials for a given storageos host",
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDelete(storageosCli, opt, args)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opt.host, "host", "", "The host to remove the credentials for")
	flags.Lookup("host").Hidden = true

	return cmd
}

func getHost(discoveryHost string, opt logoutOptions, args []string) (string, error) {
	var join string

	switch {
	case len(args) == 1 && opt.host != "":
		return "", errors.New("Conflicting options: either specify --host or provide positional arg, not both")

	case len(args) == 1:
		join = args[0]

	case opt.host != "":
		join = opt.host

	default:
		join = api.DefaultHost
	}

	if errs := jointools.VerifyJOIN(discoveryHost, join); errs != nil {
		return "", fmt.Errorf("error: %+v", errs)
	}
	return jointools.ExpandJOIN(discoveryHost, join), nil

}

func runDelete(storageosCli *command.StorageOSCli, opt logoutOptions, args []string) error {
	host, err := getHost(storageosCli.GetDiscovery(), opt, args)
	if err != nil {
		return err
	}

	conf := storageosCli.ConfigFile()

	conf.CredentialsStore.DeleteCredentials(host)
	return conf.Save()

}
