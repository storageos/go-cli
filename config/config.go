// Package config provides utilities for parsing configuration settings
// required for operating the CLI.
package config

import (
	"errors"
	"fmt"
	"os"
	"time"
)

const (
	APIEndpointEnvVar = "STORAGEOS_HOST"

	UsernameEnvVar = "STORAGEOS_USERNAME"

	PasswordEnvVar = "STORAGEOS_PASSWORD"
)

const (
	// DefaultCommandTimeout is the standard timeout for a single request to
	// the CLI's API client.
	DefaultCommandTimeout = 5 * time.Second
)

var (
	ErrMissingConfigFromEnv = errors.New("missing configuration setting from environment")
)

type Environment struct{}

func (e *Environment) APIEndpoint() (string, error) {
	endpoint := os.Getenv(APIEndpointEnvVar)
	if endpoint == "" {
		return "", fmt.Errorf("%w: %s", ErrMissingConfigFromEnv, APIEndpointEnvVar)
	}

	return endpoint, nil
}

func (e *Environment) Username() (string, error) {
	endpoint := os.Getenv(UsernameEnvVar)
	if endpoint == "" {
		return "", fmt.Errorf("%w: %s", ErrMissingConfigFromEnv, UsernameEnvVar)
	}

	return endpoint, nil
}

func (e *Environment) Password() (string, error) {
	endpoint := os.Getenv(PasswordEnvVar)
	if endpoint == "" {
		return "", fmt.Errorf("%w: %s", ErrMissingConfigFromEnv, PasswordEnvVar)
	}

	return endpoint, nil

}
