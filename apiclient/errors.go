package apiclient

import "errors"

var (
	ErrBadRequest             = errors.New("bad request")
	ErrAuthenticationRequired = errors.New("authentication required")
	ErrUnauthorised           = errors.New("unauthorised")
	ErrForbidden              = errors.New("forbidden")
	ErrNotFound               = errors.New("not found")
	ErrAlreadyExists          = errors.New("already exists")
	ErrInUse                  = errors.New("in use")
	ErrStaleWrite             = errors.New("stale write")
	ErrInvalidStateTransition = errors.New("invalid state transition")

	ErrLicenceCapacityExceeded = errors.New("licence capacity exceeded")

	ErrServerError = errors.New("server error")
	ErrStoreError  = errors.New("store error")

	// TODO: Get rid of this
	ErrUnknown = errors.New("unexpected error")
)
