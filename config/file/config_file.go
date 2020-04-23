package file

import (
	"time"

	"code.storageos.net/storageos/c2-cli/output"
)

// ConfigFile represents a struct that is used to parse a yaml config file.
//
//  - `Raw*` fields are used to load values from the file. They are pointer in
//    order to distinguish zero-value set fields from unset fields.
//  - `IsSet*` fields are bool values we use to know if that field has been set
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
	RawAuthCacheDisabled *string   `json:"noAuthCache" yaml:"noAuthCache"`
	RawAPIEndpoints      *[]string `json:"endpoints" yaml:"endpoints"`
	RawCacheDir          *string   `json:"cacheDir" yaml:"cacheDir"`
	RawCommandTimeout    *string   `json:"timeout" yaml:"timeout"`
	RawUsername          *string   `json:"username" yaml:"username"`
	RawPassword          *string   `json:"password" yaml:"password"`
	RawUseIDs            *string   `json:"useIds" yaml:"useIds"`
	RawNamespace         *string   `json:"namespace" yaml:"namespace"`
	RawOutputFormat      *string   `json:"output" yaml:"output"`
	// TODO(CP-3913): Support TLS.

	IsSetAuthCacheDisabled bool
	IsSetAPIEndpoints      bool
	IsSetCacheDir          bool
	IsSetCommandTimeout    bool
	IsSetUsername          bool
	IsSetUseIDs            bool
	IsSetNamespace         bool
	IsSetOutputFormat      bool
	// TODO(CP-3913): Support TLS.

	AuthCacheDisabled bool
	APIEndpoints      []string
	CacheDir          string
	CommandTimeout    time.Duration
	UsernameStr       string
	UseIDs            bool
	NamespaceStr      string
	OutputFormat      output.Format
	// TODO(CP-3913): Support TLS.

	Error error
}
