package flags

import "time"

const (
	// APIEndpointsFlags keys the long flag from which the list of API host
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

// FlagSet describes a set of typed flag set accessors required by the
// Provider.
type FlagSet interface {
	GetDuration(name string) (time.Duration, error)
	GetString(name string) (string, error)
	GetStringArray(name string) ([]string, error)
}

// FallbackProvider defines the set of methods which need to be implemented
// by a type to be used as a fallback configuration provider.
type FallbackProvider interface {
	APIEndpoints() ([]string, error)
	CommandTimeout() (time.Duration, error)
	Username() (string, error)
	Password() (string, error)
}

// Provider exports functionality to retrieve global configuration values from
// the global flag set if available. When a configuration value is not
// available from the flag set, the configured FallbackProvider is used.
type Provider struct {
	set      FlagSet
	fallback FallbackProvider
}

func (flag *Provider) APIEndpoints() ([]string, error) {
	hosts, err := flag.set.GetStringArray(APIEndpointsFlag)
	if err != nil {
		return nil, err
	}

	if len(hosts) == 0 {
		return flag.fallback.APIEndpoints()
	}

	return hosts, nil
}

func (flag *Provider) CommandTimeout() (time.Duration, error) {
	timeout, err := flag.set.GetDuration(CommandTimeoutFlag)
	if err != nil {
		return 0, err
	}

	if timeout == 0 {
		return flag.fallback.CommandTimeout()
	}

	return timeout, nil
}

func (flag *Provider) Username() (string, error) {
	username, err := flag.set.GetString(UsernameFlag)
	if err != nil {
		return "", err
	}

	if username == "" {
		return flag.fallback.Username()
	}

	return username, nil
}

func (flag *Provider) Password() (string, error) {
	password, err := flag.set.GetString(PasswordFlag)
	if err != nil {
		return "", err
	}

	if password == "" {
		return flag.fallback.Password()
	}

	return password, nil
}

// NewProvider initialises a new flag based configuration provider sourcing its
// values from flagset, falling back on the provided fallback if the value can
// not be sourced from flagset.
func NewProvider(flagset FlagSet, fallback FallbackProvider) *Provider {
	return &Provider{
		set:      flagset,
		fallback: fallback,
	}
}
