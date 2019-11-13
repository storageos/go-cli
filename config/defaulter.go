package config

import "time"

const (
	DefaultAPIEndpoint = "http://localhost:5705"

	DefaultCommandTimeout = 5 * time.Second

	DefaultUsername = "storageos"

	DefaultPassword = "storageos"
)

// Defaulter exports functionality to retrieve default values for the global
// configuration settings where sensible and errors where not.
type Defaulter struct{}

// APIEndpoints returns a slice containing the string form of the default host
// endpoint for the apiclient, http://localhost:5705.
func (d *Defaulter) APIEndpoints() ([]string, error) {
	return []string{DefaultAPIEndpoint}, nil
}

// CommandTimeout returns the standard timeout for a single command, 5 seconds.
func (d *Defaulter) CommandTimeout() (time.Duration, error) {
	return DefaultCommandTimeout, nil
}

// Username returns a username to default to. // TODO: Probably not grand, error or something?
func (d *Defaulter) Username() (string, error) {
	return DefaultUsername, nil
}

// Password returns a password to default to. // TODO: Probably not grand, error or something?
func (d *Defaulter) Password() (string, error) {
	return DefaultPassword, nil
}

func NewDefaulter() *Defaulter {
	return &Defaulter{}
}
