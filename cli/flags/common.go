package flags

import (
	"fmt"
	"os"

	"github.com/docker/go-connections/tlsconfig"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	cliconfig "github.com/storageos/go-cli/cli/config"
)

const (
	// DefaultTrustKeyFile is the default filename for the trust key
	DefaultTrustKeyFile = "key.json"
	// DefaultCaFile is the default filename for the CA pem file
	DefaultCaFile = "ca.pem"
	// DefaultKeyFile is the default filename for the key pem file
	DefaultKeyFile = "key.pem"
	// DefaultCertFile is the default filename for the cert pem file
	DefaultCertFile = "cert.pem"
	// FlagTLSVerify is the flag name for the TLS verification option
	FlagTLSVerify = "tlsverify"
)

var (
	storageosCertPath  = os.Getenv("DOCKER_CERT_PATH")
	storageosTLSVerify = os.Getenv("DOCKER_TLS_VERIFY") != ""
)

// CommonOptions are options common to both the client and the daemon.
type CommonOptions struct {
	Debug      bool
	Hosts      string
	Username   string
	Password   string
	LogLevel   string
	TLS        bool
	TLSVerify  bool
	TLSOptions *tlsconfig.Options
	TrustKey   string
	Discovery  string
	Timeout    int
}

// NewCommonOptions returns a new CommonOptions
func NewCommonOptions() *CommonOptions {
	return &CommonOptions{}
}

// InstallFlags adds flags for the common options on the FlagSet
func (commonOpts *CommonOptions) InstallFlags(flags *pflag.FlagSet) {
	if storageosCertPath == "" {
		storageosCertPath = cliconfig.Dir()
	}

	flags.BoolVarP(&commonOpts.Debug, "debug", "D", false, "Enable debug mode")

	flags.StringVarP(&commonOpts.Hosts, "host", "H", "", fmt.Sprintf("Node endpoint(s) to connect to (will override %s env variable value)", cliconfig.EnvStorageOSHost))
	flags.StringVarP(&commonOpts.Discovery, "discovery", "d", "", fmt.Sprintf("The discovery endpoint. Defaults to https://discovery.storageos.cloud (will override %s env variable value)", cliconfig.EnvStorageOSDiscovery))

	flags.StringVarP(&commonOpts.Username, "username", "u", "", fmt.Sprintf(`API username (will override %s env variable value)`, cliconfig.EnvStorageosUsername))
	flags.StringVarP(&commonOpts.Password, "password", "p", "", fmt.Sprintf(`API password (will override %s env variable value)`, cliconfig.EnvStorageosPassword))

	flags.IntVarP(&commonOpts.Timeout, "timeout", "t", 0, fmt.Sprintf(`client timeout in seconds (will override %s env variable value if set), default 10s`, cliconfig.EnvStorageOSTimeout))
}

// SetDefaultOptions sets default values for options after flag parsing is
// complete
func (commonOpts *CommonOptions) SetDefaultOptions(flags *pflag.FlagSet) {
	// Regardless of whether the user sets it to true or false, if they
	// specify --tlsverify at all then we need to turn on TLS
	// TLSVerify can be true even if not set due to DOCKER_TLS_VERIFY env var, so we need
	// to check that here as well
	if flags.Changed(FlagTLSVerify) || commonOpts.TLSVerify {
		commonOpts.TLS = true
	}

	if !commonOpts.TLS {
		commonOpts.TLSOptions = nil
	} else {
		tlsOptions := commonOpts.TLSOptions
		tlsOptions.InsecureSkipVerify = !commonOpts.TLSVerify

		// Reset CertFile and KeyFile to empty string if the user did not specify
		// the respective flags and the respective default files were not found.
		if !flags.Changed("tlscert") {
			if _, err := os.Stat(tlsOptions.CertFile); os.IsNotExist(err) {
				tlsOptions.CertFile = ""
			}
		}
		if !flags.Changed("tlskey") {
			if _, err := os.Stat(tlsOptions.KeyFile); os.IsNotExist(err) {
				tlsOptions.KeyFile = ""
			}
		}
	}
}

// SetLogLevel sets the logrus logging level
func SetLogLevel(logLevel string) {
	if logLevel != "" {
		lvl, err := logrus.ParseLevel(logLevel)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to parse logging level: %s\n", logLevel)
			os.Exit(1)
		}
		logrus.SetLevel(lvl)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}
