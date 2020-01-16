// Package version provides a typed wrapper around StorageOS resource versions.
package version

// Version is an opaque representation of a StorageOS API resource version.
type Version string

// FromString wraps version as a Version type.
func FromString(version string) Version {
	return Version(version)
}

// String returns the string representation of v.
func (v Version) String() string {
	return string(v)
}
