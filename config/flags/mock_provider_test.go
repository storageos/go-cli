package flags

import (
	"time"

	"code.storageos.net/storageos/c2-cli/config"
)

type mockProvider struct {
	GetError error

	GetAPIEndpoints   []string
	GetCommandTimeout time.Duration
	GetUsername       string
	GetPassword       string
	GetUseIDs         bool
	GetNamespace      string
}

var _ config.Provider = (*mockProvider)(nil)

func (m *mockProvider) APIEndpoints() ([]string, error) {
	return m.GetAPIEndpoints, m.GetError
}

func (m *mockProvider) CommandTimeout() (time.Duration, error) {
	return m.GetCommandTimeout, m.GetError
}

func (m *mockProvider) Username() (string, error) {
	return m.GetUsername, m.GetError
}

func (m *mockProvider) Password() (string, error) {
	return m.GetPassword, m.GetError
}

func (m *mockProvider) UseIDs() (bool, error) {
	return m.GetUseIDs, m.GetError
}

func (m *mockProvider) Namespace() (string, error) {
	return m.GetNamespace, m.GetError
}
