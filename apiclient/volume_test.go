package apiclient

import (
	"context"
	"errors"
	"reflect"
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

		configProvider *mockConfigProvider
		transport      *mockTransport

		volumeName string

		wantResource *volume.Resource
		wantErr      error
	}{
		{
			name: "ok",

			configProvider: &mockConfigProvider{},
			transport: &mockTransport{
				ListVolumesResource: []*volume.Resource{
					&volume.Resource{
						Name: "possibly-arthur",
					},
					&volume.Resource{
						Name: "definitely-alan",
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

			configProvider: &mockConfigProvider{},
			transport: &mockTransport{
				ListVolumesResource: []*volume.Resource{
					&volume.Resource{
						Name: "possibly-arthur",
					},
					&volume.Resource{
						Name: "not-alan",
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

			configProvider: &mockConfigProvider{},
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

			client := New(tt.configProvider)
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
				},
				ListVolumesResource: []*volume.Resource{
					{
						ID: "volume-1",
					},
					{
						ID: "volume-2",
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
				ListVolumesResource: []*volume.Resource{
					{
						ID: "volume-1",
					},
					{
						ID: "volume-2",
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

			client := New(&mockConfigProvider{})
			if err := client.ConfigureTransport(tt.transport); err != nil {
				t.Fatalf("got error configuring client transport: %v", err)
			}

			gotVolumes, gotErr := client.fetchAllVolumes(context.Background())
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
			name: "dont filter when no names provided",

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
			name: "dont filter when no uids provided",

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
