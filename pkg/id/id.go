// Package id exports a collection of typed identifiers for working with
// StorageOS resources, enabling consumers to utilise a thin layer of
// type-safety regarding identifiers.
package id

// Cluster is a unique resource identifier for a StorageOS cluster.
type Cluster string

// String returns the string representation of c.
func (c Cluster) String() string {
	return string(c)
}

// Node is a unique resource identifier for a StorageOS node.
type Node string

// String returns the string representation of n.
func (n Node) String() string {
	return string(n)
}

// Volume is a unique resource identifier for a StorageOS volume.
type Volume string

// String returns the string representation of v.
func (v Volume) String() string {
	return string(v)
}

// Deployment is a unique resource identifier for a deployment belonging to a
// StorageOS volume.
type Deployment string

// String returns the string representation of d.
func (d Deployment) String() string {
	return string(d)
}

// Namespace is a unique resource identifier for a StorageOS namespace.
type Namespace string

// String returns the string representation of n.
func (n Namespace) String() string {
	return string(n)
}

// User is a unique resource identifier for a StorageOS user account.
type User string

// String returns the string representation of u.
func (u User) String() string {
	return string(u)
}

// PolicyGroup is a unique resource identifier for a StorageOS policy group.
type PolicyGroup string

// String returns the string representation of pg.
func (pg PolicyGroup) String() string {
	return string(pg)
}
