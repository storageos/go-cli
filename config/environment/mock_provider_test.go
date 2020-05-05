package environment

import (
	"time"

	"code.storageos.net/storageos/c2-cli/config"
	"code.storageos.net/storageos/c2-cli/output"
)

type mockProvider struct {
	GetError error

	GetAuthCacheDisabled bool
	GetAPIEndpoints      []string
	GetCacheDir          string
	GetCommandTimeout    time.Duration
	GetUsername          string
	GetPassword          string
	GetUseIDs            bool
	GetNamespace         string
	GetOutput            output.Format
	GetConfigFilePath    string
}

var _ config.Provider = (*mockProvider)(nil)

func (m *mockProvider) AuthCacheDisabled() (bool, error) {
	return m.GetAuthCacheDisabled, m.GetError
}

func (m *mockProvider) APIEndpoints() ([]string, error) {
	return m.GetAPIEndpoints, m.GetError
}

func (m *mockProvider) CacheDir() (string, error) {
	return m.GetCacheDir, m.GetError
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

func (m *mockProvider) OutputFormat() (output.Format, error) {
	return m.GetOutput, m.GetError
}

func (m *mockProvider) ConfigFilePath() (string, error) {
	return m.GetConfigFilePath, m.GetError
}
