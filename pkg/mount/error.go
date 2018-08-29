package mount

// Error is a typed error to be used upon the occurence
// of a mount error. It additionally holds information
// about whether the mount error is fatal or not.
type Error struct {
	Message string
	Fatal   bool
}

func (m *Error) String() string {
	if m.Fatal {
		return m.Message + " (FATAL)"
	}
	return m.Message
}

func (m *Error) Error() string { return m.Message }
