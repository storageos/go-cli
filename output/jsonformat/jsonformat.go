package jsonformat

import (
	"context"
	"encoding/json"
	"io"

	"code.storageos.net/storageos/c2-cli/cluster"
	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/volume"
)

// DefaultEncodingIndent is the encoding indent string which consumers of the
// output package can default to when initialising Displayer types.
const DefaultEncodingIndent = "\t"

type Displayer struct {
	encoderIndent string
}

func (d *Displayer) GetCluster(ctx context.Context, w io.Writer, resource *cluster.Resource) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", d.encoderIndent)
	return enc.Encode(resource)
}

func (d *Displayer) GetNode(ctx context.Context, w io.Writer, resource *node.Resource) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", d.encoderIndent)
	return enc.Encode(resource)
}

func (d *Displayer) GetNodeList(ctx context.Context, w io.Writer, resources []*node.Resource) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", d.encoderIndent)
	return enc.Encode(resources)
}

func (d *Displayer) GetVolume(ctx context.Context, w io.Writer, resource *volume.Resource) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", d.encoderIndent)
	return enc.Encode(resource)
}

func (d *Displayer) GetVolumeList(ctx context.Context, w io.Writer, resources []*volume.Resource) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", d.encoderIndent)
	return enc.Encode(resources)
}

func (d *Displayer) DescribeNode(ctx context.Context, w io.Writer, state *node.State) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", d.encoderIndent)
	return enc.Encode(state)
}

func (d *Displayer) DescribeNodeList(ctx context.Context, w io.Writer, states []*node.State) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", d.encoderIndent)
	return enc.Encode(states)
}

func NewDisplayer(encoderIndent string) *Displayer {
	return &Displayer{
		encoderIndent: encoderIndent,
	}
}
