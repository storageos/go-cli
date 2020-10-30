package openapi

import (
	"reflect"
	"testing"
	"time"

	"github.com/kr/pretty"

	"code.storageos.net/storageos/c2-cli/cluster"
	"code.storageos.net/storageos/c2-cli/licence"
	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/pkg/health"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/labels"
	"code.storageos.net/storageos/c2-cli/pkg/version"
	"code.storageos.net/storageos/c2-cli/policygroup"
	"code.storageos.net/storageos/c2-cli/user"
	"code.storageos.net/storageos/c2-cli/volume"
	"code.storageos.net/storageos/openapi"
)

func TestDecodeLicence(t *testing.T) {
	t.Parallel()

	mockExpiryTime := time.Date(2020, 01, 01, 0, 0, 0, 2, time.UTC)

	tests := []struct {
		name string

		model openapi.Licence

		wantResource *licence.Resource
		wantErr      error
	}{
		{
			name: "ok",

			model: openapi.Licence{
				ClusterID:            "bananas",
				ExpiresAt:            mockExpiryTime,
				ClusterCapacityBytes: 42,
				UsedBytes:            42 / 2,
				Kind:                 "mockLicence",
				CustomerName:         "go testing framework",
				Features:             &[]string{"nfs"},
				Version:              "bananaVersion",
			},

			wantResource: &licence.Resource{
				ClusterID:            "bananas",
				ExpiresAt:            mockExpiryTime,
				ClusterCapacityBytes: 42,
				UsedBytes:            42 / 2,
				Kind:                 "mockLicence",
				CustomerName:         "go testing framework",
				Features:             []string{"nfs"},
				Version:              version.Version("bananaVersion"),
			},
			wantErr: nil,
		},
		{
			name: "nil features",

			model: openapi.Licence{
				ClusterID:            "bananas",
				ExpiresAt:            mockExpiryTime,
				ClusterCapacityBytes: 42,
				UsedBytes:            42 / 2,
				Kind:                 "mockLicence",
				CustomerName:         "go testing framework",
				Features:             nil,
				Version:              "bananaVersion",
			},

			wantResource: &licence.Resource{
				ClusterID:            "bananas",
				ExpiresAt:            mockExpiryTime,
				ClusterCapacityBytes: 42,
				UsedBytes:            42 / 2,
				Kind:                 "mockLicence",
				CustomerName:         "go testing framework",
				Features:             []string{},
				Version:              version.Version("bananaVersion"),
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := &codec{}

			gotResource, gotErr := c.decodeLicence(tt.model)

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
			}

			if !reflect.DeepEqual(gotResource, tt.wantResource) {
				pretty.Ldiff(t, gotResource, tt.wantResource)
				t.Errorf("got decoded cluster config %v, want %v", pretty.Sprint(gotResource), pretty.Sprint(tt.wantResource))
			}
		})
	}
}

func TestDecodeCluster(t *testing.T) {
	t.Parallel()

	mockCreatedAtTime := time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC)
	mockUpdatedAtTime := time.Date(2020, 01, 01, 0, 0, 0, 1, time.UTC)

	tests := []struct {
		name string

		model openapi.Cluster

		wantResource *cluster.Resource
		wantErr      error
	}{
		{
			name: "ok",

			model: openapi.Cluster{
				Id: "bananas",

				DisableTelemetry:      true,
				DisableCrashReporting: true,
				DisableVersionCheck:   true,
				LogLevel:              openapi.LOGLEVEL_DEBUG,
				LogFormat:             openapi.LOGFORMAT_JSON,
				CreatedAt:             mockCreatedAtTime,
				UpdatedAt:             mockUpdatedAtTime,
				Version:               "NDIK",
			},

			wantResource: &cluster.Resource{
				ID: "bananas",

				DisableTelemetry:      true,
				DisableCrashReporting: true,
				DisableVersionCheck:   true,

				LogLevel:  cluster.LogLevelFromString("debug"),
				LogFormat: cluster.LogFormatFromString("json"),

				CreatedAt: mockCreatedAtTime,
				UpdatedAt: mockUpdatedAtTime,
				Version:   "NDIK",
			},
			wantErr: nil,
		},
		{
			name:         "does not panic with no fields",
			model:        openapi.Cluster{},
			wantResource: &cluster.Resource{},
			wantErr:      nil,
		},
	}

	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := &codec{}

			gotResource, gotErr := c.decodeCluster(tt.model)

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
			}

			if !reflect.DeepEqual(gotResource, tt.wantResource) {
				pretty.Ldiff(t, gotResource, tt.wantResource)
				t.Errorf("got decoded cluster config %v, want %v", pretty.Sprint(gotResource), pretty.Sprint(tt.wantResource))
			}
		})
	}
}

func TestDecodeNode(t *testing.T) {
	t.Parallel()

	mockCreatedAtTime := time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC)
	mockUpdatedAtTime := time.Date(2020, 01, 01, 0, 0, 0, 1, time.UTC)

	tests := []struct {
		name string

		model openapi.Node

		wantResource *node.Resource
		wantErr      error
	}{
		{
			name: "ok",

			model: openapi.Node{
				Id:                 "banananodeid",
				Name:               "banananodename",
				Health:             openapi.NODEHEALTH_ONLINE,
				IoEndpoint:         "arbitraryIOEndpoint",
				SupervisorEndpoint: "arbitrarySupervisorEndpoint",
				GossipEndpoint:     "arbitraryGossipEndpoint",
				ClusteringEndpoint: "arbitraryClusteringEndpoint",
				Labels: map[string]string{
					"storageos.com/label": "value",
				},
				CreatedAt: mockCreatedAtTime,
				UpdatedAt: mockUpdatedAtTime,
				Version:   "NDIK",
			},

			wantResource: &node.Resource{
				ID:     "banananodeid",
				Name:   "banananodename",
				Health: health.NodeOnline,

				Labels: labels.Set{
					"storageos.com/label": "value",
				},

				IOAddr:         "arbitraryIOEndpoint",
				SupervisorAddr: "arbitrarySupervisorEndpoint",
				GossipAddr:     "arbitraryGossipEndpoint",
				ClusteringAddr: "arbitraryClusteringEndpoint",

				CreatedAt: mockCreatedAtTime,
				UpdatedAt: mockUpdatedAtTime,
				Version:   "NDIK",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := &codec{}

			gotResource, gotErr := c.decodeNode(tt.model)

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
			}

			if !reflect.DeepEqual(gotResource, tt.wantResource) {
				pretty.Ldiff(t, gotResource, tt.wantResource)
				t.Errorf("got decoded node config %v, want %v", pretty.Sprint(gotResource), pretty.Sprint(tt.wantResource))
			}
		})
	}
}

func TestDecodeVolume(t *testing.T) {
	t.Parallel()

	mockCreatedAtTime := time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC)
	mockUpdatedAtTime := time.Date(2020, 01, 01, 0, 0, 0, 1, time.UTC)

	tests := []struct {
		name string

		model openapi.Volume

		wantResource *volume.Resource
		wantErr      error
	}{
		{
			name: "ok with replicas",

			model: openapi.Volume{
				Id:          "my-volume-id",
				Name:        "my-volume",
				Description: "some arbitrary description",
				AttachedOn:  "some-arbitrary-node-id",
				Nfs: openapi.NfsConfig{
					Exports: &[]openapi.NfsExportConfig{
						{
							ExportID:   1,
							Path:       "/",
							PseudoPath: "/",
							Acls: []openapi.NfsAcl{
								{
									Identity: openapi.NfsAclIdentity{
										IdentityType: "cidr",
										Matcher:      "10.0.0.0/8",
									},
									SquashConfig: openapi.NfsAclSquashConfig{
										Gid:    0,
										Uid:    0,
										Squash: "root",
									},
									AccessLevel: "rw",
								},
							},
						},
					},
					ServiceEndpoint: func(s string) *string { return &s }("10.0.0.1:/"),
				},
				NamespaceID: "some-arbitrary-namespace-id",
				Labels: map[string]string{
					"storageos.com/label": "value",
				},
				FsType: openapi.FSTYPE_EXT4,
				Master: openapi.MasterDeploymentInfo{
					Id:         "master-id",
					NodeID:     "some-arbitrary-node-id",
					Health:     openapi.MASTERHEALTH_ONLINE,
					Promotable: true,
				},
				Replicas: &[]openapi.ReplicaDeploymentInfo{
					{
						Id:         "replica-a-id",
						NodeID:     "some-second-node-id",
						Health:     openapi.REPLICAHEALTH_SYNCING,
						Promotable: false,
					},
					{
						Id:         "replica-b-id",
						NodeID:     "some-third-node-id",
						Health:     openapi.REPLICAHEALTH_READY,
						Promotable: true,
					},
				},
				SizeBytes: 1337,
				CreatedAt: mockCreatedAtTime,
				UpdatedAt: mockUpdatedAtTime,
				Version:   "NDIK",
			},

			wantResource: &volume.Resource{
				ID:          "my-volume-id",
				Name:        "my-volume",
				Description: "some arbitrary description",
				SizeBytes:   1337,

				AttachedOn: "some-arbitrary-node-id",
				Nfs: volume.NFSConfig{
					Exports: []volume.NFSExportConfig{
						{
							ExportID:   1,
							Path:       "/",
							PseudoPath: "/",
							ACLs: []volume.NFSExportConfigACL{
								{
									Identity: volume.NFSExportConfigACLIdentity{
										IdentityType: "cidr",
										Matcher:      "10.0.0.0/8",
									},
									SquashConfig: volume.NFSExportConfigACLSquashConfig{
										GID:    0,
										UID:    0,
										Squash: "root",
									},
									AccessLevel: "rw",
								},
							},
						},
					},
					ServiceEndpoint: "10.0.0.1:/",
				},
				Namespace: "some-arbitrary-namespace-id",
				Labels: labels.Set{
					"storageos.com/label": "value",
				},
				Filesystem: volume.FsTypeFromString("ext4"),

				Master: &volume.Deployment{
					ID:         "master-id",
					Node:       "some-arbitrary-node-id",
					Health:     health.MasterOnline,
					Promotable: true,
				},
				Replicas: []*volume.Deployment{
					{
						ID:         "replica-a-id",
						Node:       "some-second-node-id",
						Health:     health.ReplicaSyncing,
						Promotable: false,
					},
					{
						ID:         "replica-b-id",
						Node:       "some-third-node-id",
						Health:     health.ReplicaReady,
						Promotable: true,
					},
				},

				CreatedAt: mockCreatedAtTime,
				UpdatedAt: mockUpdatedAtTime,
				Version:   "NDIK",
			},
			wantErr: nil,
		},
		{
			name: "ok no replicas",

			model: openapi.Volume{
				Id:          "my-volume-id",
				Name:        "my-volume",
				Description: "some arbitrary description",
				AttachedOn:  "some-arbitrary-node-id",
				Nfs: openapi.NfsConfig{
					Exports: &[]openapi.NfsExportConfig{
						{
							ExportID:   1,
							Path:       "/",
							PseudoPath: "/",
							Acls: []openapi.NfsAcl{
								{
									Identity: openapi.NfsAclIdentity{
										IdentityType: "cidr",
										Matcher:      "10.0.0.0/8",
									},
									SquashConfig: openapi.NfsAclSquashConfig{
										Gid:    0,
										Uid:    0,
										Squash: "root",
									},
									AccessLevel: "rw",
								},
							},
						},
					},
					ServiceEndpoint: func(s string) *string { return &s }("10.0.0.1:/"),
				},
				NamespaceID: "some-arbitrary-namespace-id",
				Labels: map[string]string{
					"storageos.com/label": "value",
				},
				FsType: openapi.FSTYPE_EXT4,
				Master: openapi.MasterDeploymentInfo{
					Id:         "master-id",
					NodeID:     "some-arbitrary-node-id",
					Health:     openapi.MASTERHEALTH_ONLINE,
					Promotable: true,
				},
				SizeBytes: 1337,
				CreatedAt: mockCreatedAtTime,
				UpdatedAt: mockUpdatedAtTime,
				Version:   "NDIK",
			},

			wantResource: &volume.Resource{
				ID:          "my-volume-id",
				Name:        "my-volume",
				Description: "some arbitrary description",
				SizeBytes:   1337,

				AttachedOn: "some-arbitrary-node-id",
				Nfs: volume.NFSConfig{
					Exports: []volume.NFSExportConfig{
						{
							ExportID:   1,
							Path:       "/",
							PseudoPath: "/",
							ACLs: []volume.NFSExportConfigACL{
								{
									Identity: volume.NFSExportConfigACLIdentity{
										IdentityType: "cidr",
										Matcher:      "10.0.0.0/8",
									},
									SquashConfig: volume.NFSExportConfigACLSquashConfig{
										GID:    0,
										UID:    0,
										Squash: "root",
									},
									AccessLevel: "rw",
								},
							},
						},
					},
					ServiceEndpoint: "10.0.0.1:/",
				},
				Namespace: "some-arbitrary-namespace-id",
				Labels: labels.Set{
					"storageos.com/label": "value",
				},
				Filesystem: volume.FsTypeFromString("ext4"),

				Master: &volume.Deployment{
					ID:         "master-id",
					Node:       "some-arbitrary-node-id",
					Health:     health.MasterOnline,
					Promotable: true,
				},
				Replicas: []*volume.Deployment{},

				CreatedAt: mockCreatedAtTime,
				UpdatedAt: mockUpdatedAtTime,
				Version:   "NDIK",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := &codec{}

			gotResource, gotErr := c.decodeVolume(tt.model)

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
			}

			if !reflect.DeepEqual(gotResource, tt.wantResource) {
				pretty.Ldiff(t, gotResource, tt.wantResource)
				t.Errorf("got decoded volume config %v, want %v", pretty.Sprint(gotResource), pretty.Sprint(tt.wantResource))
			}
		})
	}
}

func TestDecodeNamespace(t *testing.T) {
	t.Parallel()

	mockCreatedAtTime := time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC)
	mockUpdatedAtTime := time.Date(2020, 01, 01, 0, 0, 0, 1, time.UTC)

	tests := []struct {
		name string

		model openapi.Namespace

		wantResource *namespace.Resource
		wantErr      error
	}{
		{
			name: "ok",

			model: openapi.Namespace{
				Id:   "my-namespace-id",
				Name: "my-namespace",
				Labels: map[string]string{
					"storageos.com/label": "value",
				},
				CreatedAt: mockCreatedAtTime,
				UpdatedAt: mockUpdatedAtTime,
				Version:   "NDIK",
			},

			wantResource: &namespace.Resource{
				ID:   "my-namespace-id",
				Name: "my-namespace",
				Labels: labels.Set{
					"storageos.com/label": "value",
				},

				CreatedAt: mockCreatedAtTime,
				UpdatedAt: mockUpdatedAtTime,
				Version:   "NDIK",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := &codec{}

			gotResource, gotErr := c.decodeNamespace(tt.model)

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
			}

			if !reflect.DeepEqual(gotResource, tt.wantResource) {
				pretty.Ldiff(t, gotResource, tt.wantResource)
				t.Errorf("got decoded namespace config %v, want %v", pretty.Sprint(gotResource), pretty.Sprint(tt.wantResource))
			}
		})
	}
}

func TestDecodePolicyGroup(t *testing.T) {
	t.Parallel()

	mockCreatedAtTime := time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC)
	mockUpdatedAtTime := time.Date(2020, 01, 01, 0, 0, 0, 1, time.UTC)

	tests := []struct {
		name string

		model openapi.PolicyGroup

		wantResource *policygroup.Resource
		wantErr      error
	}{
		{
			name: "ok with users and specs",

			model: openapi.PolicyGroup{
				Id:   "policy-group-id",
				Name: "policy-group-name",
				Users: []openapi.PolicyGroupUsers{
					{
						Id:       "user-id",
						Username: "username",
					},
					{
						Id:       "user-id-2",
						Username: "username-2",
					},
				},
				Specs: &[]openapi.PoliciesIdSpecs{
					{
						NamespaceID:  "namespace-id",
						ResourceType: "resource-type",
						ReadOnly:     true,
					},
				},
				CreatedAt: mockCreatedAtTime,
				UpdatedAt: mockUpdatedAtTime,
				Version:   "version",
			},

			wantResource: &policygroup.Resource{
				ID:   "policy-group-id",
				Name: "policy-group-name",
				Users: []*policygroup.Member{
					{
						ID:       "user-id",
						Username: "username",
					},
					{
						ID:       "user-id-2",
						Username: "username-2",
					},
				},
				Specs: []*policygroup.Spec{
					{
						NamespaceID:  "namespace-id",
						ResourceType: "resource-type",
						ReadOnly:     true,
					},
				},
				CreatedAt: mockCreatedAtTime,
				UpdatedAt: mockUpdatedAtTime,
				Version:   "version",
			},
			wantErr: nil,
		},
		{
			name: "ok with no users or specs",

			model: openapi.PolicyGroup{
				Id:        "policy-group-id",
				Name:      "policy-group-name",
				Users:     nil,
				Specs:     nil,
				CreatedAt: mockCreatedAtTime,
				UpdatedAt: mockUpdatedAtTime,
				Version:   "version",
			},

			wantResource: &policygroup.Resource{
				ID:        "policy-group-id",
				Name:      "policy-group-name",
				Users:     []*policygroup.Member{},
				Specs:     []*policygroup.Spec{},
				CreatedAt: mockCreatedAtTime,
				UpdatedAt: mockUpdatedAtTime,
				Version:   "version",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := codec{}

			gotResource, gotErr := c.decodePolicyGroup(tt.model)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
			}

			if !reflect.DeepEqual(gotResource, tt.wantResource) {
				pretty.Ldiff(t, gotResource, tt.wantResource)
				t.Errorf("got decoded policy group %v, want %v", pretty.Sprint(gotResource), pretty.Sprint(tt.wantResource))
			}
		})
	}
}

func TestDecodeUser(t *testing.T) {
	t.Parallel()

	mockCreatedAtTime := time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC)
	mockUpdatedAtTime := time.Date(2020, 01, 01, 0, 0, 0, 1, time.UTC)

	tests := []struct {
		name string

		model openapi.User

		wantResource *user.Resource
		wantErr      error
	}{
		{
			name: "ok with groups",

			model: openapi.User{
				Id:       "my-user-id",
				Username: "my-username",
				IsAdmin:  true,
				Groups: &[]string{
					"group-a-id",
					"group-b-id",
				},
				CreatedAt: mockCreatedAtTime,
				UpdatedAt: mockUpdatedAtTime,
				Version:   "NDIK",
			},

			wantResource: &user.Resource{
				ID:       "my-user-id",
				Username: "my-username",

				IsAdmin: true,
				Groups: []id.PolicyGroup{
					"group-a-id",
					"group-b-id",
				},

				CreatedAt: mockCreatedAtTime,
				UpdatedAt: mockUpdatedAtTime,
				Version:   "NDIK",
			},
			wantErr: nil,
		},
		{
			name: "ok no groups",

			model: openapi.User{
				Id:        "my-user-id",
				Username:  "my-username",
				IsAdmin:   true,
				Groups:    nil,
				CreatedAt: mockCreatedAtTime,
				UpdatedAt: mockUpdatedAtTime,
				Version:   "NDIK",
			},

			wantResource: &user.Resource{
				ID:       "my-user-id",
				Username: "my-username",

				IsAdmin: true,
				Groups:  []id.PolicyGroup{},

				CreatedAt: mockCreatedAtTime,
				UpdatedAt: mockUpdatedAtTime,
				Version:   "NDIK",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := &codec{}

			gotResource, gotErr := c.decodeUser(tt.model)

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
			}

			if !reflect.DeepEqual(gotResource, tt.wantResource) {
				pretty.Ldiff(t, gotResource, tt.wantResource)
				t.Errorf("got decoded user config %v, want %v", pretty.Sprint(gotResource), pretty.Sprint(tt.wantResource))
			}
		})
	}
}
