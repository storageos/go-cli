package configfile

import (
	"errors"
)

var (
	ErrNotFound    = errors.New("The requested password was not found in the keychain")
	ErrNotDarwin   = errors.New("Keychain integration only available on darwin")
	ErrUnknownHost = errors.New("The provided host is not in the config file")
)
