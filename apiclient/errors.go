package apiclient

import "errors"

// ErrCommandTimedOut is returned when a command's execution deadline is
// exceeded.
var ErrCommandTimedOut = errors.New("timed out performing command")

// BadRequestError indicates that the request made by the client is invalid.
type BadRequestError struct {
	msg string
}

func (e BadRequestError) Error() string {
	if e.msg == "" {
		return "bad request"
	}
	return e.msg
}

func NewBadRequestError(msg string) BadRequestError {
	return BadRequestError{
		msg: msg,
	}
}

// AuthenticationError indicates that the requested operation could not be
// performed for the client due to an issue with the authentication credentials
// provided by the client.
type AuthenticationError struct {
	msg string
}

func (e AuthenticationError) Error() string {
	if e.msg == "" {
		return "authentication error"
	}
	return e.msg
}

func NewAuthenticationError(msg string) AuthenticationError {
	return AuthenticationError{
		msg: msg,
	}
}

// UnauthorisedError indicates that the requested operation is disallowed
// for the user which the client is authenticated as.
type UnauthorisedError struct {
	msg string
}

func (e UnauthorisedError) Error() string {
	if e.msg == "" {
		return "unauthorised"
	}
	return e.msg
}

func NewUnauthorisedError(msg string) UnauthorisedError {
	return UnauthorisedError{
		msg: msg,
	}
}

// NotFoundError indicates that a resource involved in carrying out the API
// request was not found.
type NotFoundError struct {
	msg string
}

func (e NotFoundError) Error() string {
	if e.msg == "" {
		return "not found"
	}
	return e.msg
}

func NewNotFoundError(msg string) NotFoundError {
	return NotFoundError{
		msg: msg,
	}
}

// ConflictError indicates that the requested operation could not be carried
// out due to a conflict between the current state and the desired state.
type ConflictError struct {
	msg string
}

func (e ConflictError) Error() string {
	if e.msg == "" {
		return "conflict"
	}
	return e.msg
}

func NewConflictError(msg string) ConflictError {
	return ConflictError{
		msg: msg,
	}
}

// StaleWriteError indicates that the target resource for the requested
// operation has been concurrently updated, invalidating the request. The client
// should fetch the latest version of the resource before attempting to perform
// another update.
type StaleWriteError struct {
	msg string
}

func (e StaleWriteError) Error() string {
	if e.msg == "" {
		return "stale write"
	}
	return e.msg
}

func NewStaleWriteError(msg string) StaleWriteError {
	return StaleWriteError{
		msg: msg,
	}
}

// InvalidStateTransitionError indicates that the requested operation cannot
// be performed for the target resource in its current state.
type InvalidStateTransitionError struct {
	msg string
}

func (e InvalidStateTransitionError) Error() string {
	if e.msg == "" {
		return "invalid state transition"
	}
	return e.msg
}

func NewInvalidStateTransitionError(msg string) InvalidStateTransitionError {
	return InvalidStateTransitionError{
		msg: msg,
	}
}

// LicenceCapabilityError indicates that the requested operation cannot be
// carried out due to a licensing issue with the cluster.
type LicenceCapabilityError struct {
	msg string
}

// TODO(CP-3925): This should be more helpful. Maybe decorate with upgrade
// links or suggested actions.
func (e LicenceCapabilityError) Error() string {
	if e.msg == "" {
		return "licence capability error"
	}
	return e.msg
}

func NewLicenceCapabilityError(msg string) LicenceCapabilityError {
	return LicenceCapabilityError{
		msg: msg,
	}
}

// ServerError indicates that an unrecoverable error occurred while attempting
// to perform the requested operation.
type ServerError struct {
	msg string
}

func (e ServerError) Error() string {
	if e.msg == "" {
		return "server error"
	}
	return e.msg
}

func NewServerError(msg string) ServerError {
	return ServerError{
		msg: msg,
	}
}

// StoreError indicates that the requested operation could not be performed due
// to a store outage.
type StoreError struct {
	msg string
}

func (e StoreError) Error() string {
	if e.msg == "" {
		return "store error"
	}
	return e.msg
}

func NewStoreError(msg string) StoreError {
	return StoreError{
		msg: msg,
	}
}
