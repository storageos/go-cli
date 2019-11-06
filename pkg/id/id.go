package id

import (
	"errors"
	"strings"
)

type Cluster string

func (c Cluster) String() string {
	return string(c)
}

type Node string

func (n Node) String() string {
	return string(n)
}

type Volume string

func (v Volume) String() string {
	return string(v)
}

type Deployment string

func (d Deployment) String() string {
	return string(d)
}

type Namespace string

func (n Namespace) String() string {
	return string(n)
}

// ParseFQVN parses name as a fully qualified volume name, returning the
// constituent IDs.
func ParseFQVN(name string) (Namespace, Volume, error) {
	parts := strings.Split(name, "/")
	if len(parts) != 2 {
		return "", "", errors.New("invalid fully qualified volume name")
	}

	return Namespace(parts[0]), Volume(parts[1]), nil
}
