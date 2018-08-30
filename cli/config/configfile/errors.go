package configfile

import (
	"errors"
)

// Errors to be returned by the credentials store when a requested operation could
// not be completed.
var (
	ErrNotFound    = errors.New("the requested password was not found in the keychain")
	ErrNotDarwin   = errors.New("keychain integration only available on darwin")
	ErrUnknownHost = errors.New("the provided host is not in the config file")
)
