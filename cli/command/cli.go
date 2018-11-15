package command

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/dnephin/cobra"

	api "github.com/storageos/go-api"
	"github.com/storageos/go-api/serror"
	cliconfig "github.com/storageos/go-cli/cli/config"
	"github.com/storageos/go-cli/cli/config/configfile"
	cliflags "github.com/storageos/go-cli/cli/flags"
	"github.com/storageos/go-cli/pkg/jointools"
	"github.com/storageos/go-cli/version"
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
	hosts           string
	discovery       string
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

// GetHosts returns the client's endpoints
func (cli *StorageOSCli) GetHosts() string {
	return cli.hosts
}

// GetDiscovery returns the client's discovery endpoint
func (cli *StorageOSCli) GetDiscovery() string {
	return cli.discovery
}

// GetUsername returns the client's username
func (cli *StorageOSCli) GetUsername() string {
	return cli.username
}

// GetPassword returns the client's password
func (cli *StorageOSCli) GetPassword() string {
	return cli.password
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
	cli.client.SkipServerVersionCheck = true
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

	cli.discovery = getDiscovery(opt.Common.Discovery)

	hosts, err := getServerHost(opt.Common.Hosts, opt.Common.TLS, cli.discovery)
	if err != nil {
		return err
	}
	cli.hosts = hosts

	client, err := NewAPIClientFromFlags(hosts, opt.Common, cli.configFile)
	if err != nil {
		return err
	}
	cli.client = client

	cli.defaultVersion = cli.client.ClientVersion()

	if !envHasProxy() {
		// Attempt to set HTTP proxy for the client, if set
		if proxy := cli.configFile.ProxyURL; proxy != "" {
			proxyURL, err := url.Parse(proxy)
			if err != nil {
				return errors.New("invalid proxy url")
			}
			err = cli.client.SetProxy(proxyURL)
			if err != nil {
				return err
			}
		}
	}

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
func NewAPIClientFromFlags(host string, opt *cliflags.CommonOptions, configFile *configfile.ConfigFile) (*api.Client, error) {

	if host == "" {
		return &api.Client{}, fmt.Errorf("STORAGEOS_HOST environment variable not set")
	}

	verStr := api.DefaultVersionStr
	if tmpStr := os.Getenv(cliconfig.EnvStorageosAPIVersion); tmpStr != "" {
		verStr = tmpStr
	}

	client, err := api.NewVersionedClient(host, verStr)
	if err != nil {
		return &api.Client{}, err
	}

	// Set StorageOS CLI UserAgent for all the API requests.
	client.SetUserAgent(strings.Join([]string{version.UserAgent, version.Version}, "/"))

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

func getDiscovery(discoveryFlag string) string {
	discoveryHost := os.Getenv(cliconfig.EnvStorageOSDiscovery)

	// Override the env var with flag.
	if discoveryFlag != "" {
		discoveryHost = discoveryFlag
	}
	return discoveryHost
}

func getServerHost(hosts string, tls bool, discoveryHost string) (host string, err error) {
	host = os.Getenv(cliconfig.EnvStorageOSHost)
	if hosts != "" {
		host = hosts
	}

	if host == "" {
		host = api.DefaultHost
	}

	// Verify and expand the join value in the host var
	if errs := jointools.VerifyJOIN(discoveryHost, host); errs != nil {
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
	return jointools.ExpandJOIN(discoveryHost, host), nil
}

var proxyEnvVars = [...]string{"HTTP_PROXY", "http_proxy", "HTTPS_PROXY", "https_proxy", "NO_PROXY", "no_proxy"}

// envHasProxy checks if any of the environment variables
// relating to HTTP/HTTPS proxies are set.
func envHasProxy() bool {
	for _, v := range proxyEnvVars {
		if os.Getenv(v) != "" {
			return true
		}
	}
	return false
}

// Standard alias definitions
var (
	CreateAliases  = []string{"c"}
	InspectAliases = []string{"i"}
	ListAliases    = []string{"list"}
	UpdateAliases  = []string{"u"}
	RemoveAliases  = []string{"remove"}
	HealthAliases  = []string{"h"}
	ApplyAliases   = []string{"a"}
)

// WithAlias adds the aliases given to the command.
func WithAlias(c *cobra.Command, aliases ...string) *cobra.Command {
	c.Aliases = append(c.Aliases, aliases...)
	return c
}
