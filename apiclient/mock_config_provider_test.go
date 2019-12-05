package apiclient

type mockConfigProvider struct {
	GetUsername      string
	GetUsernameError error

	GetPassword      string
	GetPasswordError error
}

var _ ConfigProvider = (*mockConfigProvider)(nil)

func (m *mockConfigProvider) Username() (string, error) {
	return m.GetUsername, m.GetUsernameError
}

func (m *mockConfigProvider) Password() (string, error) {
	return m.GetPassword, m.GetPasswordError
}
