package file

import (
	"errors"
	"fmt"
)

var (
	errMissingCacheDir      = errors.New("cache directory path can't be set and empty")
	errMissingEndpoints     = errors.New("endpoints can't be set and empty")
	errMissingUsername      = errors.New("username can't be set and empty")
	errMissingNamespace     = errors.New("namespace can't be set and empty")
	errMissingSetConfigFile = errors.New("config file path has been set but doesn't exist")
	errPasswordForbidden    = errors.New("password is not allowed in config file")
)

type parseError struct {
	inner error
	path  string
}

func (e parseError) Error() string {
	return fmt.Sprintf("failed to parse config file (%q): %v", e.path, e.inner)
}

func (e parseError) Unwrap() error {
	return e.inner
}

func newParseError(err error, path string) parseError {
	return parseError{
		inner: err,
		path:  path,
	}
}
