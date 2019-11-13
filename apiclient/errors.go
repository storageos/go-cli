package apiclient

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
// TODO: Difference between invalid state transition
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

type LicenceCapabilityError struct {
	msg string
}

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
