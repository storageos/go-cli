package runwrappers

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/c2-cli/licence"
)

// LicenceClient is a type capable of fetching a StorageOS cluster API
// resource containing the licencing information for that installation.
type LicenceClient interface {
	GetLicence(ctx context.Context) (*licence.Resource, error)
}

// LicenceLimitError is a command line interface error wrapper for an
// apiclient.LicenceCapabilityError which uses a licence configuration to
// provide a detailed error message.
type LicenceLimitError struct {
	err     apiclient.LicenceCapabilityError
	licence *licence.Resource
}

func (e LicenceLimitError) Error() string {
	resolution := `for a licence with greater capacity or additional features contact us via https://storageos.com/contact`
	if e.licence.Kind == "basic" {
		// TODO(CP-3908): Include instructions for automatic application
		resolution = `for a free capacity upgrade, create an account and register this cluster at https://my.storageos.com`
		resolution += "\n\nFor access to additional licence features contact us via https://storageos.com/contact"
	}

	return fmt.Sprintf(`the requested operation cannot be performed with the current licence configuration. 

Reason given: %v

To resolve: %v

Current licence:

%v`, e.err.Error(), resolution, e.licence)
}

// NewLicenceLimitError uses licence to decorate err with a detailed error message.
func NewLicenceLimitError(err apiclient.LicenceCapabilityError, licence *licence.Resource) LicenceLimitError {
	return LicenceLimitError{
		err:     err,
		licence: licence,
	}
}

// HandleLicenceError returns a wrapper function that, when a licence error is
// encountered by the run function given to it, uses client to fetch the cluster
// licence configuration and decorate the returned error with extra help.
func HandleLicenceError(client LicenceClient) WrapRunEWithContext {
	return func(next RunEWithContext) RunEWithContext {
		return func(ctx context.Context, cmd *cobra.Command, args []string) error {
			err := next(ctx, cmd, args)
			switch v := err.(type) {

			case apiclient.LicenceCapabilityError:
				lic, err := client.GetLicence(ctx)
				// If the client fails to get the licence information return the
				// original error.
				if err != nil {
					return v
				}

				return NewLicenceLimitError(v, lic)

			default:
				return err
			}
		}
	}
}
