// Package cmdcontext exports a helper for setting up context.Context values
// to use timeout durations from a configuration provider.
package cmdcontext

import (
	"context"
	"time"

	"code.storageos.net/storageos/c2-cli/config"
)

// TimeoutProvider abstracts a type which is capable of retrieving a timeout
// duration for CLI commands from a configuration source.
type TimeoutProvider interface {
	CommandTimeout() (time.Duration, error)
}

// WithTimeoutFromConfig derives a new Context from ctx with a timeout sourced
// from provider. If an error is encountered using provider then the value of
// config.DefaultCommandTimeout is used.
func WithTimeoutFromConfig(ctx context.Context, provider TimeoutProvider) (context.Context, context.CancelFunc) {
	timeout, err := provider.CommandTimeout()
	if err != nil {
		return context.WithTimeout(ctx, config.DefaultCommandTimeout)
	}

	return context.WithTimeout(ctx, timeout)
}
