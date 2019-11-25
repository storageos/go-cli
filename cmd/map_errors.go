package cmd

import (
	"context"
	"errors"

	"code.storageos.net/storageos/c2-cli/apiclient"
)

// ExitCodeForError returns the appropriate application exit code for err.
func ExitCodeForError(err error) int {
	switch {
	case errors.Is(err, apiclient.BadRequestError{}):
		return 1 // TODO: Pick code
	case errors.Is(err, apiclient.AuthenticationError{}):
		return 1 // TODO: Pick code
	case errors.Is(err, apiclient.UnauthorisedError{}):
		return 1 // TODO: Pick code
	case errors.Is(err, apiclient.NotFoundError{}):
		return 1 // TODO: Pick code
	case errors.Is(err, apiclient.ConflictError{}):
		return 1 // TODO: Pick code
	case errors.Is(err, apiclient.StaleWriteError{}):
		return 1 // TODO: Pick code
	case errors.Is(err, apiclient.InvalidStateTransitionError{}):
		return 1 // TODO: Pick code
	case errors.Is(err, apiclient.LicenceCapabilityError{}):
		return 1 // TODO: Pick code
	case errors.Is(err, apiclient.ServerError{}):
		return 1 // TODO: Pick code
	case errors.Is(err, apiclient.StoreError{}):
		return 1 // TODO: Pick code
	case errors.Is(err, context.DeadlineExceeded),
		errors.Is(err, apiclient.ErrCommandTimedOut):
		return 124
	default:
		return 1
	}
}
