package apiclient

import (
	"context"
	"errors"
	"reflect"
	"sort"
	"testing"

	"github.com/kr/pretty"

	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/volume"
)

func TestGetVolumeByName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string

		transport *mockTransport

		volumeName string

		wantResource *volume.Resource
		wantErr      error
	}{
		{
			name: "ok",

			transport: &mockTransport{
				ListVolumesResource: map[id.Namespace][]*volume.Resource{
					"arbitrary-namespace-id": {
						&volume.Resource{
							Name: "possibly-arthur",
						},
						&volume.Resource{
							Name: "definitely-alan",
						},
					},
				},
			},

			volumeName: "definitely-alan",

			wantResource: &volume.Resource{
				Name: "definitely-alan",
			},
			wantErr: nil,
		},
		{
			name: "volume with name does not exist",

			transport: &mockTransport{
				ListVolumesResource: map[id.Namespace][]*volume.Resource{
					"arbitrary-namespace-id": {
						&volume.Resource{
							Name: "possibly-arthur",
						},
						&volume.Resource{
							Name: "not-alan",
						},
					},
				},
			},

			volumeName: "definitely-alan",

			wantResource: nil,
			wantErr: VolumeNotFoundError{
				msg:  "volume with name definitely-alan not found for target namespace",
				name: "definitely-alan",
			},
		},
		{
			name: "error getting list of volumes",

			transport: &mockTransport{
				ListVolumesError: errors.New("bananas"),
			},

			volumeName: "a-volume",

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

			gotResource, gotErr := client.GetVolumeByName(context.Background(), "arbitrary-namespace-id", tt.volumeName)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
			}

			if !reflect.DeepEqual(gotResource, tt.wantResource) {
				t.Errorf("got volume resource %v, want %v", gotResource, tt.wantResource)
			}
		})
	}
}

func TestFetchAllVolumes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string

		transport *mockTransport

		wantVolumes []*volume.Resource
		wantErr     error
	}{
		{
			name: "ok",

			transport: &mockTransport{
				// Must have at least one namespace returned to start listing volumes
				ListNamespacesResource: []*namespace.Resource{
					{
						ID: "namespace-42",
					},
					{
						ID: "namespace-43",
					},
					{
						ID: "namespace-44",
					},
				},
				ListVolumesResource: map[id.Namespace][]*volume.Resource{
					"namespace-42": {
						{
							ID: "volume-1",
						},
						{
							ID: "volume-2",
						},
					},
					"namespace-43": {
						{
							ID: "volume-3",
						},
						{
							ID: "volume-4",
						},
					},
					"namespace-44": {
						{
							ID: "volume-5",
						},
						{
							ID: "volume-6",
						},
					},
				},
			},
			wantVolumes: []*volume.Resource{
				{
					ID: "volume-1",
				},
				{
					ID: "volume-2",
				},
				{
					ID: "volume-3",
				},
				{
					ID: "volume-4",
				},
				{
					ID: "volume-5",
				},
				{
					ID: "volume-6",
				},
			},
			wantErr: nil,
		},
		{
			name: "continues when unauthorised for namespace",

			transport: &mockTransport{
				// Must have at least one namespace returned to start listing volumes
				ListNamespacesResource: []*namespace.Resource{
					{
						ID: "namespace-42",
					},
				},
				ListVolumesResource: map[id.Namespace][]*volume.Resource{
					"namespace-42": {
						{
							ID: "volume-1",
						},
						{
							ID: "volume-2",
						},
					},
				},
				ListVolumesError: UnauthorisedError{"not allowed"},
			},

			// Should keep building results when unauthorised for a namespace,
			// the client only cares about what it can access indiscriminately
			wantVolumes: []*volume.Resource{
				{
					ID: "volume-1",
				},
				{
					ID: "volume-2",
				},
			},
			wantErr: nil,
		},
		{
			name: "unexpected error listing namespaces",

			transport: &mockTransport{
				ListNamespacesError: errors.New("bananas"),
			},

			wantVolumes: nil,
			wantErr:     errors.New("bananas"),
		},
		{
			name: "unexpected error listing volumes",

			transport: &mockTransport{
				// Must have at least one namespace returned to start listing volumes
				ListNamespacesResource: []*namespace.Resource{
					{
						ID: "namespace-42",
					},
				},
				// Fatal listing error for volumes should cause a back-out with no results
				ListVolumesError: errors.New("bananas"),
			},

			wantVolumes: nil,
			wantErr:     errors.New("bananas"),
		},
	}

	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {

			client := New()
			if err := client.ConfigureTransport(tt.transport); err != nil {
				t.Fatalf("got error configuring client transport: %v", err)
			}

			gotVolumes, gotErr := client.fetchAllVolumesParallel(context.Background())
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
			}

			// sort in order to compare
			sort.Slice(gotVolumes, func(i, j int) bool {
				return gotVolumes[i].ID.String() < gotVolumes[j].ID.String()
			})

			if !reflect.DeepEqual(gotVolumes, tt.wantVolumes) {
				pretty.Ldiff(t, gotVolumes, tt.wantVolumes)
				t.Errorf("got volumes %v, want %v", pretty.Sprint(gotVolumes), pretty.Sprint(tt.wantVolumes))
			}
		})
	}
}

func TestFilterVolumesForNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string

		volumes []*volume.Resource
		names   []string

		wantVolumes []*volume.Resource
		wantErr     error
	}{
		{
			name: "don't filter when no names provided",

			volumes: []*volume.Resource{
				&volume.Resource{
					Name: "volume-a",
				},
				&volume.Resource{
					Name: "volume-b",
				},
				&volume.Resource{
					Name: "volume-c",
				},
			},
			names: nil,

			wantVolumes: []*volume.Resource{
				&volume.Resource{
					Name: "volume-a",
				},
				&volume.Resource{
					Name: "volume-b",
				},
				&volume.Resource{
					Name: "volume-c",
				},
			},
			wantErr: nil,
		},
		{
			name: "filters for provided names",

			volumes: []*volume.Resource{
				&volume.Resource{
					Name: "volume-a",
				},
				&volume.Resource{
					Name: "volume-b",
				},
				&volume.Resource{
					Name: "volume-c",
				},
			},
			names: []string{"volume-a", "volume-c"},

			wantVolumes: []*volume.Resource{
				&volume.Resource{
					Name: "volume-a",
				},
				&volume.Resource{
					Name: "volume-c",
				},
			},
			wantErr: nil,
		},
		{
			name: "errors when a provided name is not present",

			volumes: []*volume.Resource{
				&volume.Resource{
					Name: "volume-a",
				},
				&volume.Resource{
					Name: "volume-b",
				},
				&volume.Resource{
					Name: "volume-c",
				},
			},
			names: []string{"volume-a", "definitely-steve"},

			wantVolumes: nil,
			wantErr: VolumeNotFoundError{
				msg:  "volume with name definitely-steve not found for target namespace",
				name: "definitely-steve",
			},
		},
	}

	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotVolumes, gotErr := filterVolumesForNames(tt.volumes, tt.names...)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
			}

			if !reflect.DeepEqual(gotVolumes, tt.wantVolumes) {
				pretty.Ldiff(t, gotVolumes, tt.wantVolumes)
				t.Errorf("got volumes %v, want %v", pretty.Sprint(gotVolumes), pretty.Sprint(tt.wantVolumes))
			}
		})
	}
}

func TestFilterVolumesForUIDs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string

		volumes []*volume.Resource
		uids    []id.Volume

		wantVolumes []*volume.Resource
		wantErr     error
	}{
		{
			name: "don't filter when no uids provided",

			volumes: []*volume.Resource{
				&volume.Resource{
					ID: "volume-1",
				},
				&volume.Resource{
					ID: "volume-2",
				},
				&volume.Resource{
					ID: "volume-3",
				},
			},
			uids: nil,

			wantVolumes: []*volume.Resource{
				&volume.Resource{
					ID: "volume-1",
				},
				&volume.Resource{
					ID: "volume-2",
				},
				&volume.Resource{
					ID: "volume-3",
				},
			},
			wantErr: nil,
		},
		{
			name: "filters for provided uids",

			volumes: []*volume.Resource{
				&volume.Resource{
					ID: "volume-1",
				},
				&volume.Resource{
					ID: "volume-2",
				},
				&volume.Resource{
					ID: "volume-3",
				},
			},
			uids: []id.Volume{"volume-1", "volume-3"},

			wantVolumes: []*volume.Resource{
				&volume.Resource{
					ID: "volume-1",
				},
				&volume.Resource{
					ID: "volume-3",
				},
			},
			wantErr: nil,
		},
		{
			name: "errors when a provided uid is not present",

			volumes: []*volume.Resource{
				&volume.Resource{
					ID: "volume-1",
				},
				&volume.Resource{
					ID: "volume-2",
				},
				&volume.Resource{
					ID: "volume-3",
				},
			},
			uids: []id.Volume{"volume-1", "volume-42"},

			wantVolumes: nil,
			wantErr: VolumeNotFoundError{
				msg: "volume with ID volume-42 not found for target namespace",
				uid: "volume-42",
			},
		},
	}

	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotVolumes, gotErr := filterVolumesForUIDs(tt.volumes, tt.uids...)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
			}

			if !reflect.DeepEqual(gotVolumes, tt.wantVolumes) {
				pretty.Ldiff(t, gotVolumes, tt.wantVolumes)
				t.Errorf("got volumes %v, want %v", pretty.Sprint(gotVolumes), pretty.Sprint(tt.wantVolumes))
			}
		})
	}
}

func TestClientAttachVolume(t *testing.T) {
	t.Parallel()

	var mockErr = errors.New("banana error")

	tests := []struct {
		name string

		transport *mockTransport

		nsID   id.Namespace
		volID  id.Volume
		nodeID id.Node

		wantErr         error
		wantNamespaceID id.Namespace
		wantVolumeID    id.Volume
		wantNodeID      id.Node
	}{
		{
			name: "ok",

			transport: &mockTransport{
				AuthenticateError: nil,
				AttachError:       nil,
			},

			nsID:   "bananaNamespace",
			volID:  "bananaVolume",
			nodeID: "bananaNode",

			wantErr:         nil,
			wantNamespaceID: "bananaNamespace",
			wantVolumeID:    "bananaVolume",
			wantNodeID:      "bananaNode",
		},
		{
			name: "attach transport error",

			transport: &mockTransport{
				AuthenticateError: nil,
				AttachError:       mockErr,
			},

			nsID:   "bananaNamespace",
			volID:  "bananaVolume",
			nodeID: "bananaNode",

			wantErr:         mockErr,
			wantNamespaceID: "bananaNamespace",
			wantVolumeID:    "bananaVolume",
			wantNodeID:      "bananaNode",
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

			gotErr := client.AttachVolume(context.Background(), tt.nsID, tt.volID, tt.nodeID)

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
			}

			if tt.transport.AttachGotNamespace != tt.wantNamespaceID {
				t.Errorf("got %v, want %v", tt.transport.AttachGotNamespace, tt.wantNamespaceID)
			}

			if tt.transport.AttachGotVolume != tt.wantVolumeID {
				t.Errorf("got %v, want %v", tt.transport.AttachGotVolume, tt.wantVolumeID)
			}

			if tt.transport.AttachGotNode != tt.wantNodeID {
				t.Errorf("got %v, want %v", tt.transport.AttachGotNode, tt.wantNodeID)
			}
		})
	}
}
