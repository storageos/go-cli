// Package jsonformat implements a JSON format output mechanism for StorageOS
// API resources.
package jsonformat

import (
	"context"
	"encoding/json"
	"io"

	"code.storageos.net/storageos/c2-cli/cluster"
	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/user"
	"code.storageos.net/storageos/c2-cli/volume"
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

//-----------------------------------------------------------------------------
// CREATE
//-----------------------------------------------------------------------------

// CreateUser encodes resource as JSON, writing the result to w.
func (d *Displayer) CreateUser(ctx context.Context, w io.Writer, resource *user.Resource) error {
	return d.encode(w, resource)
}

// CreateVolume encodes resource as JSON, writing the result to w.
func (d *Displayer) CreateVolume(ctx context.Context, w io.Writer, resource *volume.Resource) error {
	return d.encode(w, resource)
}

//-----------------------------------------------------------------------------
// GET
//-----------------------------------------------------------------------------

// GetCluster encodes resource as JSON, writing the result to w.
func (d *Displayer) GetCluster(ctx context.Context, w io.Writer, resource *cluster.Resource) error {
	return d.encode(w, resource)
}

// GetNode encodes resource as JSON, writing the result to w.
func (d *Displayer) GetNode(ctx context.Context, w io.Writer, resource *node.Resource) error {
	return d.encode(w, resource)
}

// GetNodeList encodes resources as JSON, writing the result to w.
func (d *Displayer) GetNodeList(ctx context.Context, w io.Writer, resources []*node.Resource) error {
	return d.encode(w, resources)
}

// GetVolume encodes resource as JSON, writing the result to w.
func (d *Displayer) GetVolume(ctx context.Context, w io.Writer, resource *volume.Resource) error {
	return d.encode(w, resource)
}

// GetVolumeList encodes resources as JSON, writing the result to w.
func (d *Displayer) GetVolumeList(ctx context.Context, w io.Writer, resources []*volume.Resource) error {
	return d.encode(w, resources)
}

//-----------------------------------------------------------------------------
// DESCRIBE
//-----------------------------------------------------------------------------

// DescribeNode encodes state as JSON, writing the result to w.
func (d *Displayer) DescribeNode(ctx context.Context, w io.Writer, state *node.State) error {
	return d.encode(w, state)
}

// DescribeNodeList encodes states as JSON, writing the result to w.
func (d *Displayer) DescribeNodeList(ctx context.Context, w io.Writer, states []*node.State) error {
	return d.encode(w, states)
}

// NewDisplayer initialises a Displayer which encodes StorageOS resources as
// JSON, using encoderIndent as the indentation string.
func NewDisplayer(encoderIndent string) *Displayer {
	return &Displayer{
		encoderIndent: encoderIndent,
	}
}
