package textformat

import "time"

type mockTimeFormatter struct {
	Str string
}

func (m *mockTimeFormatter) TimeToHuman(t time.Time) string {
	return m.Str
}
