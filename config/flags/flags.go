package flags

import "time"

const (
	// APIEndpointsFlags keys the long flag from which we source the API host
	// endpoints.
	APIEndpointsFlag = "endpoints"
	// DialTimeoutFlag keys the long flag from which we source the timeout for
	// API operations.
	DialTimeoutFlag = "timeout"
	// TODO: Maybe these don't belong here?
	UsernameFlag = "username"
	PasswordFlag = "password"
)

type FlagSet interface {
	GetDuration(name string) (time.Duration, error)
	GetString(name string) (string, error)
	GetStringArray(name string) ([]string, error)
}

type FallbackProvider interface {
	APIEndpoints() ([]string, error)
	DialTimeout() (time.Duration, error)
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

func (flag *Provider) DialTimeout() (time.Duration, error) {
	timeout, err := flag.set.GetDuration(DialTimeoutFlag)
	if err != nil {
		return 0, err
	}

	if timeout == 0 {
		return flag.fallback.DialTimeout()
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

func NewProvider(flagset FlagSet, fallback FallbackProvider) *Provider {
	return &Provider{
		set:      flagset,
		fallback: fallback,
	}
}
