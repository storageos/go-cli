package cmd

import (
	"errors"
	"fmt"
	"testing"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/c2-cli/cluster"
	"code.storageos.net/storageos/c2-cli/cmd/argwrappers"
	"code.storageos.net/storageos/c2-cli/cmd/runwrappers"
	"code.storageos.net/storageos/c2-cli/pkg/labels"
	"code.storageos.net/storageos/c2-cli/pkg/selectors"
)

func TestExitCodeForError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		forErr   error
		wantCode int
	}{
		// ---------------------------------
		// Capability/Identity Error Mapping
		// ---------------------------------
		{
			name:     "authentication error from api client",
			forErr:   apiclient.NewAuthenticationError("failed to authenticate"),
			wantCode: AuthenticationErrorCode,
		},
		{
			name:     "unauthorised error from api client",
			forErr:   apiclient.NewUnauthorisedError("not allowed to do that"),
			wantCode: UnauthorisedErrorCode,
		},
		{
			name:     "licence error from api client",
			forErr:   apiclient.NewLicenceCapabilityError("not licensed to do that"),
			wantCode: LicenceCapabilityErrorCode,
		},
		{
			name: "decorated licence error from cli",
			forErr: runwrappers.NewLicenceLimitError(
				apiclient.NewLicenceCapabilityError("not licensed to do that"),
				&cluster.Licence{CustomerName: "buddy"},
			),
			wantCode: LicenceCapabilityErrorCode,
		},
		// ---------------------
		// Invalid Input Mapping
		// ---------------------
		{
			name:     "invalid user creation from api",
			forErr:   apiclient.NewInvalidUserCreationError("username does not check out"),
			wantCode: InvalidInputCode,
		},
		{
			name:     "invalid volume creation from api",
			forErr:   apiclient.NewInvalidVolumeCreationError("not able to create that volume"),
			wantCode: InvalidInputCode,
		},
		{
			name:     "invalid labels",
			forErr:   labels.NewErrInvalidLabelFormatWithDetails("not supported format"),
			wantCode: InvalidInputCode,
		},
		{
			name:     "invalid labels",
			forErr:   labels.NewErrLabelKeyConflictWithDetails("duplicate keys"),
			wantCode: InvalidInputCode,
		},
		{
			name:     "invalid labels",
			forErr:   labels.NewErrLabelKeyConflictWithDetails("duplicate keys"),
			wantCode: InvalidInputCode,
		},
		{
			name:     "invalid selector",
			forErr:   fmt.Errorf("arbitrary message: %w", selectors.ErrInvalidSelectorFormat),
			wantCode: InvalidInputCode,
		},
		{
			name:     "invalid specified target due to conflict",
			forErr:   runwrappers.ErrTargetOrSelector,
			wantCode: InvalidInputCode,
		},
		{
			name:     "invalid argument values",
			forErr:   argwrappers.NewInvalidArgsError(errors.New("not a valid argument")),
			wantCode: InvalidInputCode,
		},
		// --------------------------
		// Current State Caused Error
		// --------------------------
		{
			name:     "namespace id not found by api client",
			forErr:   apiclient.NewNamespaceNotFoundError("namespace-id"),
			wantCode: NotFoundCode,
		},
		{
			name:     "namespace name not found by api client",
			forErr:   apiclient.NewNamespaceNameNotFoundError("namespace-name"),
			wantCode: NotFoundCode,
		},
		{
			name:     "node id not found by api client",
			forErr:   apiclient.NewNodeNotFoundError("node-id"),
			wantCode: NotFoundCode,
		},
		{
			name:     "node name not found by api client",
			forErr:   apiclient.NewNodeNameNotFoundError("node-name"),
			wantCode: NotFoundCode,
		},
		{
			name:     "volume id not found by api client",
			forErr:   apiclient.NewVolumeNotFoundError("volume-id"),
			wantCode: NotFoundCode,
		},
		{
			name:     "volume name not found by api client",
			forErr:   apiclient.NewVolumeNameNotFoundError("volume-name"),
			wantCode: NotFoundCode,
		},
		{
			name:     "user already exists",
			forErr:   apiclient.NewUserExistsError("jim"),
			wantCode: AlreadyExistsCode,
		},
		{
			name:     "volume already exists",
			forErr:   apiclient.NewVolumeExistsError("cloud-storage", "namespace-id"),
			wantCode: AlreadyExistsCode,
		},
		{
			name:     "invalid state for action",
			forErr:   apiclient.NewInvalidStateTransitionError("cannot perform that action in this state"),
			wantCode: InvalidStateCode,
		},
		// ----------------
		// Transient Errors
		// ----------------
		{
			name:     "command timeout error",
			forErr:   ErrCommandTimedOut,
			wantCode: CommandTimedOutCode,
		},
		{
			name:     "stale write, retry",
			forErr:   apiclient.NewStaleWriteError("stale write"),
			wantCode: TryAgainCode,
		},
		{
			name:     "store error, retry",
			forErr:   apiclient.NewStoreError("store error"),
			wantCode: TryAgainCode,
		},
		// -----------------
		// Unexpected Errors
		// -----------------
		{
			name:     "server error",
			forErr:   apiclient.NewServerError("unexpected server error"),
			wantCode: InternalErrorCode,
		},
		{
			name:     "arbitrary error",
			forErr:   errors.New("very arbitrary"),
			wantCode: 1,
		},
	}

	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotCode := ExitCodeForError(tt.forErr)

			if gotCode != tt.wantCode {
				t.Errorf("got exit code %v, want %v", gotCode, tt.wantCode)
			}
		})
	}
}
