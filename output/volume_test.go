package output

import (
	"reflect"
	"testing"
	"time"

	"github.com/kr/pretty"

	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/pkg/health"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/labels"
	"code.storageos.net/storageos/c2-cli/pkg/version"
	"code.storageos.net/storageos/c2-cli/volume"
)

func TestNewVolume(t *testing.T) {
	t.Parallel()

	labelsFromPairs := func(t *testing.T, pairs ...string) labels.Set {
		set, err := labels.NewSetFromPairs(pairs)
		if err != nil {
			t.Errorf("failed to set up test case: %v", err)
		}
		return set
	}

	tests := []struct {
		name string

		vol   *volume.Resource
		ns    *namespace.Resource
		nodes map[id.Node]*node.Resource

		wantOutputVol *Volume
		wantErr       error
	}{
		{
			name: "ok when master with nil sync progress",

			vol: &volume.Resource{
				ID:          "vol-id",
				Name:        "vol-name",
				Description: "vol-description",
				AttachedOn:  "attached-node",

				Namespace:  "namespace-id",
				Labels:     labelsFromPairs(t, "a=b", "b=c"),
				Filesystem: volume.FsTypeFromString("BLOCK"),
				SizeBytes:  42,

				Master: &volume.Deployment{
					ID:           "deploy-id",
					Node:         "node-id",
					Health:       health.MasterOnline,
					Promotable:   true,
					SyncProgress: nil, // explicitly nil
				},

				Replicas:  nil,
				CreatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				Version:   version.FromString("some-version"),
			},
			ns: &namespace.Resource{
				Name: "namespace-name",
			},
			nodes: map[id.Node]*node.Resource{
				"attached-node": &node.Resource{
					Name: "attached-node-name",
				},
				"node-id": &node.Resource{
					Name: "node-name",
				},
			},

			wantOutputVol: &Volume{
				ID:             "vol-id",
				Name:           "vol-name",
				Description:    "vol-description",
				AttachedOn:     "attached-node",
				AttachedOnName: "attached-node-name",
				Namespace:      "namespace-id",
				NamespaceName:  "namespace-name",
				Labels:         labelsFromPairs(t, "a=b", "b=c"),
				Filesystem:     volume.FsTypeFromString("BLOCK"),
				SizeBytes:      42,
				Master: &Deployment{
					ID:           "deploy-id",
					Node:         "node-id",
					NodeName:     "node-name",
					Health:       health.MasterOnline,
					Promotable:   true,
					SyncProgress: nil,
				},
				Replicas:  []*Deployment{},
				CreatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				Version:   version.FromString("some-version"),
			},
			wantErr: nil,
		},
		{
			name: "ok when replicas both with sync progress and without",

			vol: &volume.Resource{
				ID:          "vol-id",
				Name:        "vol-name",
				Description: "vol-description",
				AttachedOn:  "attached-node",

				Namespace:  "namespace-id",
				Labels:     labelsFromPairs(t, "b=b", "a=c"),
				Filesystem: volume.FsTypeFromString("BLOCK"),
				SizeBytes:  42,

				Master: &volume.Deployment{
					ID:         "deploy-id",
					Node:       "node-id",
					Health:     health.MasterOnline,
					Promotable: true,
				},

				Replicas: []*volume.Deployment{
					{
						ID:           "repl-1",
						Node:         "node-1",
						Health:       health.ReplicaReady,
						Promotable:   true,
						SyncProgress: nil,
					},
					{
						ID:         "repl-2",
						Node:       "node-2",
						Health:     health.ReplicaSyncing,
						Promotable: false,
						SyncProgress: &volume.SyncProgress{
							BytesRemaining:            6,
							ThroughputBytes:           4,
							EstimatedSecondsRemaining: 2,
						},
					},
				},
				CreatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				Version:   version.FromString("some-version"),
			},
			ns: &namespace.Resource{
				Name: "namespace-name",
			},
			nodes: map[id.Node]*node.Resource{
				"attached-node": &node.Resource{
					Name: "attached-node-name",
				},
				"node-id": &node.Resource{
					Name: "node-name",
				},
				"node-1": &node.Resource{
					Name: "node-1-name",
				},
				"node-2": &node.Resource{
					Name: "node-2-name",
				},
			},

			wantOutputVol: &Volume{
				ID:             "vol-id",
				Name:           "vol-name",
				Description:    "vol-description",
				AttachedOn:     "attached-node",
				AttachedOnName: "attached-node-name",
				Namespace:      "namespace-id",
				NamespaceName:  "namespace-name",
				Labels:         labelsFromPairs(t, "b=b", "a=c"),
				Filesystem:     volume.FsTypeFromString("BLOCK"),
				SizeBytes:      42,
				Master: &Deployment{
					ID:         "deploy-id",
					Node:       "node-id",
					NodeName:   "node-name",
					Health:     health.MasterOnline,
					Promotable: true,
				},
				Replicas: []*Deployment{
					{
						ID:           "repl-1",
						Node:         "node-1",
						NodeName:     "node-1-name",
						Health:       health.ReplicaReady,
						Promotable:   true,
						SyncProgress: nil,
					},
					{
						ID:         "repl-2",
						Node:       "node-2",
						NodeName:   "node-2-name",
						Health:     health.ReplicaSyncing,
						Promotable: false,
						SyncProgress: &SyncProgress{
							BytesRemaining:            6,
							ThroughputBytes:           4,
							EstimatedSecondsRemaining: 2,
						},
					},
				},
				CreatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				Version:   version.FromString("some-version"),
			},
			wantErr: nil,
		},
		{
			name: "missing attached on node information",

			vol: &volume.Resource{
				ID:          "vol-id",
				Name:        "vol-name",
				Description: "vol-description",
				AttachedOn:  "attached-node",

				Namespace:  "namespace-id",
				Labels:     labelsFromPairs(t, "b=b", "a=c"),
				Filesystem: volume.FsTypeFromString("BLOCK"),
				SizeBytes:  42,

				Master: &volume.Deployment{
					ID:         "deploy-id",
					Node:       "node-id",
					Health:     health.MasterOnline,
					Promotable: true,
				},

				Replicas:  []*volume.Deployment{},
				CreatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				Version:   version.FromString("some-version"),
			},
			ns: &namespace.Resource{
				Name: "namespace-name",
			},
			nodes: map[id.Node]*node.Resource{
				"node-id": &node.Resource{
					Name: "node-name",
				},
			},

			wantOutputVol: nil,
			wantErr:       NewMissingRequiredNodeErr("attached-node"),
		},
	}

	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotOutputVol, gotErr := NewVolume(tt.vol, tt.ns, tt.nodes)
			if gotErr != tt.wantErr {
				t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
			}

			if !reflect.DeepEqual(gotOutputVol, tt.wantOutputVol) {
				pretty.Ldiff(t, gotOutputVol, tt.wantOutputVol)
				t.Errorf("got output vol %v, want %v", pretty.Sprint(gotOutputVol), pretty.Sprint(tt.wantOutputVol))
			}
		})
	}
}

func TestNewDeployment(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string

		inDeployment *volume.Deployment
		nodes        map[id.Node]*node.Resource

		wantOutputDeployment *Deployment
		wantErr              error
	}{
		{
			name: "ok",

			inDeployment: &volume.Deployment{
				ID:         "id",
				Node:       "node-id",
				Health:     "health",
				Promotable: true,
			},
			nodes: map[id.Node]*node.Resource{
				"node-id": &node.Resource{
					Name: "node-name",
				},
			},

			wantOutputDeployment: &Deployment{
				ID:         "id",
				Node:       "node-id",
				NodeName:   "node-name",
				Health:     "health",
				Promotable: true,
			},
			wantErr: nil,
		},
		{
			name: "missing",

			inDeployment: &volume.Deployment{
				ID:         "id",
				Node:       "node-id",
				Health:     "health",
				Promotable: true,
			},
			nodes: map[id.Node]*node.Resource{
				"some-other-node-id": &node.Resource{
					Name: "node-name",
				},
			},

			wantOutputDeployment: nil,
			wantErr:              NewMissingRequiredNodeErr("node-id"),
		},
	}

	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotOutputDeployment, gotErr := newDeployment(tt.inDeployment, tt.nodes)
			if gotErr != tt.wantErr {
				t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
			}

			if !reflect.DeepEqual(gotOutputDeployment, tt.wantOutputDeployment) {
				pretty.Ldiff(t, gotOutputDeployment, tt.wantOutputDeployment)
				t.Errorf("got output %v, want %v", gotOutputDeployment, tt.wantOutputDeployment)
			}
		})
	}
}
