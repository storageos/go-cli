// Package cmdcontext exports a helper for setting up context.Context values
// to use timeout durations from a configuration provider.
package cmdcontext

import (
	"context"
	"fmt"
	"io"
	"time"

	"code.storageos.net/storageos/c2-cli/config"
)

// TimeoutProvider abstracts a type which is capable of retrieving a timeout
// duration for CLI commands from a configuration source.
type TimeoutProvider interface {
	CommandTimeout() (time.Duration, error)
}

// MinimumTimeoutProvider wraps a timeout provider, selecting the larger
// timeout duration value from the configured minimum and the value returned by
// the inner provider.
type MinimumTimeoutProvider struct {
	inner   TimeoutProvider
	minimum time.Duration
	output  io.Writer
}

// CommandTimeout returns the configured command timeout value, ensuring it
// meets the configured minimum duration. If not, the minimum is returned and
// a notification of the change is written to the output.
func (m *MinimumTimeoutProvider) CommandTimeout() (time.Duration, error) {
	provided, err := m.inner.CommandTimeout()
	if err != nil {
		return 0, err
	}

	if provided >= m.minimum {
		return provided, nil
	}

	fmt.Fprintf(m.output, "increasing command timeout to %v\n", m.minimum)

	return m.minimum, nil
}

// NewMinimumTimeoutProvider constructs a timeout provider returning either the
// specified minimum duration or the value returned by inner, whichever is
// larger.
// If the timeout providers value is overriden, a notification of this is
// written to output.
func NewMinimumTimeoutProvider(inner TimeoutProvider, minimum time.Duration, output io.Writer) *MinimumTimeoutProvider {
	return &MinimumTimeoutProvider{
		inner:   inner,
		minimum: minimum,
		output:  output,
	}
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
