package command

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"

	"github.com/dnephin/cobra"

	api "github.com/storageos/go-api"
	cliconfig "github.com/storageos/go-cli/cli/config"
	"github.com/storageos/go-cli/cli/config/configfile"
	cliflags "github.com/storageos/go-cli/cli/flags"
	"github.com/storageos/go-cli/cli/opts"
	"github.com/storageos/go-cli/secret"
)

// Streams is an interface which exposes the standard input and output streams
type Streams interface {
	In() *InStream
	Out() *OutStream
	Err() io.Writer
}

// Cli represents the storageos command line client.
type Cli interface {
	Client() api.Client
	Out() *OutStream
	Err() io.Writer
	In() *InStream
	ConfigFile() *configfile.ConfigFile
}

// StorageOSCli is an instance the storageos command line client.
// Instances of the client can be returned from NewStorageOSCli.
type StorageOSCli struct {
	configFile      *configfile.ConfigFile
	username        string
	password        string
	in              *InStream
	out             *OutStream
	err             io.Writer
	keyFile         string
	client          *api.Client
	hasExperimental bool
	defaultVersion  string
}

// HasExperimental returns true if experimental features are accessible.
func (cli *StorageOSCli) HasExperimental() bool {
	return cli.hasExperimental
}

// DefaultVersion returns api.defaultVersion of DOCKER_API_VERSION if specified.
func (cli *StorageOSCli) DefaultVersion() string {
	return cli.defaultVersion
}

// Client returns the APIClient
func (cli *StorageOSCli) Client() *api.Client {
	return cli.client
}

// Out returns the writer used for stdout
func (cli *StorageOSCli) Out() *OutStream {
	return cli.out
}

// Err returns the writer used for stderr
func (cli *StorageOSCli) Err() io.Writer {
	return cli.err
}

// In returns the reader used for stdin
func (cli *StorageOSCli) In() *InStream {
	return cli.in
}

// ShowHelp shows the command help.
func (cli *StorageOSCli) ShowHelp(cmd *cobra.Command, args []string) error {
	cmd.SetOutput(cli.err)
	cmd.HelpFunc()(cmd, args)
	return nil
}

// ConfigFile returns the ConfigFile
func (cli *StorageOSCli) ConfigFile() *configfile.ConfigFile {
	return cli.configFile
}

// Initialize the dockerCli runs initialization that must happen after command
// line flags are parsed.
func (cli *StorageOSCli) Initialize(opt *cliflags.ClientOptions) error {
	cli.configFile = LoadDefaultConfigFile(cli.err)

	var err error
	// cli.client, err = NewAPIClientFromFlags(opts.Common, cli.configFile)
	cli.client, err = NewAPIClientFromFlags(opt.Common)
	if err != nil {
		return err
	}

	cli.defaultVersion = cli.client.ClientVersion()

	// if opts.Common.TrustKey == "" {
	// 	cli.keyFile = filepath.Join(cliconfig.Dir(), cliflags.DefaultTrustKeyFile)
	// } else {
	// 	cli.keyFile = opts.Common.TrustKey
	// }
	//
	// if ping, err := cli.client.Ping(context.Background()); err == nil {
	// 	cli.hasExperimental = ping.Experimental
	//
	// 	// since the new header was added in 1.25, assume server is 1.24 if header is not present.
	// 	if ping.APIVersion == "" {
	// 		ping.APIVersion = "1.24"
	// 	}
	//
	// 	// if server version is lower than the current cli, downgrade
	// 	if versions.LessThan(ping.APIVersion, cli.client.ClientVersion()) {
	// 		cli.client.UpdateClientVersion(ping.APIVersion)
	// 	}
	// }
	return nil
}

// NewStorageOSCli returns a StorageOSCli instance with IO output and error streams set by in, out and err.
func NewStorageOSCli(in io.ReadCloser, out, err io.Writer) *StorageOSCli {
	return &StorageOSCli{in: NewInStream(in), out: NewOutStream(out), err: err}
}

// LoadDefaultConfigFile attempts to load the default config file and returns
// an initialized ConfigFile struct if none is found.
func LoadDefaultConfigFile(err io.Writer) *configfile.ConfigFile {
	configFile, e := cliconfig.Load(cliconfig.Dir())
	if e != nil {
		fmt.Fprintf(err, "WARNING: Error loading config file:%v\n", e)
	}
	// if !configFile.ContainsAuth() {
	// 	credentials.DetectDefaultStore(configFile)
	// }
	return configFile
}

// NewAPIClientFromFlags creates a new APIClient from command line flags
// func NewAPIClientFromFlags(opts *cliflags.CommonOptions, configFile *configfile.ConfigFile) (client.APIClient, error) {
func NewAPIClientFromFlags(opt *cliflags.CommonOptions) (*api.Client, error) {
	host, err := getServerHost(opt.Hosts, opt.TLS)
	if err != nil {
		return &api.Client{}, err
	}

	verStr := api.DefaultVersionStr
	if tmpStr := os.Getenv(cliconfig.EnvStorageosAPIVersion); tmpStr != "" {
		verStr = tmpStr
	}

	client, err := api.NewVersionedClient(host, verStr)
	if err != nil {
		return &api.Client{}, err
	}

	username := getUsername(host, opt.Username)
	password := getPassword(host, opt.Password)
	if username != "" && password != "" {
		client.SetAuth(username, password)
	}

	return client, nil
}

func getServerHost(hosts []string, tls bool) (host string, err error) {
	switch len(hosts) {
	case 0:
		host = os.Getenv(cliconfig.EnvStorageOSHost)
	case 1:
		host = hosts[0]
	default:
		return "", errors.New("Please specify only one -H")
	}

	host, err = opts.ParseHost(tls, host)
	return
}

func getUsername(host, username string) string {
	if len(username) == 0 {
		envUser := os.Getenv(cliconfig.EnvStorageosUsername)

		if len(envUser) == 0 {
			if u, err := url.Parse(host); err == nil {
				user, err := secret.PromptOrGetSecret("username", secret.DefaultUserSecretName+"_"+u.Hostname(), false)
				if err != nil {
					fmt.Println("Failed to get username: " + err.Error())
				}
				return user
			}
		}

		return envUser
	}
	return username
}

func getPassword(host, password string) string {
	if len(password) == 0 {
		envPass := os.Getenv(cliconfig.EnvStorageosPassword)

		if len(envPass) == 0 {
			if u, err := url.Parse(host); err == nil {
				pass, err := secret.PromptOrGetSecret("password", secret.DefaultPassSecretName+"_"+u.Hostname(), true)
				if err != nil {
					fmt.Println("Failed to get password: " + err.Error())
				}
				return pass
			}
		}

		return envPass
	}
	return password
}
