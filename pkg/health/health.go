package health

import "strings"

type State string

const (
	Unknown State = "unknown"
	Online        = "online"
	Offline       = "offline"
)

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
