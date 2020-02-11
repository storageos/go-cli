// Package health provides a type which is able to represent the health
// status of supported StorageOS resources.
package health

// State represents the health of a StorageOS resource.
type State string

// FromString wraps health as a State type.
func FromString(health string) State {
	return State(health)
}

// String returns the string representation of s.
func (s State) String() string {
	return string(s)
}
