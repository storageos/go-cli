package apiclient

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/kr/pretty"

	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/pkg/id"
)

func TestGetNodeByName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string

		transport *mockTransport

		nodeName string

		wantResource *node.Resource
		wantErr      error
	}{
		{
			name: "ok",

			transport: &mockTransport{
				ListNodesResource: []*node.Resource{
					&node.Resource{
						Name: "possibly-dave",
					},
					&node.Resource{
						Name: "definitely-steve",
					},
				},
			},

			nodeName: "definitely-steve",

			wantResource: &node.Resource{
				Name: "definitely-steve",
			},
			wantErr: nil,
		},
		{
			name: "node with name does not exist",

			transport: &mockTransport{
				ListNodesResource: []*node.Resource{
					&node.Resource{
						Name: "possibly-dave",
					},
					&node.Resource{
						Name: "not-steve",
					},
				},
			},

			nodeName: "definitely-steve",

			wantResource: nil,
			wantErr: NodeNotFoundError{
				name: "definitely-steve",
			},
		},
		{
			name: "error getting list of nodes",

			transport: &mockTransport{
				ListNodesError: errors.New("bananas"),
			},

			nodeName: "a-node",

			wantResource: nil,
			wantErr:      errors.New("bananas"),
		},
	}

	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := New()
			if err := client.ConfigureTransport(tt.transport); err != nil {
				t.Fatalf("got error configuring client transport: %v", err)
			}

			gotResource, gotErr := client.GetNodeByName(context.Background(), tt.nodeName)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
			}

			if !reflect.DeepEqual(gotResource, tt.wantResource) {
				t.Errorf("got node resource %v, want %v", gotResource, tt.wantResource)
			}
		})
	}
}

func TestFilterNodesForNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string

		nodes []*node.Resource
		names []string

		wantNodes []*node.Resource
		wantErr   error
	}{
		{
			name: "don't filter when no names provided",

			nodes: []*node.Resource{
				&node.Resource{
					Name: "node-a",
				},
				&node.Resource{
					Name: "node-b",
				},
				&node.Resource{
					Name: "node-c",
				},
			},
			names: nil,

			wantNodes: []*node.Resource{
				&node.Resource{
					Name: "node-a",
				},
				&node.Resource{
					Name: "node-b",
				},
				&node.Resource{
					Name: "node-c",
				},
			},
			wantErr: nil,
		},
		{
			name: "filters for provided names",

			nodes: []*node.Resource{
				&node.Resource{
					Name: "node-a",
				},
				&node.Resource{
					Name: "node-b",
				},
				&node.Resource{
					Name: "node-c",
				},
			},
			names: []string{"node-a", "node-c"},

			wantNodes: []*node.Resource{
				&node.Resource{
					Name: "node-a",
				},
				&node.Resource{
					Name: "node-c",
				},
			},
			wantErr: nil,
		},
		{
			name: "errors when a provided name is not present",

			nodes: []*node.Resource{
				&node.Resource{
					Name: "node-a",
				},
				&node.Resource{
					Name: "node-b",
				},
				&node.Resource{
					Name: "node-c",
				},
			},
			names: []string{"node-a", "definitely-steve"},

			wantNodes: nil,
			wantErr: NodeNotFoundError{
				name: "definitely-steve",
			},
		},
	}

	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotNodes, gotErr := filterNodesForNames(tt.nodes, tt.names...)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
			}

			if !reflect.DeepEqual(gotNodes, tt.wantNodes) {
				pretty.Ldiff(t, gotNodes, tt.wantNodes)
				t.Errorf("got nodes %v, want %v", pretty.Sprint(gotNodes), pretty.Sprint(tt.wantNodes))
			}
		})
	}
}

func TestFilterNodesForUIDs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string

		nodes []*node.Resource
		uids  []id.Node

		wantNodes []*node.Resource
		wantErr   error
	}{
		{
			name: "don't filter when no uids provided",

			nodes: []*node.Resource{
				&node.Resource{
					ID: "node-1",
				},
				&node.Resource{
					ID: "node-2",
				},
				&node.Resource{
					ID: "node-3",
				},
			},
			uids: nil,

			wantNodes: []*node.Resource{
				&node.Resource{
					ID: "node-1",
				},
				&node.Resource{
					ID: "node-2",
				},
				&node.Resource{
					ID: "node-3",
				},
			},
			wantErr: nil,
		},
		{
			name: "filters for provided uids",

			nodes: []*node.Resource{
				&node.Resource{
					ID: "node-1",
				},
				&node.Resource{
					ID: "node-2",
				},
				&node.Resource{
					ID: "node-3",
				},
			},
			uids: []id.Node{"node-1", "node-3"},

			wantNodes: []*node.Resource{
				&node.Resource{
					ID: "node-1",
				},
				&node.Resource{
					ID: "node-3",
				},
			},
			wantErr: nil,
		},
		{
			name: "errors when a provided uid is not present",

			nodes: []*node.Resource{
				&node.Resource{
					ID: "node-1",
				},
				&node.Resource{
					ID: "node-2",
				},
				&node.Resource{
					ID: "node-3",
				},
			},
			uids: []id.Node{"node-1", "node-42"},

			wantNodes: nil,
			wantErr: NodeNotFoundError{
				uid: "node-42",
			},
		},
	}

	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotNodes, gotErr := filterNodesForUIDs(tt.nodes, tt.uids...)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
			}

			if !reflect.DeepEqual(gotNodes, tt.wantNodes) {
				pretty.Ldiff(t, gotNodes, tt.wantNodes)
				t.Errorf("got nodes %v, want %v", pretty.Sprint(gotNodes), pretty.Sprint(tt.wantNodes))
			}
		})
	}
}
