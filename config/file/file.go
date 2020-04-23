// Package file exports an implementation of a configuration settings
// provider which operates using the user config file.
package file

import (
	"os"
	"strconv"
	"sync"
	"time"

	"gopkg.in/yaml.v3"

	"code.storageos.net/storageos/c2-cli/config"
	"code.storageos.net/storageos/c2-cli/output"
)

// ConfigFilePathProvider abstracts the functionality required by the Provider.
type ConfigFilePathProvider interface {
	ConfigFilePath() (string, error)
}

// Provider exports functionality to retrieve global configuration values from
// environment variables if available. When a configuration value is not
// available from the environment, the configured fallback is used.
type Provider struct {
	mu   sync.Mutex
	once sync.Once

	configProvider ConfigFilePathProvider
	fallback       config.Provider
	configFile     *ConfigFile

	err error
}

// AuthCacheDisabled sources the setting which determines whether the CLI must
// disable use of the auth cache from the config file, if the field is set.
// Otherwise the fallback provider is used.
func (c *Provider) AuthCacheDisabled() (bool, error) {
	if err := c.lazyInit(); err != nil {
		return false, c.err
	}

	if c.configFile.IsSetAuthCacheDisabled {
		return c.configFile.AuthCacheDisabled, nil
	}

	return c.fallback.AuthCacheDisabled()
}

// APIEndpoints sources the list of comma-separated target API endpoints from
// the config file, if the field is set. Otherwise the fallback provider
// is used.
func (c *Provider) APIEndpoints() ([]string, error) {
	if err := c.lazyInit(); err != nil {
		return nil, c.err
	}

	if c.configFile.IsSetAPIEndpoints {
		return c.configFile.APIEndpoints, nil
	}
	return c.fallback.APIEndpoints()
}

// CacheDir sources the path to the directory for the CLI to use when caching
// data from the config file, if the field is set. Otherwise the fallback provider is used.
func (c *Provider) CacheDir() (string, error) {
	if err := c.lazyInit(); err != nil {
		return "", c.err
	}

	if c.configFile.IsSetCacheDir {
		return c.configFile.CacheDir, nil
	}

	return c.fallback.CacheDir()
}

// CommandTimeout sources the command timeout duration from the config file,
// if the field is set. Otherwise the fallback provider is used.
func (c *Provider) CommandTimeout() (time.Duration, error) {
	if err := c.lazyInit(); err != nil {
		return 0, c.err
	}

	if c.configFile.IsSetCommandTimeout {
		return c.configFile.CommandTimeout, nil
	}
	return c.fallback.CommandTimeout()
}

// Username sources the StorageOS account username to authenticate with from
// the config file, if the field is set. Otherwise the fallback provider
// is used.
func (c *Provider) Username() (string, error) {
	if err := c.lazyInit(); err != nil {
		return "", c.err
	}

	if c.configFile.IsSetUsername {
		return c.configFile.UsernameStr, nil
	}
	return c.fallback.Username()
}

// Password returns the result of the fallback provider.
// We can't handle password through config file.
func (c *Provider) Password() (string, error) {
	if err := c.lazyInit(); err != nil {
		return "", c.err
	}

	return c.fallback.Password()
}

// UseIDs sources the configuration setting to specify existing API resources
// by their unique identifier instead of name from the config file, if the field
// is set. Otherwise the fallback provider is used.
func (c *Provider) UseIDs() (bool, error) {
	if err := c.lazyInit(); err != nil {
		return false, c.err
	}

	if c.configFile.IsSetUseIDs {
		return c.configFile.UseIDs, nil
	}
	return c.fallback.UseIDs()
}

// Namespace sources the StorageOS namespace to operate within from
// the config file, if the field is set. Otherwise the fallback provider
// is used.
func (c *Provider) Namespace() (string, error) {
	if err := c.lazyInit(); err != nil {
		return "", c.err
	}

	if c.configFile.IsSetNamespace {
		return c.configFile.NamespaceStr, nil
	}
	return c.fallback.Namespace()
}

// OutputFormat returns the output format type taken from
// the config file, if the field is set. Otherwise the fallback provider
// is used.
func (c *Provider) OutputFormat() (output.Format, error) {
	if err := c.lazyInit(); err != nil {
		return output.Unknown, c.err
	}

	if c.configFile.IsSetOutputFormat {
		return c.configFile.OutputFormat, nil
	}
	return c.fallback.OutputFormat()
}

// ConfigFilePath returns the config file path of the fallback provider, because
// it's impossible to define the config file path from the config file itself.
func (c *Provider) ConfigFilePath() (string, error) {
	return c.fallback.ConfigFilePath()
}

// NewProvider returns a configuration provider which sources
// its configuration setting values from the config file if this exists,
// otherwise it acts like a noop provider, passing the value from the fallback
// provider.
func NewProvider(fallback config.Provider) *Provider {
	return &Provider{
		fallback:       fallback,
		configProvider: nil,
		configFile:     nil,
	}
}

// SetConfigProvider set the config provider that will be used to retrieve the
// config file path in the lazy initialization.
func (c *Provider) SetConfigProvider(configProvider ConfigFilePathProvider) {
	c.configProvider = configProvider
}

func (c *Provider) lazyInit() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// if we never loaded the file, we do it once
	c.once.Do(func() {
		if c.configFile == nil {
			c.err = c.parse()
		}
	})

	return c.err
}

func (c *Provider) parse() error {
	// ensure at the end of this method,
	// the provider has its configFile struct
	defer func() {
		if c.configFile == nil {
			c.configFile = &ConfigFile{
				// empty struct
				// all IsSet* fields will be false
				// and fallback methods will be used
			}
		}
	}()

	configFile, err := c.configProvider.ConfigFilePath()
	if err != nil {
		return err
	}

	reader, err := os.Open(configFile)
	if err != nil {
		if os.IsNotExist(err) {

			if configFile == config.GetDefaultConfigFile() {

				// default file doesn't exists
				// likely the user is not using it
				return nil
			}

			// config file path has been changed from default
			// but it doesn't exist. This is unwanted situation.
			return errMissingSetConfigFile
		}

		// any other error from Open()
		return err
	}

	dec := yaml.NewDecoder(reader)

	// If the config file contains some unknown fields
	// we want to raise an error
	dec.KnownFields(true)

	conf := &ConfigFile{}
	err = dec.Decode(conf)
	if err != nil {
		return errBadYAMLFile(err)
	}

	if conf.RawAuthCacheDisabled != nil {
		b, err := strconv.ParseBool(*conf.RawAuthCacheDisabled)
		if err != nil {
			return err
		}

		conf.IsSetAuthCacheDisabled = true
		conf.AuthCacheDisabled = b
	}

	if conf.RawAPIEndpoints != nil {
		if len(*conf.RawAPIEndpoints) == 0 {
			return errMissingEndpoints
		}

		conf.IsSetAPIEndpoints = true
		conf.APIEndpoints = *conf.RawAPIEndpoints
	}

	if conf.RawCacheDir != nil {
		if len(*conf.RawCacheDir) == 0 {
			return errMissingCacheDir
		}

		conf.IsSetCacheDir = true
		conf.CacheDir = *conf.RawCacheDir
	}

	if conf.RawCommandTimeout != nil {
		dur, err := time.ParseDuration(*conf.RawCommandTimeout)
		if err != nil {
			return err
		}

		conf.IsSetCommandTimeout = true
		conf.CommandTimeout = dur
	}

	if conf.RawUsername != nil {
		if *conf.RawUsername == "" {
			return errMissingUsername
		}

		conf.IsSetUsername = true
		conf.UsernameStr = *conf.RawUsername
	}

	if conf.RawPassword != nil {
		return errPasswordForbidden
	}

	if conf.RawUseIDs != nil {
		b, err := strconv.ParseBool(*conf.RawUseIDs)
		if err != nil {
			return err
		}
		conf.IsSetUseIDs = true
		conf.UseIDs = b
	}

	if conf.RawNamespace != nil {
		if *conf.RawNamespace == "" {
			return errMissingNamespace
		}

		conf.IsSetNamespace = true
		conf.NamespaceStr = *conf.RawNamespace
	}

	if conf.RawOutputFormat != nil {
		outputType, err := output.FormatFromString(*conf.RawOutputFormat)
		if err != nil {
			return err
		}

		conf.IsSetOutputFormat = true
		conf.OutputFormat = outputType
	}

	// TODO(CP-3913): Support TLS.

	c.configFile = conf

	return nil
}
