package cmd

import (
	"context"
	"errors"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/c2-cli/cmd/runwrappers"
)

const (
	// LicenceCapabilityErrorCode is the exit code which an API client licence
	// capability error maps to. This exit code indicates that a requested
	// operation could not be completed with the current product licence.
	LicenceCapabilityErrorCode = 5
	// CommandTimedOutCode is the exit code which a command time out error
	// maps to.
	CommandTimedOutCode = 124
)

// ErrCommandTimedOut is returned when a command's execution deadline is
// exceeded.
var ErrCommandTimedOut = errors.New("timed out performing command")

// MapCommandError attempts to map err to a user friendly error type. If
// err is not a known application-level error mapping it is returned as
// is.
func MapCommandError(err error) error {
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return ErrCommandTimedOut
	default:
		return err
	}
}

// ExitCodeForError returns the appropriate application exit code for err.
func ExitCodeForError(err error) int {
	switch {
	case errors.As(err, &apiclient.AuthenticationError{}):
		return 1 // TODO(CP-3973): Pick code

	case errors.As(err, &apiclient.UnauthorisedError{}):
		return 1 // TODO(CP-3973): Pick code
	case errors.As(err, &apiclient.NamespaceNotFoundError{}),
		errors.As(err, &apiclient.NodeNotFoundError{}),
		errors.As(err, &apiclient.VolumeNotFoundError{}):
		return 1 // TODO(CP-3973): Pick code

	case errors.As(err, &apiclient.StaleWriteError{}):
		return 1 // TODO(CP-3973): Pick code

	case errors.As(err, &apiclient.InvalidStateTransitionError{}):
		return 1 // TODO(CP-3973): Pick code

	case errors.As(err, &apiclient.LicenceCapabilityError{}),
		errors.As(err, &runwrappers.LicenceLimitError{}):
		return LicenceCapabilityErrorCode

	case errors.As(err, &apiclient.ServerError{}):
		return 1 // TODO(CP-3973): Pick code

	case errors.As(err, &apiclient.StoreError{}):
		return 1 // TODO(CP-3973): Pick code

	case errors.Is(err, context.DeadlineExceeded),
		errors.Is(err, ErrCommandTimedOut):
		return CommandTimedOutCode

	default:
		return 1
	}
}
