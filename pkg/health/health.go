// Package health provides a type which is able to represent the health
// status of supported StorageOS resources.
package health

import "strings"

// State represents the health of a StorageOS resource.
type State string

const (
	// Unknown indicates that either the status of the resource could not be
	// determined or is not recognised.
	Unknown State = "unknown"
	// Online indicates the resource is functional.
	Online = "online"
	// Offline indicates the resource is not available.
	Offline = "offline"
)

// FromString returns resource state determined by the value of health.
func FromString(health string) State {
	switch strings.ToLower(health) {
	case "online":
		return Online
	case "offline":
		return Offline
	default:
		return Unknown
	}
}
