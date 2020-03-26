// Package jsonformat implements a JSON format output mechanism for StorageOS
// API resources.
package jsonformat

import (
	"context"
	"encoding/json"
	"io"

	"code.storageos.net/storageos/c2-cli/cluster"
	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/output"
)

// DefaultEncodingIndent is the encoding indent string which consumers of the
// output package can default to when initialising Displayer types.
const DefaultEncodingIndent = "\t"

// Displayer is a type which encodes StorageOS resources to JSON and writes the
// result to io.Writers.
type Displayer struct {
	encoderIndent string
}

func (d *Displayer) encode(w io.Writer, v interface{}) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", d.encoderIndent)
	return enc.Encode(v)
}

// -----------------------------------------------------------------------------
// CREATE
// -----------------------------------------------------------------------------

// CreateUser encodes user as JSON, writing the result to w.
func (d *Displayer) CreateUser(ctx context.Context, w io.Writer, user *output.User) error {
	return d.encode(w, user)
}

// CreateVolume encodes volume as JSON, writing the result to w.
func (d *Displayer) CreateVolume(ctx context.Context, w io.Writer, volume *output.Volume) error {
	return d.encode(w, volume)
}

// -----------------------------------------------------------------------------
// UPDATE
// -----------------------------------------------------------------------------

// UpdateLicence encodes licence as JSON, writing the result to w.
func (d *Displayer) UpdateLicence(ctx context.Context, w io.Writer, licence *cluster.Licence) error {
	return d.encode(w, licence)
}

// -----------------------------------------------------------------------------
// GET
// -----------------------------------------------------------------------------

// GetCluster encodes resource as JSON, writing the result to w.
func (d *Displayer) GetCluster(ctx context.Context, w io.Writer, resource *output.Cluster) error {
	return d.encode(w, resource)
}

// GetDiagnostics encodes outputPath as JSON, writing the result to w.
func (d *Displayer) GetDiagnostics(ctx context.Context, w io.Writer, outputPath string) error {
	output := struct {
		OutputPath string `json:"outputPath"`
	}{
		OutputPath: outputPath,
	}

	return d.encode(w, output)
}

// GetNode encodes resource as JSON, writing the result to w.
func (d *Displayer) GetNode(ctx context.Context, w io.Writer, resource *output.Node) error {
	return d.encode(w, resource)
}

// GetListNodes encodes resources as JSON, writing the result to w.
func (d *Displayer) GetListNodes(ctx context.Context, w io.Writer, resources []*output.Node) error {
	return d.encode(w, resources)
}

// GetNamespace encodes resource as JSON, writing the result to w.
func (d *Displayer) GetNamespace(ctx context.Context, w io.Writer, resource *output.Namespace) error {
	return d.encode(w, resource)
}

// GetListNamespaces encodes resources as JSON, writing the result to w.
func (d *Displayer) GetListNamespaces(ctx context.Context, w io.Writer, resources []*output.Namespace) error {
	return d.encode(w, resources)
}

// GetVolume encodes resource as JSON, writing the result to w.
func (d *Displayer) GetVolume(ctx context.Context, w io.Writer, volume *output.Volume) error {
	return d.encode(w, volume)
}

// GetListVolumes encodes resources as JSON, writing the result to w.
func (d *Displayer) GetListVolumes(ctx context.Context, w io.Writer, volumes []*output.Volume) error {
	return d.encode(w, volumes)
}

// -----------------------------------------------------------------------------
// DESCRIBE
// -----------------------------------------------------------------------------

// DescribeNode encodes state as JSON, writing the result to w.
func (d *Displayer) DescribeNode(ctx context.Context, w io.Writer, state *node.State) error {
	return d.encode(w, state)
}

// DescribeListNodes encodes states as JSON, writing the result to w.
func (d *Displayer) DescribeListNodes(ctx context.Context, w io.Writer, states []*node.State) error {
	return d.encode(w, states)
}

// DescribeVolume encodes volume as JSON, writing the result to w
func (d *Displayer) DescribeVolume(ctx context.Context, w io.Writer, volume *output.Volume) error {
	return d.encode(w, volume)
}

// DescribeListVolumes encodes volumes as JSON, writing the result to w
func (d *Displayer) DescribeListVolumes(ctx context.Context, w io.Writer, volumes []*output.Volume) error {
	return d.encode(w, volumes)
}

// AttachVolume writes nothing to the writer.
func (d *Displayer) AttachVolume(ctx context.Context, w io.Writer) error {
	return nil
}

// DetachVolume writes nothing to the writer.
func (d *Displayer) DetachVolume(ctx context.Context, w io.Writer) error {
	return nil
}

// NewDisplayer initialises a Displayer which encodes StorageOS resources as
// JSON, using encoderIndent as the indentation string.
func NewDisplayer(encoderIndent string) *Displayer {
	return &Displayer{
		encoderIndent: encoderIndent,
	}
}
