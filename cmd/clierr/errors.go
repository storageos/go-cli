package clierr

import "errors"

var (
	// ErrNoNamespaceSpecified is returned when mandatory parameter -namespace
	// is not set.
	ErrNoNamespaceSpecified = errors.New("no namespace specified")
)
