// Package
package environment

import (
	"os"
	"strings"
	"time"
)

const (
	// APIEndpointsVar keys the environment variable from which we source the
	// API host endpoints.
	APIEndpointsVar = "STORAGEOS_HOST"
	// CommandTimeoutVar keys the environment variable from which we source the
	// timeout for API operations.
	CommandTimeoutVar = "STORAGEOS_TIMEOUT"
	// UsernameVar keys the environment variable from which we source the
	// username of the StorageOS account to authenticate with.
	UsernameVar = "STORAGEOS_USER_NAME"
	// PasswordVar keys the environment variable from which we source the
	// password of the StorageOS account to authenticate with.
	PasswordVar = "STORAGEOS_PASSWORD"
	// TODO(CP-3919):
	//
	// PasswordCommandVar keys the environment variable from which we optionally
	// source the password of the StorageOS account to authenticate with through
	// command execution.
	PasswordCommandVar = "STORAGEOS_PASSWORD_COMMAND"
)

type FallbackProvider interface {
	APIEndpoints() ([]string, error)
	CommandTimeout() (time.Duration, error)
	Username() (string, error)
	Password() (string, error)
}

// Provider exports functionality to retrieve global configuration values from
// environment variables if available. When a configuration value is not
// available from the environment, the configured FallbackProvider is used.
type Provider struct {
	fallback FallbackProvider
}

func (env *Provider) APIEndpoints() ([]string, error) {
	hostString := os.Getenv(APIEndpointsVar)
	if hostString == "" {
		// If there is no value to parse then fall back
		return env.fallback.APIEndpoints()
	}
	endpoints := strings.Split(hostString, ",")

	return endpoints, nil
}

func (env *Provider) CommandTimeout() (time.Duration, error) {
	timeoutString := os.Getenv(CommandTimeoutVar)
	if timeoutString == "" {
		return env.fallback.CommandTimeout()
	}

	return time.ParseDuration(timeoutString)
}

func (env *Provider) Username() (string, error) {
	username := os.Getenv(UsernameVar)
	if username == "" {
		return env.fallback.Username()
	}

	return username, nil
}

func (env *Provider) Password() (string, error) {
	password := os.Getenv(PasswordVar)
	if password == "" {
		return env.fallback.Password()
	}

	return password, nil

}

func NewProvider(fallback FallbackProvider) *Provider {
	return &Provider{
		fallback: fallback,
	}
}
