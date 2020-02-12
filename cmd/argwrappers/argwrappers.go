// Package argwrappers contains wrapper functions which implement shared
// functionality for command argument checking or input functions.
package argwrappers

import "github.com/spf13/cobra"

// InvalidArgsError is a type wrapper around the inner error denoting it as
// an invalid argument error.
type InvalidArgsError struct {
	inner error
}

func (e InvalidArgsError) Error() string {
	return e.inner.Error()
}

// NewInvalidArgsError wraps inner as an InvalidArgsError.
func NewInvalidArgsError(inner error) InvalidArgsError {
	return InvalidArgsError{
		inner: inner,
	}
}

// WrapInvalidArgsError returns a thin wrapper function around next which wraps
// the returned error as an InvalidArgsError when not nil.
func WrapInvalidArgsError(next cobra.PositionalArgs) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		err := next(cmd, args)
		if err != nil {
			err = NewInvalidArgsError(err)
		}
		return err
	}
}
