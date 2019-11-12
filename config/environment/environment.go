package environment

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	// APIEndpointEnvVar is the
	APIEndpointEnvVar = "STORAGEOS_HOST"

	UsernameEnvVar = "STORAGEOS_USERNAME"

	PasswordEnvVar = "STORAGEOS_PASSWORD"
)

var (
	ErrMissingConfigFromEnv = errors.New("missing configuration setting from environment")
)

type Provider struct{}

func (e *Provider) APIEndpoints() ([]string, error) {
	hostString := os.Getenv(APIEndpointEnvVar)
	if hostString == "" {
		return nil, fmt.Errorf("%w: %s", ErrMissingConfigFromEnv, APIEndpointEnvVar)
	}

	endpoints := strings.Split(hostString, ",")

	return endpoints, nil
}

func (e *Provider) Username() (string, error) {
	endpoint := os.Getenv(UsernameEnvVar)
	if endpoint == "" {
		return "", fmt.Errorf("%w: %s", ErrMissingConfigFromEnv, UsernameEnvVar)
	}

	return endpoint, nil
}

func (e *Provider) Password() (string, error) {
	endpoint := os.Getenv(PasswordEnvVar)
	if endpoint == "" {
		return "", fmt.Errorf("%w: %s", ErrMissingConfigFromEnv, PasswordEnvVar)
	}

	return endpoint, nil

}

func NewProvider() *Provider {
	// TODO: Need to use the config provider priority chain etc.
	return &Provider{}
}
