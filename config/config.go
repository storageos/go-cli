// Package config provides utilities for parsing configuration settings
// required for operating the CLI.
package config

import "time"

// Provider defines the required set of configuration setting accessors
// which a type must implement in order to be used for configuring the
// application.
type Provider interface {
	APIEndpoints() ([]string, error)
	CommandTimeout() (time.Duration, error)
	Username() (string, error)
	Password() (string, error)
}
