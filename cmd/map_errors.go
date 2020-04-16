package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/c2-cli/cmd/argwrappers"
	"code.storageos.net/storageos/c2-cli/cmd/runwrappers"
	"code.storageos.net/storageos/c2-cli/pkg/labels"
	"code.storageos.net/storageos/c2-cli/pkg/selectors"
)

const (
	// LicenceCapabilityErrorCode is the exit code used by the CLI
	// when a command failed due to an API client licence capability error
	// occurring. This exit code indicates that a requested operation could
	// not be completed with the current product licence.
	LicenceCapabilityErrorCode = 5

	// AuthenticationErrorCode is the exit code used by the CLI when a
	// command failed due to an authentication error.
	AuthenticationErrorCode = 6

	// UnauthorisedErrorCode is the exit code used by the CLI when a command
	// failed due to the authenticated user not being authorised to perform
	// the required action.
	UnauthorisedErrorCode = 7

	// NotFoundCode is the exit code used by the CLI when a command cannot be
	// performed because the target could not be found.
	NotFoundCode = 8

	// AlreadyExistsCode is the exit code used by the CLI when a command could
	// not be performed because the result would conflict with an already
	// existent entity.
	AlreadyExistsCode = 9

	// InvalidInputCode is the exit code used by the CLI when a user provides
	// invalid parameter values to a command, preventing it from being carried
	// out. The input parameters for the command must be changed before
	// retrying in order to be successful.
	InvalidInputCode = 10

	// InvalidStateCode is the exit code used by the CLI when a command cannot
	// be performed on the provided target as a result of its current status.
	InvalidStateCode = 11

	// TryAgainCode is the exit code used by the CLI to indicate that the
	// command failed, but it is safe to retry it without modification.
	TryAgainCode = 12

	// InternalErrorCode is the exit code used by the CLI when a command fails
	// due to an unexpected fatal error.
	InternalErrorCode = 13

	// CommandTimedOutCode is the exit code which a command time out error
	// maps to.
	CommandTimedOutCode = 124
)

var exitCodeHelp = []struct {
	code int
	help string
}{
	{
		code: LicenceCapabilityErrorCode,
		help: "The product licence in use prevents the command from completing the operation at present",
	},
	{
		code: AuthenticationErrorCode,
		help: "Unable to authenticate with the API using the provided credentials",
	},
	{
		code: UnauthorisedErrorCode,
		help: "The authenticated user does not have the rights to perform the requested action",
	},
	{
		code: NotFoundCode,
		help: "A resource required for successful completion of the command is not found",
	},
	{
		code: AlreadyExistsCode,
		help: "The result of completing the command would conflict with an already existing resource, so it cannot be performed",
	},
	{
		code: InvalidInputCode,
		help: "An invalid value given to the command prevents completion of the operation and must be addressed before re-attempting the operation",
	},
	{
		code: InvalidStateCode,
		help: "The target's current state does not support an operation required to complete the command",
	},
	{
		code: TryAgainCode,
		help: "Indicates command failure as the result of a transient error. It can be re-attempted without modification",
	},
	{
		code: InternalErrorCode,
		help: "Command failure is prevented from completion as the result of an unexpected fatal error",
	},
	{
		code: CommandTimedOutCode,
		help: "The command reached the configured time out duration and exited early. It is recommended to check the current state before re-attempting",
	},
}

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

	//
	// Capability/Identity Errors
	//

	case errors.As(err, &apiclient.AuthenticationError{}):

		return AuthenticationErrorCode

	case errors.As(err, &apiclient.UnauthorisedError{}):

		return UnauthorisedErrorCode

	case errors.As(err, &apiclient.LicenceCapabilityError{}),
		errors.As(err, &runwrappers.LicenceLimitError{}):

		return LicenceCapabilityErrorCode

	//
	// Invalid Input
	//

	case errors.As(err, &apiclient.InvalidUserCreationError{}),
		errors.As(err, &apiclient.InvalidVolumeCreationError{}),
		errors.Is(err, labels.ErrInvalidLabelFormat),
		errors.Is(err, labels.ErrLabelKeyConflict),
		errors.Is(err, selectors.ErrInvalidSelectorFormat),
		errors.Is(err, runwrappers.ErrTargetOrSelector),
		errors.As(err, &argwrappers.InvalidArgsError{}):

		return InvalidInputCode

	//
	// Failed due to current state
	//

	case errors.As(err, &apiclient.NamespaceNotFoundError{}),
		errors.As(err, &apiclient.NodeNotFoundError{}),
		errors.As(err, &apiclient.VolumeNotFoundError{}):

		return NotFoundCode

	case errors.As(err, &apiclient.UserExistsError{}),
		errors.As(err, &apiclient.VolumeExistsError{}):

		return AlreadyExistsCode

	case errors.As(err, &apiclient.InvalidStateTransitionError{}):

		return InvalidStateCode

	//
	// Transient errors
	//

	case errors.As(err, &apiclient.StaleWriteError{}),
		errors.As(err, &apiclient.StoreError{}):

		return TryAgainCode

	case errors.Is(err, ErrCommandTimedOut):

		return CommandTimedOutCode

	//
	// Unexpected errors
	//

	case errors.As(err, &apiclient.ServerError{}):

		return InternalErrorCode

	default:
		return 1
	}
}

// newExitCodeHelpCommand initialises a non-runnable command which displays
// help documentation for the application's defined exit codes.
func newExitCodeHelpTopic() *cobra.Command {
	w := &bytes.Buffer{}

	fmt.Fprintf(w, "Certain classes of error encountered by the StorageOS CLI will set an exit code with special meaning.\n\nDefined Exit Codes:\n")

	table := uitable.New()
	table.MaxColWidth = 80
	table.Separator = "  "
	table.Wrap = true

	for _, doc := range exitCodeHelp {
		table.AddRow(fmt.Sprintf("  %d", doc.code), doc.help)
	}

	fmt.Fprintln(w, table)

	return &cobra.Command{
		Use:   "exitcodes",
		Short: "View documentation for the exit codes used by the StorageOS CLI",
		Long:  w.String(),
	}
}
