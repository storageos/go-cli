// Package diagnostics contains abstractions for working with diagnostic bundles
// retrieved from a StorageOS cluster.
package diagnostics

import "io"

// BundleReadCloser extends an I/O ReadCloser with a bundle name accessor.
type BundleReadCloser struct {
	io.ReadCloser
	name string
}

// Named indicates if the bundle has an associated name, returning it if so.
func (b *BundleReadCloser) Named() (string, bool) {
	if b.name == "" {
		return "", false
	}
	return b.name, true
}

// NewBundleReadCloser creates a BundleReadCloser with name.
func NewBundleReadCloser(data io.ReadCloser, name string) *BundleReadCloser {
	return &BundleReadCloser{
		ReadCloser: data,
		name:       name,
	}
}
