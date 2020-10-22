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
	"code.storageos.net/storageos/c2-cli/pkg/size"
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

func TestClientAttachNFSVolume(t *testing.T) {
	t.Parallel()

	var mockErr = errors.New("banana error")

	tests := []struct {
		name string

		transport *mockTransport

		nsID   id.Namespace
		volID  id.Volume
		params *AttachNFSVolumeRequestParams

		wantErr         error
		wantNamespaceID id.Namespace
		wantVolumeID    id.Volume
		wantParams      *AttachNFSVolumeRequestParams
	}{
		{
			name: "ok, params complete",

			transport: &mockTransport{},

			nsID:  "bananaNamespace",
			volID: "bananaVolume",
			params: &AttachNFSVolumeRequestParams{
				CASVersion: "42",
				AsyncMax:   42,
			},

			wantErr:         nil,
			wantNamespaceID: "bananaNamespace",
			wantVolumeID:    "bananaVolume",
			wantParams: &AttachNFSVolumeRequestParams{
				CASVersion: "42",
				AsyncMax:   42,
			},
		},
		{
			name: "ok, no params",

			transport: &mockTransport{
				GetVolumeResource: &volume.Resource{
					ID:        "bananaVolume",
					Namespace: "bananaNamespace",
					Version:   "42",
				},
			},

			nsID:   "bananaNamespace",
			volID:  "bananaVolume",
			params: nil,

			wantErr:         nil,
			wantNamespaceID: "bananaNamespace",
			wantVolumeID:    "bananaVolume",
			wantParams: &AttachNFSVolumeRequestParams{
				CASVersion: "42",
				AsyncMax:   0,
			},
		},
		{
			name: "only async, ignore version",

			transport: &mockTransport{
				GetVolumeResource: &volume.Resource{
					ID:        "bananaVolume",
					Namespace: "bananaNamespace",
					Version:   "42",
				},
			},

			nsID:  "bananaNamespace",
			volID: "bananaVolume",
			params: &AttachNFSVolumeRequestParams{
				CASVersion: "",
				AsyncMax:   42,
			},

			wantErr:         nil,
			wantNamespaceID: "bananaNamespace",
			wantVolumeID:    "bananaVolume",
			wantParams: &AttachNFSVolumeRequestParams{
				CASVersion: "42",
				AsyncMax:   42,
			},
		},
		{
			name: "attach transport error",

			transport: &mockTransport{
				AttachNFSError: mockErr,
			},

			nsID:  "bananaNamespace",
			volID: "bananaVolume",
			params: &AttachNFSVolumeRequestParams{
				CASVersion: "42",
				AsyncMax:   42,
			},

			wantErr:         mockErr,
			wantNamespaceID: "bananaNamespace",
			wantVolumeID:    "bananaVolume",
			wantParams: &AttachNFSVolumeRequestParams{
				CASVersion: "42",
				AsyncMax:   42,
			},
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

			gotErr := client.AttachNFSVolume(context.Background(), tt.nsID, tt.volID, tt.params)

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
			}

			if tt.transport.AttachNFSGotNamespace != tt.wantNamespaceID {
				t.Errorf("got %v, want %v", tt.transport.AttachGotNamespace, tt.wantNamespaceID)
			}

			if tt.transport.AttachNFSGotVolume != tt.wantVolumeID {
				t.Errorf("got %v, want %v", tt.transport.AttachGotVolume, tt.wantVolumeID)
			}

			if !reflect.DeepEqual(tt.transport.AttachNFSGotParams, tt.wantParams) {
				t.Errorf("got %v, want %v", tt.transport.AttachNFSGotParams, tt.wantParams)
			}
		})
	}
}

func TestClient_UpdateNFSVolumeMountEndpoint(t *testing.T) {
	t.Parallel()

	var mockErr = errors.New("banana error")

	tests := []struct {
		name string

		transport *mockTransport

		nsID     id.Namespace
		volID    id.Volume
		endpoint string
		params   *UpdateNFSVolumeMountEndpointRequestParams

		wantErr         error
		wantNamespaceID id.Namespace
		wantVolumeID    id.Volume
		wantParams      *UpdateNFSVolumeMountEndpointRequestParams
	}{
		{
			name: "ok, params complete",

			transport: &mockTransport{},

			nsID:  "bananaNamespace",
			volID: "bananaVolume",
			params: &UpdateNFSVolumeMountEndpointRequestParams{
				CASVersion: "42",
				AsyncMax:   42,
			},
			endpoint: "10.0.0.1:/",

			wantErr:         nil,
			wantNamespaceID: "bananaNamespace",
			wantVolumeID:    "bananaVolume",
			wantParams: &UpdateNFSVolumeMountEndpointRequestParams{
				CASVersion: "42",
				AsyncMax:   42,
			},
		},
		{
			name: "ok, no params",

			transport: &mockTransport{
				GetVolumeResource: &volume.Resource{
					ID:        "bananaVolume",
					Namespace: "bananaNamespace",
					Version:   "42",
				},
			},
			endpoint: "10.0.0.1:/",

			nsID:   "bananaNamespace",
			volID:  "bananaVolume",
			params: nil,

			wantErr:         nil,
			wantNamespaceID: "bananaNamespace",
			wantVolumeID:    "bananaVolume",
			wantParams: &UpdateNFSVolumeMountEndpointRequestParams{
				CASVersion: "42",
				AsyncMax:   0,
			},
		},
		{
			name: "only async, ignore version",

			transport: &mockTransport{
				GetVolumeResource: &volume.Resource{
					ID:        "bananaVolume",
					Namespace: "bananaNamespace",
					Version:   "42",
				},
			},

			nsID:  "bananaNamespace",
			volID: "bananaVolume",
			params: &UpdateNFSVolumeMountEndpointRequestParams{
				CASVersion: "",
				AsyncMax:   42,
			},
			endpoint: "10.0.0.1:/",

			wantErr:         nil,
			wantNamespaceID: "bananaNamespace",
			wantVolumeID:    "bananaVolume",
			wantParams: &UpdateNFSVolumeMountEndpointRequestParams{
				CASVersion: "42",
				AsyncMax:   42,
			},
		},
		{
			name: "update NFS volume mount endpoint transport error",

			transport: &mockTransport{
				UpdateNFSVolumeMountEndpointError: mockErr,
			},

			nsID:  "bananaNamespace",
			volID: "bananaVolume",
			params: &UpdateNFSVolumeMountEndpointRequestParams{
				CASVersion: "42",
				AsyncMax:   42,
			},
			endpoint: "10.0.0.1:/",

			wantErr:         mockErr,
			wantNamespaceID: "bananaNamespace",
			wantVolumeID:    "bananaVolume",
			wantParams: &UpdateNFSVolumeMountEndpointRequestParams{
				CASVersion: "42",
				AsyncMax:   42,
			},
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

			gotErr := client.UpdateNFSVolumeMountEndpoint(context.Background(), tt.nsID, tt.volID, tt.endpoint, tt.params)

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
			}

			if tt.transport.UpdateNFSVolumeMountEndpointGotNamespaceID != tt.wantNamespaceID {
				t.Errorf("got %v, want %v", tt.transport.UpdateNFSVolumeMountEndpointGotNamespaceID, tt.wantNamespaceID)
			}

			if tt.transport.UpdateNFSVolumeMountEndpointGotVolumeID != tt.wantVolumeID {
				t.Errorf("got %v, want %v", tt.transport.UpdateNFSVolumeMountEndpointGotVolumeID, tt.wantVolumeID)
			}

			if tt.transport.UpdateNFSVolumeMountEndpointGotEndpoint != tt.endpoint {
				t.Errorf("got %v, want %v", tt.transport.UpdateNFSVolumeMountEndpointGotEndpoint, tt.endpoint)
			}

			if !reflect.DeepEqual(tt.transport.UpdateNFSVolumeMountEndpointGotParams, tt.wantParams) {
				t.Errorf("got %v, want %v", tt.transport.UpdateNFSVolumeMountEndpointGotParams, tt.wantParams)
			}
		})
	}
}

func TestClient_UpdateNFSVolumeExports(t *testing.T) {
	t.Parallel()

	var mockErr = errors.New("banana error")

	tests := []struct {
		name string

		transport *mockTransport

		nsID    id.Namespace
		volID   id.Volume
		exports []volume.NFSExportConfig
		params  *UpdateNFSVolumeExportsRequestParams

		wantErr         error
		wantNamespaceID id.Namespace
		wantVolumeID    id.Volume
		wantParams      *UpdateNFSVolumeExportsRequestParams
	}{
		{
			name: "ok, params complete",

			transport: &mockTransport{},

			nsID:  "bananaNamespace",
			volID: "bananaVolume",
			params: &UpdateNFSVolumeExportsRequestParams{
				CASVersion: "42",
				AsyncMax:   42,
			},
			exports: []volume.NFSExportConfig{
				{
					ExportID:   0,
					Path:       "/",
					PseudoPath: "/",
					ACLs: []volume.NFSExportConfigACL{
						{
							Identity: volume.NFSExportConfigACLIdentity{
								IdentityType: "cidr",
								Matcher:      "10.0.0.1/8",
							},
							SquashConfig: volume.NFSExportConfigACLSquashConfig{
								GID:    1001,
								UID:    1000,
								Squash: "all",
							},
							AccessLevel: "rw",
						},
					},
				},
				{
					ExportID:   1,
					Path:       "/path",
					PseudoPath: "/pseudo",
					ACLs: []volume.NFSExportConfigACL{
						{
							Identity: volume.NFSExportConfigACLIdentity{
								IdentityType: "hostname",
								Matcher:      "*.storageos.com",
							},
							SquashConfig: volume.NFSExportConfigACLSquashConfig{
								GID:    1001,
								UID:    1000,
								Squash: "root",
							},
							AccessLevel: "ro",
						},
					},
				},
			},

			wantErr:         nil,
			wantNamespaceID: "bananaNamespace",
			wantVolumeID:    "bananaVolume",
			wantParams: &UpdateNFSVolumeExportsRequestParams{
				CASVersion: "42",
				AsyncMax:   42,
			},
		},
		{
			name: "ok, no params",

			transport: &mockTransport{
				GetVolumeResource: &volume.Resource{
					ID:        "bananaVolume",
					Namespace: "bananaNamespace",
					Version:   "42",
				},
			},
			exports: []volume.NFSExportConfig{
				{
					ExportID:   0,
					Path:       "/",
					PseudoPath: "/",
					ACLs: []volume.NFSExportConfigACL{
						{
							Identity: volume.NFSExportConfigACLIdentity{
								IdentityType: "cidr",
								Matcher:      "10.0.0.1/8",
							},
							SquashConfig: volume.NFSExportConfigACLSquashConfig{
								GID:    1001,
								UID:    1000,
								Squash: "all",
							},
							AccessLevel: "rw",
						},
					},
				},
				{
					ExportID:   1,
					Path:       "/path",
					PseudoPath: "/pseudo",
					ACLs: []volume.NFSExportConfigACL{
						{
							Identity: volume.NFSExportConfigACLIdentity{
								IdentityType: "hostname",
								Matcher:      "*.storageos.com",
							},
							SquashConfig: volume.NFSExportConfigACLSquashConfig{
								GID:    1001,
								UID:    1000,
								Squash: "root",
							},
							AccessLevel: "ro",
						},
					},
				},
			},

			nsID:   "bananaNamespace",
			volID:  "bananaVolume",
			params: nil,

			wantErr:         nil,
			wantNamespaceID: "bananaNamespace",
			wantVolumeID:    "bananaVolume",
			wantParams: &UpdateNFSVolumeExportsRequestParams{
				CASVersion: "42",
				AsyncMax:   0,
			},
		},
		{
			name: "only async, ignore version",

			transport: &mockTransport{
				GetVolumeResource: &volume.Resource{
					ID:        "bananaVolume",
					Namespace: "bananaNamespace",
					Version:   "42",
				},
			},

			nsID:  "bananaNamespace",
			volID: "bananaVolume",
			params: &UpdateNFSVolumeExportsRequestParams{
				CASVersion: "",
				AsyncMax:   42,
			},
			exports: []volume.NFSExportConfig{
				{
					ExportID:   0,
					Path:       "/",
					PseudoPath: "/",
					ACLs: []volume.NFSExportConfigACL{
						{
							Identity: volume.NFSExportConfigACLIdentity{
								IdentityType: "cidr",
								Matcher:      "10.0.0.1/8",
							},
							SquashConfig: volume.NFSExportConfigACLSquashConfig{
								GID:    1001,
								UID:    1000,
								Squash: "all",
							},
							AccessLevel: "rw",
						},
					},
				},
				{
					ExportID:   1,
					Path:       "/path",
					PseudoPath: "/pseudo",
					ACLs: []volume.NFSExportConfigACL{
						{
							Identity: volume.NFSExportConfigACLIdentity{
								IdentityType: "hostname",
								Matcher:      "*.storageos.com",
							},
							SquashConfig: volume.NFSExportConfigACLSquashConfig{
								GID:    1001,
								UID:    1000,
								Squash: "root",
							},
							AccessLevel: "ro",
						},
					},
				},
			},

			wantErr:         nil,
			wantNamespaceID: "bananaNamespace",
			wantVolumeID:    "bananaVolume",
			wantParams: &UpdateNFSVolumeExportsRequestParams{
				CASVersion: "42",
				AsyncMax:   42,
			},
		},
		{
			name: "update NFS volume exports transport error",

			transport: &mockTransport{
				UpdateNFSVolumeExportsError: mockErr,
			},

			nsID:  "bananaNamespace",
			volID: "bananaVolume",
			params: &UpdateNFSVolumeExportsRequestParams{
				CASVersion: "42",
				AsyncMax:   42,
			},
			exports: []volume.NFSExportConfig{
				{
					ExportID:   0,
					Path:       "/",
					PseudoPath: "/",
					ACLs: []volume.NFSExportConfigACL{
						{
							Identity: volume.NFSExportConfigACLIdentity{
								IdentityType: "cidr",
								Matcher:      "10.0.0.1/8",
							},
							SquashConfig: volume.NFSExportConfigACLSquashConfig{
								GID:    1001,
								UID:    1000,
								Squash: "all",
							},
							AccessLevel: "rw",
						},
					},
				},
				{
					ExportID:   1,
					Path:       "/path",
					PseudoPath: "/pseudo",
					ACLs: []volume.NFSExportConfigACL{
						{
							Identity: volume.NFSExportConfigACLIdentity{
								IdentityType: "hostname",
								Matcher:      "*.storageos.com",
							},
							SquashConfig: volume.NFSExportConfigACLSquashConfig{
								GID:    1001,
								UID:    1000,
								Squash: "root",
							},
							AccessLevel: "ro",
						},
					},
				},
			},

			wantErr:         mockErr,
			wantNamespaceID: "bananaNamespace",
			wantVolumeID:    "bananaVolume",
			wantParams: &UpdateNFSVolumeExportsRequestParams{
				CASVersion: "42",
				AsyncMax:   42,
			},
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

			gotErr := client.UpdateNFSVolumeExports(context.Background(), tt.nsID, tt.volID, tt.exports, tt.params)

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
			}

			if tt.transport.UpdateNFSVolumeExportsGotNamespaceID != tt.wantNamespaceID {
				t.Errorf("got %v, want %v", tt.transport.UpdateNFSVolumeExportsGotNamespaceID, tt.wantNamespaceID)
			}

			if tt.transport.UpdateNFSVolumeExportsGotVolumeID != tt.wantVolumeID {
				t.Errorf("got %v, want %v", tt.transport.UpdateNFSVolumeExportsGotVolumeID, tt.wantVolumeID)
			}

			if !reflect.DeepEqual(tt.transport.UpdateNFSVolumeExportsGotExports, tt.exports) {
				t.Errorf("got %v, want %v", tt.transport.UpdateNFSVolumeExportsGotExports, tt.exports)
			}

			if !reflect.DeepEqual(tt.transport.UpdateNFSVolumeExportsGotParams, tt.wantParams) {
				t.Errorf("got %v, want %v", tt.transport.UpdateNFSVolumeExportsGotParams, tt.wantParams)
			}
		})
	}
}

func TestClient_ResizeVolume(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string

		Transport *mockTransport
		params    *ResizeVolumeRequestParams

		wantParams *ResizeVolumeRequestParams
	}{
		{
			name: "params input nil",
			Transport: &mockTransport{
				GetVolumeResource: &volume.Resource{Version: "42"},
			},
			params: nil,
			wantParams: &ResizeVolumeRequestParams{
				AsyncMax:   0,
				CASVersion: "42",
			},
		},
		{
			name: "params input both empty",
			Transport: &mockTransport{
				GetVolumeResource: &volume.Resource{Version: "42"},
			},
			params: &ResizeVolumeRequestParams{
				AsyncMax:   0,
				CASVersion: "",
			},
			wantParams: &ResizeVolumeRequestParams{
				AsyncMax:   0,
				CASVersion: "42",
			},
		},
		{
			name: "params input version set, async empty",
			Transport: &mockTransport{
				GetVolumeResource: &volume.Resource{Version: "43"},
			},
			params: &ResizeVolumeRequestParams{
				AsyncMax:   0,
				CASVersion: "42",
			},
			wantParams: &ResizeVolumeRequestParams{
				AsyncMax:   0,
				CASVersion: "42",
			},
		},
		{
			name: "params input version empty, async set",
			Transport: &mockTransport{
				GetVolumeResource: &volume.Resource{Version: "42"},
			},
			params: &ResizeVolumeRequestParams{
				AsyncMax:   42,
				CASVersion: "",
			},
			wantParams: &ResizeVolumeRequestParams{
				AsyncMax:   42,
				CASVersion: "42",
			},
		},
		{
			name: "params input both set",
			Transport: &mockTransport{
				GetVolumeResource: &volume.Resource{Version: "43"},
			},
			params: &ResizeVolumeRequestParams{
				AsyncMax:   42,
				CASVersion: "42",
			},
			wantParams: &ResizeVolumeRequestParams{
				AsyncMax:   42,
				CASVersion: "42",
			},
		},
	}
	for _, tt := range tests {
		var tt = tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := New()
			err := c.ConfigureTransport(tt.Transport)
			if err != nil {
				t.Fatal(err)
			}

			_, err = c.ResizeVolume(context.Background(), "bananaNS", "bananaVolume", size.GiB, tt.params)
			if err != nil {
				t.Errorf("Resize returns error: %q", err)
			}

			if !reflect.DeepEqual(tt.Transport.ResizeVolumeGotParams, tt.wantParams) {
				t.Errorf("ResizeVolume() got = %+v, want %+v", tt.Transport.ResizeVolumeGotParams, tt.wantParams)
			}
		})
	}
}
