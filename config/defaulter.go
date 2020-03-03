package config

import (
	"time"

	"code.storageos.net/storageos/c2-cli/output"
)

const (
	// DefaultAPIEndpoint is the default endpoint which the CLI will
	// attempt to perform API operations with.
	DefaultAPIEndpoint = "http://localhost:5705"
	// DefaultCommandTimeout is a sensible command timeout duration to use when
	// none has been specified.
	DefaultCommandTimeout = 5 * time.Second
	// DefaultUsername defines a fallback username which the CLI will
	// attempt to use in the credentials presented to the StorageOS API for
	// authentication.
	DefaultUsername = "storageos" // #nosec G101
	// DefaultPassword defines a fallback password which the CLI will
	// attempt to use in the credentials presented to the StorageOS API for
	// authentication.
	DefaultPassword = "storageos" // #nosec G101
	// DefaultUseIDs defines the default setting for using unique identifiers
	// over names when specifying existing API resources. The default is to
	// use names.
	DefaultUseIDs = false
	// DefaultNamespaceName defines the name of the default StorageOS namespace
	// which is used as a fallback when no namespace specifier is provided.
	DefaultNamespaceName = "default"
	// DefaultOutput defines the default output type for commands. Default is to
	// use text
	DefaultOutput = output.Text
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

// UseIDs returns the default value for whether API resources must be specified
// by their unique identifiers instead of names.
func (d *Defaulter) UseIDs() (bool, error) {
	return DefaultUseIDs, nil
}

// Namespace returns the namespace name "default" to use for operations which
// required a namespace to be specified.
func (d *Defaulter) Namespace() (string, error) {
	return DefaultNamespaceName, nil
}

// OutputFormat returns the default output format of the command, that is output.Text
func (d *Defaulter) OutputFormat() (output.Format, error) {
	return DefaultOutput, nil
}

var _ Provider = (*Defaulter)(nil) // Ensure that the defaulter satisfies the exported interface

// NewDefaulter returns an initialised default config provider.
func NewDefaulter() *Defaulter {
	return &Defaulter{}
}
