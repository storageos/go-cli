package apiclient

// TODO(CP-3962): These should be converted to more application specific error
// types that can be used to provide additional error context to the user at
// the terminal based on the operation/resource being attempted, consider
// housing them alongside the appropriate apiclient methods.

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

// NewAuthenticationError returns a new AuthenticationError using msg as an
// optional error message if given.
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
		return "authenticated user is not authorised to perform that action"
	}
	return e.msg
}

// NewUnauthorisedError returns a new UnauthorisedError using msg as an
// optional error message if given.
func NewUnauthorisedError(msg string) UnauthorisedError {
	return UnauthorisedError{
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
		return "stale write attempted"
	}
	return e.msg
}

// NewStaleWriteError returns a new StaleWriteError using msg as an optional
// error message if given.
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
		return "target resource is in an invalid state for carrying out the request"
	}
	return e.msg
}

// NewInvalidStateTransitionError returns a new InvalidStateTransitionError
// using msg as an optional error message if given.
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

// NewLicenceCapabilityError returns a new LicenceCapabilityError using msg as
// an optional error message if given.
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
		return "server encountered internal error"
	}
	return e.msg
}

// NewServerError returns a new ServerError using msg as an optional error
// message if given.
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
		return "server encountered store outage"
	}
	return e.msg
}

// NewStoreError returns a new StoreError using msg as an optional error
// message if given.
func NewStoreError(msg string) StoreError {
	return StoreError{
		msg: msg,
	}
}
