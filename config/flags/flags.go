// Package flags exports an implementation of a configuration settings provider
// which operates using a set of flags.
package flags

import (
	"time"

	"code.storageos.net/storageos/c2-cli/config"
)

const (
	// APIEndpointsFlag keys the long flag from which the list of API host
	// endpoints are sourced, if set.
	APIEndpointsFlag = "endpoints"
	// CommandTimeoutFlag keys the long flag from which the timeout for API
	// operations is sourced, if set.
	CommandTimeoutFlag = "timeout"
	// UsernameFlag keys the long flag from which the username part of the
	// credentials used for authentication is sourced, if set.
	UsernameFlag = "username"
	// PasswordFlag keys the long flag from which the password part of the
	// credentials used for authentication is sourced, if set.
	PasswordFlag = "password"
)

// FlagSet describes a set of typed flag set accessors/setters required by the
// Provider.
type FlagSet interface {
	Changed(name string) bool

	Duration(name string, value time.Duration, usage string) *time.Duration
	String(name string, value string, usage string) *string
	StringArray(name string, value []string, usage string) *[]string

	GetDuration(name string) (time.Duration, error)
	GetString(name string) (string, error)
	GetStringArray(name string) ([]string, error)
}

// Provider exports functionality to retrieve global configuration values from
// the global flag set if available. When a configuration value is not
// available from the flag set, the configured fallback is used.
type Provider struct {
	set      FlagSet
	fallback config.Provider
}

// APIEndpoints sources the list of comma-separated target API endpoints from
// flag's FlagSet. If the value stored has not changed then flag's fallback
// is used.
func (flag *Provider) APIEndpoints() ([]string, error) {
	hosts, err := flag.set.GetStringArray(APIEndpointsFlag)
	if err != nil {
		return nil, err
	}

	if len(hosts) == 0 || !flag.set.Changed(APIEndpointsFlag) {
		return flag.fallback.APIEndpoints()
	}

	return hosts, nil
}

// CommandTimeout sources the command timeout duration from flag's FlagSet.
// If the value stored has not changed then flag's fallback is used.
func (flag *Provider) CommandTimeout() (time.Duration, error) {
	timeout, err := flag.set.GetDuration(CommandTimeoutFlag)
	if err != nil {
		return 0, err
	}

	if timeout == 0 || !flag.set.Changed(CommandTimeoutFlag) {
		return flag.fallback.CommandTimeout()
	}

	return timeout, nil
}

// Username sources the StorageOS account username to authenticate with from
// flag's FlagSet. If the value stored has not changed then flag's fallback
// is used.
func (flag *Provider) Username() (string, error) {
	username, err := flag.set.GetString(UsernameFlag)
	if err != nil {
		return "", err
	}

	if username == "" || !flag.set.Changed(UsernameFlag) {
		return flag.fallback.Username()
	}

	return username, nil
}

// Password sources the StorageOS account password to authenticate with from
// flag's FlagSet. If the value stored has not changed then flag's fallback
// is used.
func (flag *Provider) Password() (string, error) {
	password, err := flag.set.GetString(PasswordFlag)
	if err != nil {
		return "", err
	}

	if password == "" || !flag.set.Changed(PasswordFlag) {
		return flag.fallback.Password()
	}

	return password, nil
}

// NewProvider initialises a new flag based configuration provider backed by
// flagset, falling back on the provided fallback if the value can
// not be sourced from flagset.
//
//
func NewProvider(flagset FlagSet, fallback config.Provider) *Provider {

	// Set up the flags for the config provider
	flagset.StringArray(
		APIEndpointsFlag,
		[]string{config.DefaultAPIEndpoint},
		"set the list of endpoints which are used when connecting to the StorageOS API",
	)
	flagset.Duration(
		CommandTimeoutFlag,
		config.DefaultCommandTimeout,
		"set the timeout duration to use for execution of the command",
	)
	flagset.String(
		UsernameFlag,
		config.DefaultUsername,
		"set the StorageOS account username to authenticate as",
	)
	flagset.String(
		PasswordFlag,
		config.DefaultPassword,
		"set the StorageOS account password to authenticate with",
	)

	return &Provider{
		set:      flagset,
		fallback: fallback,
	}
}
