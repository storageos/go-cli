package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/storageos/go-cli/cli/config/configfile"
	"github.com/storageos/go-cli/pkg/homedir"
)

const (
	// ConfigFileName is the name of config file
	ConfigFileName = "config.json"
	configFileDir  = ".storageos"
)

// env vars
const (
	EnvStorageOSHost       = "STORAGEOS_HOST"
	EnvStorageosUsername   = "STORAGEOS_USERNAME"
	EnvStorageosPassword   = "STORAGEOS_PASSWORD"
	EnvStorageosAPIVersion = "STORAGEOS_API_VERSION"
)

var (
	configDir = os.Getenv("STORAGEOS_CONFIG")
)

// default features
const (
	FeatureReplicas = "storageos.feature.replicas"
)

// DeviceRootPath defines the directory in which the raw StorageOS volumes are
// created.
const DeviceRootPath = "/var/lib/storageos/volumes"

// DefaultFSType is the default filesystem we'll use if creating filesystems.
const DefaultFSType = "ext4"

func init() {
	if configDir == "" {
		configDir = filepath.Join(homedir.Get(), configFileDir)
	}
}

// Dir returns the directory the configuration file is stored in
func Dir() string {
	return configDir
}

// SetDir sets the directory the configuration file is stored in
func SetDir(dir string) {
	configDir = dir
}

// NewConfigFile initializes an empty configuration file for the given filename 'fn'
func NewConfigFile(fn string) *configfile.ConfigFile {
	return &configfile.ConfigFile{
		// HTTPHeaders: make(map[string]string),
		Filename: fn,
	}
}

// LoadFromReader is a convenience function that creates a ConfigFile object from
// a reader
func LoadFromReader(configData io.Reader) (*configfile.ConfigFile, error) {
	configFile := configfile.ConfigFile{}
	err := configFile.LoadFromReader(configData)
	return &configFile, err
}

// Load reads the configuration files in the given directory, and sets up
// the auth config information and returns values.
// FIXME: use the internal golang config parser
func Load(configDir string) (*configfile.ConfigFile, error) {
	if configDir == "" {
		configDir = Dir()
	}

	configFile := configfile.ConfigFile{
		Filename:         filepath.Join(configDir, ConfigFileName),
		CredentialsStore: make(configfile.CredStore),
	}

	// Try happy path first - latest config file
	if _, err := os.Stat(configFile.Filename); err == nil {
		file, err := os.Open(configFile.Filename)
		if err != nil {
			return &configFile, fmt.Errorf("%s - %v", configFile.Filename, err)
		}
		defer file.Close()
		err = configFile.LoadFromReader(file)
		if err != nil {
			err = fmt.Errorf("%s - %v", configFile.Filename, err)
		}
		return &configFile, err
	} else if !os.IsNotExist(err) {
		// if file is there but we can't stat it for any reason other
		// than it doesn't exist then stop
		return &configFile, fmt.Errorf("%s - %v", configFile.Filename, err)
	}
	return &configFile, nil
}
