package cmdcontext

import "time"

type mockTimeoutProvider struct {
	CommandTimeoutReturnDuration time.Duration
	CommandTimeoutErr            error
}

func (m *mockTimeoutProvider) CommandTimeout() (time.Duration, error) {
	return m.CommandTimeoutReturnDuration, m.CommandTimeoutErr
}
