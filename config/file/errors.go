package file

import (
	"errors"
	"fmt"
)

var (
	errMissingEndpoints     = errors.New("config file error: endpoints can't be set and empty")
	errMissingUsername      = errors.New("config file error: username can't be set and empty")
	errMissingNamespace     = errors.New("config file error: namespace can't be set and empty")
	errMissingSetConfigFile = errors.New("config file error: config file path has been set but doesn't exist")
	errPasswordForbidden    = errors.New("config file error: password is not allowed in config file")
)

func errBadYAMLFile(err error) error {
	return fmt.Errorf("config file error: decoding yaml file failed: %w", err)
}
