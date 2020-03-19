// Package yamlformat implements a YAML format output mechanism for StorageOS
// API resources.
package yamlformat

import (
	"context"
	"io"

	"gopkg.in/yaml.v3"

	"code.storageos.net/storageos/c2-cli/cluster"
	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/output"
	"code.storageos.net/storageos/c2-cli/user"
)

// Displayer is a type which encodes StorageOS resources to YAML and writes the
// result to io.Writers.
type Displayer struct {
	encoderIndent string
}

func (d *Displayer) encode(w io.Writer, v interface{}) error {
	enc := yaml.NewEncoder(w)
	enc.SetIndent(0)
	return enc.Encode(v)
}

// -----------------------------------------------------------------------------
// CREATE
// -----------------------------------------------------------------------------

// CreateUser encodes resource as YAML, writing the result to w.
func (d *Displayer) CreateUser(ctx context.Context, w io.Writer, resource *user.Resource) error {
	return d.encode(w, resource)
}

// CreateVolume encodes resource as YAML, writing the result to w.
func (d *Displayer) CreateVolume(ctx context.Context, w io.Writer, volume *output.Volume) error {
	return d.encode(w, volume)
}

// -----------------------------------------------------------------------------
// UPDATE
// -----------------------------------------------------------------------------

// UpdateLicence encodes licence as YAML, writing the result to w.
func (d *Displayer) UpdateLicence(ctx context.Context, w io.Writer, licence *cluster.Licence) error {
	return d.encode(w, licence)
}

// -----------------------------------------------------------------------------
// GET
// -----------------------------------------------------------------------------

// GetCluster encodes resource as YAML, writing the result to w.
func (d *Displayer) GetCluster(ctx context.Context, w io.Writer, resource *cluster.Resource) error {
	return d.encode(w, resource)
}

// GetNode encodes resource as YAML, writing the result to w.
func (d *Displayer) GetNode(ctx context.Context, w io.Writer, resource *node.Resource) error {
	return d.encode(w, resource)
}

// GetListNodes encodes resources as YAML, writing the result to w.
func (d *Displayer) GetListNodes(ctx context.Context, w io.Writer, resources []*node.Resource) error {
	return d.encode(w, resources)
}

// GetNamespace encodes resource as YAML, writing the result to w.
func (d *Displayer) GetNamespace(ctx context.Context, w io.Writer, resource *namespace.Resource) error {
	return d.encode(w, resource)
}

// GetListNamespaces encodes resources as YAML, writing the result to w.
func (d *Displayer) GetListNamespaces(ctx context.Context, w io.Writer, resources []*namespace.Resource) error {
	return d.encode(w, resources)
}

// GetVolume encodes resource as YAML, writing the result to w.
func (d *Displayer) GetVolume(ctx context.Context, w io.Writer, volume *output.Volume) error {
	return d.encode(w, volume)
}

// GetListVolumes encodes resources as YAML, writing the result to w.
func (d *Displayer) GetListVolumes(ctx context.Context, w io.Writer, volumes []*output.Volume) error {
	return d.encode(w, volumes)
}

// -----------------------------------------------------------------------------
// DESCRIBE
// -----------------------------------------------------------------------------

// DescribeNode encodes state as YAML, writing the result to w.
func (d *Displayer) DescribeNode(ctx context.Context, w io.Writer, state *node.State) error {
	return d.encode(w, state)
}

// DescribeListNodes encodes states as YAML, writing the result to w.
func (d *Displayer) DescribeListNodes(ctx context.Context, w io.Writer, states []*node.State) error {
	return d.encode(w, states)
}

// AttachVolume writes nothing in the writer
func (d *Displayer) AttachVolume(ctx context.Context, w io.Writer) error {
	return nil
}

// NewDisplayer initialises a Displayer which encodes StorageOS resources as
// YAML, using encoderIndent as the indentation string.
func NewDisplayer(encoderIndent string) *Displayer {
	return &Displayer{
		encoderIndent: encoderIndent,
	}
}
