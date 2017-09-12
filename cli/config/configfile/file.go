package configfile

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ConfigFile ~/.storageos/config.json file info
type ConfigFile struct {
	// HTTPHeaders          map[string]string           `json:"HttpHeaders,omitempty"`
	CredentialsStore    CredStore `json:"knownHosts,omitempty"`
	VolumesFormat       string    `json:"volumesFormat,omitempty"`
	PoolsFormat         string    `json:"poolsFormat,omitempty"`
	NamespacesFormat    string    `json:"namespacesFormat,omitempty"`
	RulesFormat         string    `json:"rulesFormat,omitempty"`
	UsersFormat         string    `json:"usersFormat,omitempty"`
	PoliciesFormat      string    `json:"policiesFormat,omitempty"`
	TemplatesFormat     string    `json:"templatesFormat,omitempty"`
	ClusterHealthFormat string    `json:"clusterHealthFormat,omitempty"`
	NodeHealthFormat    string    `json:"nodeHealthFormat,omitempty"`
	Filename            string    `json:"-"` // Note: for internal use only
}

// LoadFromReader reads the configuration data given and sets up the auth config
// information with given directory and populates the receiver object
func (configFile *ConfigFile) LoadFromReader(configData io.Reader) error {
	return json.NewDecoder(configData).Decode(&configFile)
}

// SaveToWriter encodes and writes out all the authorization information to
// the given writer
func (configFile *ConfigFile) SaveToWriter(writer io.Writer) error {
	data, err := json.MarshalIndent(configFile, "", "\t")
	if err != nil {
		return err
	}
	_, err = writer.Write(data)
	return err
}

// Save encodes and writes out all the authorization information
func (configFile *ConfigFile) Save() error {
	if configFile.Filename == "" {
		return fmt.Errorf("Can't save config with empty filename")
	}

	if err := os.MkdirAll(filepath.Dir(configFile.Filename), 0700); err != nil {
		return err
	}
	f, err := os.OpenFile(configFile.Filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return configFile.SaveToWriter(f)
}
