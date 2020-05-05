package file

import (
	"time"

	"code.storageos.net/storageos/c2-cli/config"
	"code.storageos.net/storageos/c2-cli/output"
)

// ConfigFile represents a struct that is used to parse a yaml config file.
//
//  - `Raw*` fields are used to load values from the file. They are pointer in
//    order to distinguish zero-value set fields from unset fields.
//  - `isSet*` fields are bool values we use to know if that field has been set
//    and correctly validated.
//  - other fields contains the actual value of the field, after being parsed
//    and, if necessary, converted.
//
//   If:
//   - yaml file has invalid syntax
//   - any of the fields contain an invalid value
//   - any of the fields is defined but left empty
//   - password is set
//   - in the config file there are fields we do not understand
//   we store an error in the Error field to return it as soon as the provider
//   will receive calls.
type ConfigFile struct {
	RawAuthCacheDisabled *string   `json:"noAuthCache,omitempty" yaml:"noAuthCache,omitempty"`
	RawAPIEndpoints      *[]string `json:"endpoints,omitempty" yaml:"endpoints,omitempty"`
	RawCacheDir          *string   `json:"cacheDir,omitempty" yaml:"cacheDir,omitempty"`
	RawCommandTimeout    *string   `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	RawUsername          *string   `json:"username,omitempty" yaml:"username,omitempty"`
	RawPassword          *string   `json:"password,omitempty" yaml:"password,omitempty"`
	RawUseIDs            *string   `json:"useIds,omitempty" yaml:"useIds,omitempty"`
	RawNamespace         *string   `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	RawOutputFormat      *string   `json:"output,omitempty" yaml:"output,omitempty"`
	// NOTE: If adding to or modifying this, update the below example config file help
	//
	// TODO(CP-3913): Support TLS.

	isSetAuthCacheDisabled bool
	isSetAPIEndpoints      bool
	isSetCacheDir          bool
	isSetCommandTimeout    bool
	isSetUsername          bool
	isSetUseIDs            bool
	isSetNamespace         bool
	isSetOutputFormat      bool
	// TODO(CP-3913): Support TLS.

	authCacheDisabled bool
	apiEndpoints      []string
	cacheDir          string
	commandTimeout    time.Duration
	usernameStr       string
	useIDs            bool
	namespaceStr      string
	outputFormat      output.Format
	// TODO(CP-3913): Support TLS.
}

// strptr is syntactic sugar for &strvar, preventing the need to declare a
// string variable in order to store it as a *string.
func strptr(value string) *string {
	return &value
}

// ExampleConfigFile exports an example value of a config file with set values.
var ExampleConfigFile = ConfigFile{
	RawAuthCacheDisabled: strptr("false"),
	RawAPIEndpoints:      &[]string{config.DefaultAPIEndpoint},
	RawCacheDir:          strptr(config.GetDefaultCacheDir()),
	RawCommandTimeout:    strptr(config.DefaultCommandTimeout.String()),
	RawUsername:          strptr("storageos"),
	// RawPassword is not supported - do not set it here as that is misleading.
	RawUseIDs:       strptr("false"),
	RawNamespace:    strptr("default"),
	RawOutputFormat: strptr(output.Text.String()),
}
