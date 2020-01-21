// Package runwrappers contains wrapper functions which implement shared
// functionality for command run functions.
package runwrappers

import (
	"context"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/pkg/cmdcontext"
)

// RunEWithContext is a function that extends a cobra.RunE function with a
// context parameter.
type RunEWithContext func(ctx context.Context, cmd *cobra.Command, args []string) error

// WrapRunEWithContext wraps next within another RunEWithContext.
type WrapRunEWithContext func(next RunEWithContext) RunEWithContext

// Chain chains wrappers from first to last, returning a wrap function that
// can be used to wrap an inner RunEWithContext.
func Chain(wrappers ...WrapRunEWithContext) WrapRunEWithContext {
	return func(next RunEWithContext) RunEWithContext {
		return func(ctx context.Context, cmd *cobra.Command, args []string) error {
			wrapped := next
			for i := len(wrappers) - 1; i >= 0; i-- {
				wrapped = wrappers[i](wrapped)
			}

			return wrapped(ctx, cmd, args)
		}
	}
}

// RunWithTimeout returns a wrapper function that uses provider to source a
// deadline for the context of the run function it is given to wrap.
func RunWithTimeout(provider cmdcontext.TimeoutProvider) WrapRunEWithContext {
	return func(next RunEWithContext) RunEWithContext {
		return func(ctx context.Context, cmd *cobra.Command, args []string) error {
			ctx, cancel := cmdcontext.WithTimeoutFromConfig(ctx, provider)
			defer cancel()

			return next(ctx, cmd, args)
		}
	}
}
