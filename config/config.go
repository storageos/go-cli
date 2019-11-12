// Package config provides utilities for parsing configuration settings
// required for operating the CLI.
package config

import (
	"time"
)

const (
	// DefaultCommandTimeout is the standard timeout for a single request to
	// the CLI's API client.
	DefaultCommandTimeout = 5 * time.Second
)
