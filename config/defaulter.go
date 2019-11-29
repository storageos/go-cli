package config

import "time"

const (
	// DefaultAPIEndpoint is the default endpoint which the CLI will
	// attempt to perform API operations with.
	DefaultAPIEndpoint = "http://localhost:5705"
	// DefaultCommandTimeout is a sensible command timeout duration to use when
	// none has been specified.
	DefaultCommandTimeout = 5 * time.Second
	// DefaultPassword defines a fallback username which the CLI will
	// attempt to use in the credentials presented to the StorageOS API for
	// authentication.
	DefaultUsername = "storageos"
	// DefaultPassword defines a fallback password which the CLI will
	// attempt to use in the credentials presented to the StorageOS API for
	// authentication.
	DefaultPassword = "storageos"
)

// Defaulter exports functionality to retrieve default values for the global
// configuration settings.
//
// As a default config provider, the Defaulter does not accept fallback
// configuration providers. If there is no sensible way to default for one of
// the accessors methods, an error is allowed to be returned.
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

// Username returns a username to default to.
func (d *Defaulter) Username() (string, error) {
	return DefaultUsername, nil
}

// Password returns a password to default to.
func (d *Defaulter) Password() (string, error) {
	return DefaultPassword, nil
}

// NewDefaulter returns an initialised default config provider.
func NewDefaulter() *Defaulter {
	return &Defaulter{}
}
