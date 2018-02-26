package node

import (
	"errors"
	"testing"

	"github.com/storageos/go-api/types"
)

func TestUpdateLabel(t *testing.T) {
	testcases := []struct {
		name       string
		node       *types.Node
		label      string
		wantLabels map[string]string
		wantErr    error
	}{
		{
			name:  "basic",
			node:  &types.Node{},
			label: "country=UK",
			wantLabels: map[string]string{
				"country": "UK",
			},
		},
		{
			name:    "empty label",
			node:    &types.Node{},
			label:   "",
			wantErr: errors.New("bad attribute format: "),
		},
		{
			name:    "multiple labels",
			node:    &types.Node{},
			label:   "country=UK,load=prod",
			wantErr: errors.New("Bad label format: country=UK,load=prod"),
		},
		{
			name:    "invalid label",
			node:    &types.Node{},
			label:   "country=",
			wantErr: errors.New("Bad label format: country="),
		},
		{
			name:    "invalid label",
			node:    &types.Node{},
			label:   "=UK",
			wantErr: errors.New("Bad label format: =UK"),
		},
		{
			name: "append new label",
			node: &types.Node{
				NodeConfig: types.NodeConfig{
					Labels: map[string]string{
						"label1": "val1",
						"label2": "val2",
					},
				},
			},
			label: "country=UK",
			wantLabels: map[string]string{
				"label1":  "val1",
				"label2":  "val2",
				"country": "UK",
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := updateLabel(tc.node, tc.label)
			if err == nil {
				if err != tc.wantErr {
					t.Fatalf("unexpected error while updating node label:\n\t(GOT): %v\n\t(WNT): %v", err, tc.wantErr)
				}
			} else if err.Error() != tc.wantErr.Error() {
				t.Fatalf("unexpected error while updating node label:\n\t(GOT): %v\n\t(WNT): %v", err, tc.wantErr)
			}

			if len(tc.wantLabels) != len(tc.node.Labels) {
				t.Fatalf("unexpected number of labels while updating node label:\n\t(GOT): %d\n\t(WNT): %d", len(tc.node.Labels), len(tc.wantLabels))
			}

			for key, val := range tc.wantLabels {
				gotVal, ok := tc.node.Labels[key]
				if !ok {
					t.Fatalf("expected node to be labelled with %s", key)
				} else {
					if gotVal != val {
						t.Fatalf("unexpected node label value for label %q:\n\t(GOT): %s\n\t(WNT): %s", key, gotVal, val)
					}
				}
			}
		})
	}
}
