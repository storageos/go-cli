package netutil

import (
	"errors"
	"fmt"
)

var ErrAllFailed = errors.New("failed to dial all provided addresses")
var ErrNoAddresses = errors.New("the MultiDialer instance has not been initialised with client addresses")

type InvalidNodeError struct {
	cause error
}

func (i *InvalidNodeError) Error() string {
	return fmt.Sprintf("invalid node format: %s", i.cause.Error())
}

func newInvalidNodeError(err error) error {
	return &InvalidNodeError{err}
}

var errUnsupportedScheme = errors.New("unsupported URL scheme")
var errInvalidPortNumber = errors.New("invalid port number")
