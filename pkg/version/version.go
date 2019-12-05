package version

type Version string

func FromString(version string) Version {
	return Version(version)
}

func (v Version) String() string {
	return string(v)
}
