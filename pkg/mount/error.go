package mount

type MountError struct {
	Message string
	Fatal   bool
}

func (m *MountError) String() string {
	if m.Fatal {
		return m.Message + " (FATAL)"
	}
	return m.Message
}

func (m *MountError) Error() string { return m.Message }
