// Package config provides utilities for parsing configuration settings
// required for operating the CLI.
package config

import (
	"time"

	"code.storageos.net/storageos/c2-cli/output"
)

// Provider defines the required set of configuration setting accessors
// which a type must implement in order to be used for configuring the
// application.
type Provider interface {
	APIEndpoints() ([]string, error)
	CommandTimeout() (time.Duration, error)
	Username() (string, error)
	Password() (string, error)
	UseIDs() (bool, error)
	Namespace() (string, error)
	OutputFormat() (output.Format, error)
}
