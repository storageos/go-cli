package login

import (
	"errors"
	"fmt"
	"strings"
	"syscall"

	"github.com/dnephin/cobra"

	api "github.com/storageos/go-api"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/pkg/jointools"
	"github.com/storageos/go-cli/version"
	"golang.org/x/crypto/ssh/terminal"
)

type loginOptions struct {
	host     string
	username string
	password string
}

// NewLoginCommand returns a command which, when run, will authenticate the
// provided credentials against the target host(s). Upon successful authentication,
// the details will be stored against the host address(es) in the credentials store.
// The command takes the target host(s) either via flag or positional argument,
// whilst username and password can be given either interactively or by flag.
func NewLoginCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	opt := loginOptions{}

	cmd := &cobra.Command{
		Use:   "login [HOST]",
		Short: "Store login credentials for a given storageos host",
		Args:  cli.RequiresMaxArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			parentCmd := cmd.Parent()

			if parentCmd != nil {
				parentFlags := parentCmd.Flags()
				if username, err := parentFlags.GetString("username"); err == nil && opt.username == "" {
					opt.username = username
				}
				if password, err := parentFlags.GetString("password"); err == nil && opt.password == "" {
					opt.password = password
				}
			}

			return runLogin(storageosCli, opt, args)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opt.host, "host", "H", "", "The host to store the credentials for")
	flags.Lookup("host").Hidden = true
	flags.StringVarP(&opt.username, "username", "u", "", "The username to use for this host (will override value of the global option --username)")
	flags.StringVarP(&opt.password, "password", "p", "", "The password to use for this host (will override value of the global option --password)")

	return cmd
}

func verifyCredsWithServer(username, password, host string) error {
	client, err := api.NewVersionedClient(host, api.DefaultVersionStr)
	if err != nil {
		return fmt.Errorf("Failed to verify credentials (%v)", err)
	}
	client.SetAuth(username, password)

	// Set StorageOS CLI UserAgent for all the API requests.
	client.SetUserAgent(strings.Join([]string{version.UserAgent, version.Version}, "/"))

	_, err = client.Login()
	if err != nil {
		return fmt.Errorf("Failed to verify credentials (%v)", err)
	}
	return nil
}

func getHost(opt loginOptions, args []string) (string, error) {
	var join string

	switch {
	case opt.host != "" && len(args) > 0:
		return "", errors.New("conflicting options: either specify --host or provide positional arg, not both")

	case opt.host != "":
		join = opt.host

	case len(args) > 0:
		join = args[0]

	default:
		join = api.DefaultHost
	}

	if errs := jointools.VerifyJOIN(join); errs != nil {
		return "", fmt.Errorf("error: %+v", errs)
	}
	return jointools.ExpandJOIN(join), nil
}

func promptUsername(storageosCli *command.StorageOSCli) (string, error) {
	buf := make([]byte, 1024)
	fmt.Fprint(storageosCli.Out(), "Username: ")
	i, err := storageosCli.In().Read(buf)
	if err != nil {
		return "", err
	}

	return string(buf[:i-1]), nil // i-1 strips newline
}

func promptPassword(storageosCli *command.StorageOSCli) (string, error) {
	fmt.Fprint(storageosCli.Out(), "Password: ")
	p, err := terminal.ReadPassword(int(syscall.Stdin)) // Cast to int for windows and others which use uintptr
	if err != nil {
		return "", err
	}

	fmt.Fprint(storageosCli.Out(), "\n")
	return string(p), nil
}

func runLogin(storageosCli *command.StorageOSCli, opt loginOptions, args []string) (err error) {
	opt.host, err = getHost(opt, args)
	if err != nil {
		return err
	}

	if opt.username == "" {
		opt.username, err = promptUsername(storageosCli)
		if err != nil {
			return err
		}
	}

	if opt.password == "" {
		opt.password, err = promptPassword(storageosCli)
		if err != nil {
			return err
		}
	}

	if verr := verifyCredsWithServer(opt.username, opt.password, opt.host); verr != nil {
		return verr
	}

	fmt.Fprintln(storageosCli.Out(), "Credentials verified")

	err = storageosCli.ConfigFile().CredentialsStore.SetCredentials(opt.host, opt.username, opt.password)
	if err != nil {
		return err
	}

	return storageosCli.ConfigFile().Save()
}
