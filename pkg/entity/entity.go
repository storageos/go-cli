package entity

import (
	"strings"
)

type Version string

func VersionFromString(version string) Version {
	return Version(version)
}

func (v Version) String() string {
	return string(v)
}

type Health string

const (
	Unknown Health = "unknown"
	Online         = "online"
	Offline        = "offline"
)

func HealthFromString(health string) Health {
	switch strings.ToLower(health) {
	case "online":
		return Online
	case "offline":
		return Offline
	default:
		return Unknown
	}
}
