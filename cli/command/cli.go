package command

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/dnephin/cobra"

	api "github.com/storageos/go-api"
	"github.com/storageos/go-api/serror"
	cliconfig "github.com/storageos/go-cli/cli/config"
	"github.com/storageos/go-cli/cli/config/configfile"
	cliflags "github.com/storageos/go-cli/cli/flags"
	"github.com/storageos/go-cli/pkg/jointools"
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
	cli.client, err = NewAPIClientFromFlags(opt.Common, cli.configFile)
	if err != nil {
		return err
	}

	cli.defaultVersion = cli.client.ClientVersion()
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
	return configFile
}

// NewAPIClientFromFlags creates a new APIClient from command line flags
// func NewAPIClientFromFlags(opts *cliflags.CommonOptions, configFile *configfile.ConfigFile) (client.APIClient, error) {
func NewAPIClientFromFlags(opt *cliflags.CommonOptions, configFile *configfile.ConfigFile) (*api.Client, error) {
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

	initClientAuth(host, opt, configFile, client)
	return client, nil
}

func initClientAuth(join string, opt *cliflags.CommonOptions, configFile *configfile.ConfigFile, client *api.Client) {
	var username string
	var password string

	// Env vars bind weakest to this value
	username = os.Getenv(cliconfig.EnvStorageosUsername)
	password = os.Getenv(cliconfig.EnvStorageosPassword)

	// For each host we know about, check to see if we have a stored password
	joinFragments := strings.Split(join, ",")
	for _, host := range joinFragments {
		if u, p, err := configFile.CredentialsStore.GetCredentials(host); err == nil {
			username, password = u, p
			break // don't check any more hosts, we have what we need
		}
	}

	// cli options overide the env vars and login methods
	if opt.Username != "" {
		username = opt.Username
	}
	if opt.Password != "" {
		password = opt.Password
	}

	// setup auth only if we have a password
	if username != "" && password != "" {
		client.SetAuth(username, password)
	}
}

func getServerHost(hosts string, tls bool) (host string, err error) {
	host = os.Getenv(cliconfig.EnvStorageOSHost)
	if hosts != "" {
		host = hosts
	}

	if host == "" {
		host = api.DefaultHost
	}

	// Verify and expand the join value in the host var
	if errs := jointools.VerifyJOIN(host); errs != nil {
		causes := make([]string, 0)
		help := make([]string, 0)

		for _, e := range errs {
			causes = append(causes, e.Error())
			if se, ok := e.(serror.StorageOSError); ok {
				help = append(help, se.Help())
			}
		}

		return "", serror.NewTypedStorageOSError(
			serror.InvalidHostConfig,
			errors.New(strings.Join(causes, ",")),
			"invalid host config",
			strings.Join(help, ","),
		)
	}
	return jointools.ExpandJOIN(host), nil
}

// Standard alias definitions
var (
	CreateAliases  = []string{"c"}
	InspectAliases = []string{"i"}
	ListAliases    = []string{"list"}
	UpdateAliases  = []string{"u"}
	RemoveAliases  = []string{"remove"}
	HealthAliases  = []string{"h"}
)

func WithAlias(c *cobra.Command, aliases ...string) *cobra.Command {
	c.Aliases = append(c.Aliases, aliases...)
	return c
}
