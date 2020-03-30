package textformat

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/kr/pretty"

	"code.storageos.net/storageos/c2-cli/cluster"
	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/output"
	"code.storageos.net/storageos/c2-cli/pkg/capacity"
	"code.storageos.net/storageos/c2-cli/pkg/health"
	"code.storageos.net/storageos/c2-cli/pkg/labels"
	"code.storageos.net/storageos/c2-cli/volume"
)

type mockTimeFormatter struct {
	Str string
}

func (m *mockTimeFormatter) TimeToHuman(t time.Time) string {
	return m.Str
}

func TestCreateUser(t *testing.T) {
	t.Parallel()

	var mockTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		user    *output.User
		wantW   string
		wantErr error
	}{
		{
			name: "display created user with single group ok",
			user: &output.User{
				ID:       "bananaID",
				Username: "banana-name",

				IsAdmin: true,
				Groups: []output.PolicyGroup{
					{
						ID:   "policy-group-id",
						Name: "policy-group-name",
					},
				},

				CreatedAt: mockTime,
				UpdatedAt: mockTime,
				Version:   "42",
			},
			wantW: `NAME         ROLE   AGE                GROUPS           
banana-name  admin  donkeys years ago  policy-group-name
`,
			wantErr: nil,
		},
		{
			name: "display created user with multiple groups ok",
			user: &output.User{
				ID:       "bananaID",
				Username: "banana-name",

				IsAdmin: true,
				Groups: []output.PolicyGroup{
					{
						ID:   "policy-group-id",
						Name: "policy-group-name",
					},
					{
						ID:   "policy-group-id-2",
						Name: "policy-group-name-2",
					},
				},

				CreatedAt: mockTime,
				UpdatedAt: mockTime,
				Version:   "42",
			},
			wantW: `NAME         ROLE   AGE                GROUPS                               
banana-name  admin  donkeys years ago  policy-group-name,policy-group-name-2
`,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := NewDisplayer(&mockTimeFormatter{Str: "donkeys years ago"})
			w := &bytes.Buffer{}

			gotErr := d.CreateUser(context.Background(), w, tt.user)
			if gotErr != tt.wantErr {
				t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
			}

			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("got output: \n%v\nwant: \n%v\n", gotW, tt.wantW)
			}
		})
	}

}

func TestDisplayer_GetCluster(t *testing.T) {
	t.Parallel()

	var mockTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		timeHumanizer output.TimeHumanizer
		resource      *cluster.Resource
		wantW         string
		wantErr       bool
	}{
		{
			name: "print cluster",
			resource: &cluster.Resource{
				ID: "bananaCluster",
				Licence: &cluster.Licence{
					ClusterID:            "bananaCluster",
					ExpiresAt:            mockTime,
					ClusterCapacityBytes: 42 * humanize.GiByte,
					Kind:                 "bananaLicence",
					CustomerName:         "bananaCustomer",
				},
				CreatedAt: mockTime,
			},
			wantW: `ID:               bananaCluster                      
Licence:                                             
  expiration:     2000-01-01T00:00:00Z (xx aeons ago)
  capacity:       42 GiB                             
  kind:           bananaLicence                      
  customer name:  bananaCustomer                     
Created at:       2000-01-01T00:00:00Z (xx aeons ago)
Updated at:       0001-01-01T00:00:00Z (xx aeons ago)
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		var tt = tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := NewDisplayer(&mockTimeFormatter{Str: "xx aeons ago"})
			w := &bytes.Buffer{}

			err := d.GetCluster(context.Background(), w, output.NewCluster(tt.resource))
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCluster() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("GetCluster() gotW = \n%v\n, want \n%v\n", gotW, tt.wantW)
			}
		})
	}
}

func TestDisplayer_UpdateLicence(t *testing.T) {
	t.Parallel()

	var mockTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		timeHumanizer output.TimeHumanizer
		licence       *cluster.Licence
		wantW         string
		wantErr       bool
	}{
		{
			name: "print cluster",
			licence: &cluster.Licence{
				ClusterID:            "bananaCluster",
				ExpiresAt:            mockTime,
				ClusterCapacityBytes: 42 * humanize.GiByte,
				Kind:                 "bananaLicence",
				CustomerName:         "bananaCustomer",
			},
			wantW: `Licence applied to cluster bananaCluster.

Expiration:     2000-01-01T00:00:00Z (xx aeons ago)
Capacity:       42 GiB                             
Kind:           bananaLicence                      
Customer name:  bananaCustomer                     
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		var tt = tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := NewDisplayer(&mockTimeFormatter{Str: "xx aeons ago"})
			w := &bytes.Buffer{}

			err := d.UpdateLicence(context.Background(), w, output.NewLicence(tt.licence))
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCluster() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("GetCluster() gotW = \n%v\n, want \n%v\n", gotW, tt.wantW)
			}
		})
	}
}

func TestDisplayer_GetUser(t *testing.T) {
	t.Parallel()

	var mockTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		user    *output.User
		wantW   string
		wantErr error
	}{
		{
			name: "display user with single group ok",
			user: &output.User{
				ID:       "bananaID",
				Username: "banana-name",

				IsAdmin: true,
				Groups: []output.PolicyGroup{
					{
						ID:   "policy-group-id",
						Name: "policy-group-name",
					},
				},

				CreatedAt: mockTime,
				UpdatedAt: mockTime,
				Version:   "42",
			},
			wantW: `NAME         ROLE   AGE                GROUPS           
banana-name  admin  donkeys years ago  policy-group-name
`,
			wantErr: nil,
		},
		{
			name: "display user with multiple groups ok",
			user: &output.User{
				ID:       "bananaID",
				Username: "banana-name",

				IsAdmin: true,
				Groups: []output.PolicyGroup{
					{
						ID:   "policy-group-id",
						Name: "policy-group-name",
					},
					{
						ID:   "policy-group-id-2",
						Name: "policy-group-name-2",
					},
				},

				CreatedAt: mockTime,
				UpdatedAt: mockTime,
				Version:   "42",
			},
			wantW: `NAME         ROLE   AGE                GROUPS                               
banana-name  admin  donkeys years ago  policy-group-name,policy-group-name-2
`,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := NewDisplayer(&mockTimeFormatter{Str: "donkeys years ago"})
			w := &bytes.Buffer{}

			gotErr := d.GetUser(context.Background(), w, tt.user)
			if gotErr != tt.wantErr {
				t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
			}

			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("got output: \n%v\nwant: \n%v\n", gotW, tt.wantW)
			}
		})
	}
}

func TestDisplayer_GetUsers(t *testing.T) {
	t.Parallel()

	var mockTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		user    []*output.User
		wantW   string
		wantErr error
	}{
		{
			name: "display users with single group ok",
			user: []*output.User{
				{
					ID:       "bananaID",
					Username: "banana-name",

					IsAdmin: true,
					Groups: []output.PolicyGroup{
						{
							ID:   "policy-group-id",
							Name: "policy-group-name",
						},
					},

					CreatedAt: mockTime,
					UpdatedAt: mockTime,
					Version:   "42",
				},
				{
					ID:       "kiwiID",
					Username: "kiwi-name",

					IsAdmin: true,
					Groups: []output.PolicyGroup{
						{
							ID:   "policy-group-id",
							Name: "policy-group-name",
						},
					},

					CreatedAt: mockTime,
					UpdatedAt: mockTime,
					Version:   "42",
				},
				{
					ID:       "pearID",
					Username: "pear-name",

					IsAdmin: true,
					Groups: []output.PolicyGroup{
						{
							ID:   "policy-group-id",
							Name: "policy-group-name",
						},
					},

					CreatedAt: mockTime,
					UpdatedAt: mockTime,
					Version:   "42",
				},
			},
			wantW: `NAME         ROLE   AGE                GROUPS           
banana-name  admin  donkeys years ago  policy-group-name
kiwi-name    admin  donkeys years ago  policy-group-name
pear-name    admin  donkeys years ago  policy-group-name
`,
			wantErr: nil,
		},
		{
			name: "display users with multiple groups ok",
			user: []*output.User{
				{
					ID:       "bananaID",
					Username: "banana-name",

					IsAdmin: true,
					Groups: []output.PolicyGroup{
						{
							ID:   "policy-group-id",
							Name: "policy-group-name",
						},
						{
							ID:   "policy-group-id-2",
							Name: "policy-group-name-2",
						},
					},

					CreatedAt: mockTime,
					UpdatedAt: mockTime,
					Version:   "42",
				},
				{
					ID:       "kiwiID",
					Username: "kiwi-name",

					IsAdmin: true,
					Groups: []output.PolicyGroup{
						{
							ID:   "policy-group-id",
							Name: "policy-group-name",
						},
						{
							ID:   "policy-group-id-2",
							Name: "policy-group-name-2",
						},
					},

					CreatedAt: mockTime,
					UpdatedAt: mockTime,
					Version:   "42",
				},
				{
					ID:       "pearID",
					Username: "pear-name",

					IsAdmin: true,
					Groups: []output.PolicyGroup{
						{
							ID:   "policy-group-id",
							Name: "policy-group-name",
						},
						{
							ID:   "policy-group-id-2",
							Name: "policy-group-name-2",
						},
					},

					CreatedAt: mockTime,
					UpdatedAt: mockTime,
					Version:   "42",
				},
			},
			wantW: `NAME         ROLE   AGE                GROUPS                               
banana-name  admin  donkeys years ago  policy-group-name,policy-group-name-2
kiwi-name    admin  donkeys years ago  policy-group-name,policy-group-name-2
pear-name    admin  donkeys years ago  policy-group-name,policy-group-name-2
`,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := NewDisplayer(&mockTimeFormatter{Str: "donkeys years ago"})
			w := &bytes.Buffer{}

			gotErr := d.GetUsers(context.Background(), w, tt.user)
			if gotErr != tt.wantErr {
				t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
			}

			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("got output: \n%v\nwant: \n%v\n", gotW, tt.wantW)
			}
		})
	}
}

func TestDisplayer_GetNode(t *testing.T) {
	t.Parallel()

	var mockTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		resource *node.Resource
		wantW    string
		wantErr  bool
	}{
		{
			name: "print single node",
			resource: &node.Resource{
				ID:     "bananaID",
				Name:   "bananaName",
				Health: "ready",
				Labels: map[string]string{
					"bananaLabelKey": "bananaLabelValue",
				},
				CreatedAt: mockTime,
				UpdatedAt: mockTime,
				Version:   "42",
			},
			wantW: `NAME        HEALTH  AGE           LABELS                         
bananaName  ready   xx aeons ago  bananaLabelKey=bananaLabelValue
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		var tt = tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := NewDisplayer(&mockTimeFormatter{Str: "xx aeons ago"})
			w := &bytes.Buffer{}
			err := d.GetNode(context.Background(), w, output.NewNode(tt.resource))
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("GetNode() gotW = \n%v\n, want \n%v\n", gotW, tt.wantW)
			}
		})
	}
}

func TestDisplayer_GetListNodes(t *testing.T) {
	t.Parallel()

	var mockTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		resources []*node.Resource
		wantW     string
		wantErr   bool
	}{
		{
			name: "print nodes",
			resources: []*node.Resource{
				{
					ID:     "bananaID",
					Name:   "bananaName",
					Health: "ready",
					Labels: map[string]string{
						"bananaLabelKey": "bananaLabelValue",
					},
					CreatedAt: mockTime,
					UpdatedAt: mockTime,
					Version:   "42",
				},
				{
					ID:     "kiwiID",
					Name:   "kiwiName",
					Health: "ready",
					Labels: map[string]string{
						"kiwiLabelKey": "kiwiLabelValue",
					},
					CreatedAt: mockTime,
					UpdatedAt: mockTime,
					Version:   "42",
				},
				{
					ID:     "pineappleID",
					Name:   "pineappleName",
					Health: "offline",
					Labels: map[string]string{
						"pineappleLabelKey": "pineappleLabelValue",
					},
					CreatedAt: mockTime,
					UpdatedAt: mockTime,
					Version:   "42",
				},
			},
			wantW: `NAME           HEALTH   AGE           LABELS                               
bananaName     ready    xx aeons ago  bananaLabelKey=bananaLabelValue      
kiwiName       ready    xx aeons ago  kiwiLabelKey=kiwiLabelValue          
pineappleName  offline  xx aeons ago  pineappleLabelKey=pineappleLabelValue
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			d := NewDisplayer(&mockTimeFormatter{Str: "xx aeons ago"})
			w := &bytes.Buffer{}
			err := d.GetListNodes(context.Background(), w, output.NewNodes(tt.resources))
			if (err != nil) != tt.wantErr {
				t.Errorf("GetListNodes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("GetListNodes() gotW = \n%v\n, want \n%v\n", gotW, tt.wantW)
			}
		})
	}
}

func TestDisplayer_GetNamespace(t *testing.T) {
	t.Parallel()

	var mockTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		resource *namespace.Resource
		wantW    string
		wantErr  bool
	}{
		{
			name: "print namespace",
			resource: &namespace.Resource{
				ID:        "bananaID",
				Name:      "bananaName",
				Labels:    map[string]string{"bananaKey": "bananaValue"},
				CreatedAt: mockTime,
				UpdatedAt: mockTime,
				Version:   "42",
			},
			wantW: `NAME        AGE         
bananaName  xx aeons ago
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := NewDisplayer(&mockTimeFormatter{Str: "xx aeons ago"})
			w := &bytes.Buffer{}
			err := d.GetNamespace(context.Background(), w, output.NewNamespace(tt.resource))
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNamespace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("GetNamespace() gotW = \n%v\n, want \n%v\n", gotW, tt.wantW)
			}
		})
	}
}

func TestDisplayer_GetListNamespaces(t *testing.T) {
	t.Parallel()

	var mockTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		resources []*namespace.Resource
		wantW     string
		wantErr   bool
	}{
		{
			name: "print namespaces",
			resources: []*namespace.Resource{
				{
					ID:        "bananaID",
					Name:      "bananaName",
					Labels:    map[string]string{"bananaKey": "bananaValue"},
					CreatedAt: mockTime,
					UpdatedAt: mockTime,
					Version:   "42",
				},
				{
					ID:        "pineappleID",
					Name:      "pineappleName",
					Labels:    map[string]string{"pineappleKey": "pineappleValue"},
					CreatedAt: mockTime,
					UpdatedAt: mockTime,
					Version:   "42",
				},
				{
					ID:        "kiwiaID",
					Name:      "kiwiaName",
					Labels:    map[string]string{"kiwiaKey": "kiwiaValue"},
					CreatedAt: mockTime,
					UpdatedAt: mockTime,
					Version:   "42",
				},
			},
			wantW: `NAME           AGE         
bananaName     xx aeons ago
pineappleName  xx aeons ago
kiwiaName      xx aeons ago
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		var tt = tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			d := NewDisplayer(&mockTimeFormatter{Str: "xx aeons ago"})
			w := &bytes.Buffer{}
			err := d.GetListNamespaces(context.Background(), w, output.NewNamespaces(tt.resources))
			if (err != nil) != tt.wantErr {
				t.Errorf("GetListNamespaces() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("GetListNamespaces() gotW = \n%v\n, want \n%v\n", gotW, tt.wantW)
			}
		})
	}
}

func TestDisplayer_GetVolume(t *testing.T) {
	t.Parallel()

	var mockTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		volume        *output.Volume
		namespaceName string
		wantW         string
		wantErr       bool
	}{
		{
			name: "print volume",
			volume: &output.Volume{
				ID:             "bananaID",
				Name:           "banana-name",
				Description:    "banana description",
				AttachedOnName: "banana-node-a",
				Namespace:      "banana-namespace",
				NamespaceName:  "kiwi",
				Labels: labels.Set{
					"kiwi": "42",
					"pear": "42",
				},
				Filesystem: volume.FsTypeFromString("ext4"),
				SizeBytes:  humanize.GiByte,
				Master: &output.Deployment{
					ID:         "bananaDeploymentID1",
					NodeName:   "banana-node1",
					Health:     "ready",
					Promotable: true,
				},
				Replicas: []*output.Deployment{
					{
						ID:         "bananaDeploymentID2",
						NodeName:   "banana-node2",
						Health:     "ready",
						Promotable: true,
					},
					{
						ID:         "bananaDeploymentID3",
						NodeName:   "banana-node3",
						Health:     "offline",
						Promotable: false,
					},
					{
						ID:         "bananaDeploymentID4",
						NodeName:   "banana-node4",
						Health:     "ready",
						Promotable: true,
					},
				},
				CreatedAt: mockTime,
				UpdatedAt: mockTime,
				Version:   "42",
			},
			wantW: `NAMESPACE  NAME         SIZE     LOCATION              ATTACHED ON    REPLICAS  AGE         
kiwi       banana-name  1.0 GiB  banana-node1 (ready)  banana-node-a  2/3       xx aeons ago
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := NewDisplayer(&mockTimeFormatter{Str: "xx aeons ago"})
			w := &bytes.Buffer{}

			err := d.GetVolume(context.Background(), w, tt.volume)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVolume() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("GetVolume() \nGOT: \n%v\n\nWANT: \n%v\n", gotW, tt.wantW)
			}
		})
	}
}

func TestDisplayer_GetListVolumes(t *testing.T) {
	t.Parallel()

	var mockTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		volumes []*output.Volume
		wantW   string
		wantErr bool
	}{
		{
			name: "print volumes",
			volumes: []*output.Volume{
				{
					ID:             "bananaID",
					Name:           "banana-name",
					Description:    "banana description",
					AttachedOnName: "banana-node-a",
					Namespace:      "banana-namespace",
					NamespaceName:  "BANANA",
					Labels: labels.Set{
						"kiwi": "42",
						"pear": "42",
					},
					Filesystem: volume.FsTypeFromString("ext4"),
					SizeBytes:  humanize.GiByte,
					Master: &output.Deployment{
						ID:         "bananaDeploymentID1",
						NodeName:   "banana-node1",
						Health:     "ready",
						Promotable: true,
					},
					Replicas: []*output.Deployment{
						{
							ID:         "bananaDeploymentID2",
							NodeName:   "banana-node2",
							Health:     "ready",
							Promotable: true,
						},
						{
							ID:         "bananaDeploymentID3",
							NodeName:   "banana-node3",
							Health:     "offline",
							Promotable: false,
						},
						{
							ID:         "bananaDeploymentID4",
							NodeName:   "banana-node4",
							Health:     "ready",
							Promotable: true,
						},
					},
					CreatedAt: mockTime,
					UpdatedAt: mockTime,
					Version:   "42",
				},
				{
					ID:             "kiwiID",
					Name:           "kiwi-name",
					Description:    "kiwi description",
					AttachedOnName: "kiwi-node-a",
					Namespace:      "kiwi-namespace",
					NamespaceName:  "KIWI",
					Labels: labels.Set{
						"kiwi": "foo",
						"pear": "bar",
					},
					Filesystem: volume.FsTypeFromString("ext4"),
					SizeBytes:  2 * humanize.GiByte,
					Master: &output.Deployment{
						ID:         "bananaDeploymentID1",
						NodeName:   "banana-node1",
						Health:     "ready",
						Promotable: true,
					},
					Replicas: []*output.Deployment{
						{
							ID:         "bananaDeploymentID2",
							NodeName:   "banana-node2",
							Health:     "ready",
							Promotable: true,
						},
						{
							ID:         "bananaDeploymentID3",
							NodeName:   "banana-node3",
							Health:     "offline",
							Promotable: false,
						},
						{
							ID:         "bananaDeploymentID4",
							NodeName:   "banana-node4",
							Health:     "ready",
							Promotable: true,
						},
					},
					CreatedAt: mockTime,
					UpdatedAt: mockTime,
					Version:   "43",
				},
			},
			wantW: `NAMESPACE  NAME         SIZE     LOCATION              ATTACHED ON    REPLICAS  AGE         
BANANA     banana-name  1.0 GiB  banana-node1 (ready)  banana-node-a  2/3       xx aeons ago
KIWI       kiwi-name    2.0 GiB  banana-node1 (ready)  kiwi-node-a    2/3       xx aeons ago
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := NewDisplayer(&mockTimeFormatter{Str: "xx aeons ago"})
			w := &bytes.Buffer{}

			err := d.GetListVolumes(context.Background(), w, tt.volumes)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetListVolumes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("GetListVolumes() GOT:\n%v\n\nWANT:\n%v\n", gotW, tt.wantW)
			}
		})
	}
}

func TestDescribeNode(t *testing.T) {
	t.Parallel()

	var (
		createdTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedTime = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
	)

	tests := []struct {
		name string

		node *output.NodeDescription

		wantOutput string
		wantErr    error
	}{
		{
			name: "describe node with local volumes ok",

			node: &output.NodeDescription{
				Node: output.Node{
					ID:     "bananaID",
					Name:   "bananaName",
					Health: health.NodeOnline,
					Capacity: capacity.Stats{
						Total:     humanize.GiByte,
						Available: 500 * humanize.MiByte,
						Free:      550 * humanize.MiByte,
					},
					IOAddr:         "io.addr",
					SupervisorAddr: "supervisor.addr",
					GossipAddr:     "gossip.addr",
					ClusteringAddr: "clustering.addr",
					Labels: labels.Set{
						"a": "b",
						"1": "2",
					},
					CreatedAt: createdTime,
					UpdatedAt: updatedTime,
					Version:   "some-version",
				},
				HostedVolumes: []*output.HostedVolume{
					{
						ID:            "volumeID",
						Name:          "volumeName",
						Description:   "description",
						Namespace:     "namespaceID",
						NamespaceName: "namespaceName",
						Labels:        labels.Set{},
						Filesystem:    volume.FsTypeFromString("fs"),
						SizeBytes:     humanize.GiByte,
						LocalDeployment: output.LocalDeployment{
							ID:         "deployID",
							Kind:       "master",
							Health:     health.MasterOnline,
							Promotable: true,
						},
						CreatedAt: createdTime,
						UpdatedAt: updatedTime,
						Version:   "other-version",
					},
					{
						ID:            "volumeID2",
						Name:          "volumeName2",
						Description:   "description2",
						Namespace:     "namespaceID2",
						NamespaceName: "namespaceName2",
						Labels:        labels.Set{},
						Filesystem:    volume.FsTypeFromString("fs"),
						SizeBytes:     2 * humanize.GiByte,
						LocalDeployment: output.LocalDeployment{
							ID:         "deployID2",
							Kind:       "replica",
							Health:     health.ReplicaReady,
							Promotable: true,
						},
						CreatedAt: createdTime,
						UpdatedAt: updatedTime,
						Version:   "another-version",
					},
				},
			},

			wantOutput: `ID                         bananaID                              
Name                       bananaName                            
Health                     online                                
Addresses:               
  Data Transfer address    io.addr                               
  Gossip address           gossip.addr                           
  Supervisor address       supervisor.addr                       
  Clustering address       clustering.addr                       
Labels                     1=2,a=b                               
Created at                 2000-01-01T00:00:00Z (a long time ago)
Updated at                 2001-01-01T00:00:00Z (a long time ago)
Version                    some-version                          
Available capacity         500 MiB/1.0 GiB (474 MiB in use)      

Local volume deployments:
  DEPLOYMENT ID            VOLUME                                  NAMESPACE       HEALTH  TYPE     SIZE   
  deployID                 volumeName                              namespaceName   online  master   1.0 GiB
  deployID2                volumeName2                             namespaceName2  ready   replica  2.0 GiB
`,
			wantErr: nil,
		},
		{
			name: "describe node with no volumes ok",

			node: &output.NodeDescription{
				Node: output.Node{
					ID:     "bananaID",
					Name:   "bananaName",
					Health: health.NodeOnline,
					Capacity: capacity.Stats{
						Total:     humanize.GiByte,
						Available: 500 * humanize.MiByte,
						Free:      550 * humanize.MiByte,
					},
					IOAddr:         "io.addr",
					SupervisorAddr: "supervisor.addr",
					GossipAddr:     "gossip.addr",
					ClusteringAddr: "clustering.addr",
					Labels: labels.Set{
						"a": "b",
						"1": "2",
					},
					CreatedAt: createdTime,
					UpdatedAt: updatedTime,
					Version:   "some-version",
				},
				HostedVolumes: nil,
			},

			wantOutput: `ID                       bananaID                              
Name                     bananaName                            
Health                   online                                
Addresses:             
  Data Transfer address  io.addr                               
  Gossip address         gossip.addr                           
  Supervisor address     supervisor.addr                       
  Clustering address     clustering.addr                       
Labels                   1=2,a=b                               
Created at               2000-01-01T00:00:00Z (a long time ago)
Updated at               2001-01-01T00:00:00Z (a long time ago)
Version                  some-version                          
Available capacity       500 MiB/1.0 GiB (474 MiB in use)      
`,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := NewDisplayer(&mockTimeFormatter{Str: "a long time ago"})
			w := &bytes.Buffer{}

			gotErr := d.DescribeNode(context.Background(), w, tt.node)
			if gotErr != tt.wantErr {
				t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
			}

			gotOutput := w.String()
			if gotOutput != tt.wantOutput {
				pretty.Ldiff(t, gotOutput, tt.wantOutput)
				t.Errorf("got output: \n%v\n\nwant: \n%v\n", gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestDisplayer_DescribeVolume(t *testing.T) {
	t.Parallel()

	var mockTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name string

		volume *output.Volume

		wantOutput string
		wantErr    error
	}{
		{
			name: "describe volume ok",

			volume: &output.Volume{
				ID:             "bananaID",
				Name:           "banana-name",
				Description:    "banana description",
				AttachedOn:     "banana-node-a-id",
				AttachedOnName: "banana-node-a",
				Namespace:      "banana-namespace-id",
				NamespaceName:  "banana-namespace",
				Labels: labels.Set{
					"kiwi": "42",
					"pear": "42",
				},
				Filesystem: volume.FsTypeFromString("ext4"),
				SizeBytes:  humanize.GiByte,
				Master: &output.Deployment{
					ID:         "bananaDeploymentID1",
					Node:       "banana-node1-id",
					NodeName:   "banana-node1",
					Health:     "ready",
					Promotable: true,
				},
				Replicas: []*output.Deployment{
					{
						ID:           "bananaDeploymentID2",
						Node:         "banana-node2-id",
						NodeName:     "banana-node2",
						Health:       health.ReplicaReady,
						Promotable:   true,
						SyncProgress: nil,
					},
					{
						ID:           "bananaDeploymentID3",
						Node:         "banana-node3-id",
						NodeName:     "banana-node3",
						Health:       health.ReplicaDeleted,
						Promotable:   false,
						SyncProgress: nil,
					},
					{
						ID:         "bananaDeploymentID4",
						Node:       "banana-node4-id",
						NodeName:   "banana-node4",
						Health:     health.ReplicaSyncing,
						Promotable: false,
						SyncProgress: &output.SyncProgress{
							BytesRemaining:            256 * humanize.MiByte,
							ThroughputBytes:           0,
							EstimatedSecondsRemaining: 1200,
						},
					},
				},
				CreatedAt: mockTime,
				UpdatedAt: mockTime,
				Version:   "42",
			},

			wantOutput: `ID               bananaID                                                                        
Name             banana-name                                                                     
Description      banana description                                                              
AttachedOn       banana-node-a (banana-node-a-id)                                                
Namespace        banana-namespace (banana-namespace-id)                                          
Labels           kiwi=42,pear=42                                                                 
Filesystem       ext4                                                                            
Size             1.0 GiB (1073741824 bytes)                                                      
Version          42                                                                              
Created at       2000-01-01T00:00:00Z (xx aeons ago)                                             
Updated at       2000-01-01T00:00:00Z (xx aeons ago)                                             
                                                                                                 
Master:        
  ID             bananaDeploymentID1                                                             
  Node           banana-node1 (banana-node1-id)                                                  
  Health         ready                                                                           
                                                                                                 
Replicas:      
  ID             bananaDeploymentID2                                                             
  Node           banana-node2 (banana-node2-id)                                                  
  Health         ready                                                                           
  Promotable     true                                                                            
                                                                                                 
  ID             bananaDeploymentID3                                                             
  Node           banana-node3 (banana-node3-id)                                                  
  Health         deleted                                                                         
  Promotable     false                                                                           
                                                                                                 
  ID             bananaDeploymentID4                                                             
  Node           banana-node4 (banana-node4-id)                                                  
  Health         syncing                                                                         
  Promotable     false                                                                           
  Sync Progress  768.00 MiB / 1.00 GiB [##########################........] 75.00%  -  ETA: 20m0s
`,
			wantErr: nil,
		},
		{
			name: "describe volume with no replicas ok",

			volume: &output.Volume{
				ID:             "bananaID",
				Name:           "banana-name",
				Description:    "banana description",
				AttachedOn:     "banana-node-a-id",
				AttachedOnName: "banana-node-a",
				Namespace:      "banana-namespace-id",
				NamespaceName:  "banana-namespace",
				Labels: labels.Set{
					"kiwi": "42",
					"pear": "42",
				},
				Filesystem: volume.FsTypeFromString("ext4"),
				SizeBytes:  humanize.GiByte,
				Master: &output.Deployment{
					ID:         "bananaDeploymentID1",
					Node:       "banana-node1-id",
					NodeName:   "banana-node1",
					Health:     "ready",
					Promotable: true,
				},
				Replicas:  []*output.Deployment{},
				CreatedAt: mockTime,
				UpdatedAt: mockTime,
				Version:   "42",
			},

			wantOutput: `ID           bananaID                              
Name         banana-name                           
Description  banana description                    
AttachedOn   banana-node-a (banana-node-a-id)      
Namespace    banana-namespace (banana-namespace-id)
Labels       kiwi=42,pear=42                       
Filesystem   ext4                                  
Size         1.0 GiB (1073741824 bytes)            
Version      42                                    
Created at   2000-01-01T00:00:00Z (xx aeons ago)   
Updated at   2000-01-01T00:00:00Z (xx aeons ago)   
                                                   
Master:    
  ID         bananaDeploymentID1                   
  Node       banana-node1 (banana-node1-id)        
  Health     ready                                 
`,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := NewDisplayer(&mockTimeFormatter{Str: "xx aeons ago"})
			w := &bytes.Buffer{}

			gotErr := d.DescribeVolume(context.Background(), w, tt.volume)
			if gotErr != tt.wantErr {
				t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
				return
			}

			gotOutput := w.String()
			if gotOutput != tt.wantOutput {
				pretty.Ldiff(t, gotOutput, tt.wantOutput)
				t.Errorf("got output: \n%v\n\nwant: \n%v\n", gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestDisplayer_DescribeListVolumes(t *testing.T) {
	t.Parallel()

	var mockTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name string

		volumes []*output.Volume

		wantOutput string
		wantErr    error
	}{
		{
			name: "describe volume ok",

			volumes: []*output.Volume{
				{
					ID:             "bananaID",
					Name:           "banana-name",
					Description:    "banana description",
					AttachedOn:     "banana-node-a-id",
					AttachedOnName: "banana-node-a",
					Namespace:      "banana-namespace-id",
					NamespaceName:  "banana-namespace",
					Labels: labels.Set{
						"kiwi": "42",
						"pear": "42",
					},
					Filesystem: volume.FsTypeFromString("ext4"),
					SizeBytes:  humanize.GiByte,
					Master: &output.Deployment{
						ID:         "bananaDeploymentID1",
						Node:       "banana-node1-id",
						NodeName:   "banana-node1",
						Health:     "ready",
						Promotable: true,
					},
					Replicas: []*output.Deployment{
						{
							ID:           "bananaDeploymentID2",
							Node:         "banana-node2-id",
							NodeName:     "banana-node2",
							Health:       health.ReplicaReady,
							Promotable:   true,
							SyncProgress: nil,
						},
						{
							ID:           "bananaDeploymentID3",
							Node:         "banana-node3-id",
							NodeName:     "banana-node3",
							Health:       health.ReplicaDeleted,
							Promotable:   false,
							SyncProgress: nil,
						},
						{
							ID:         "bananaDeploymentID4",
							Node:       "banana-node4-id",
							NodeName:   "banana-node4",
							Health:     health.ReplicaSyncing,
							Promotable: false,
							SyncProgress: &output.SyncProgress{
								BytesRemaining:            256 * humanize.MiByte,
								ThroughputBytes:           0,
								EstimatedSecondsRemaining: 1200,
							},
						},
					},
					CreatedAt: mockTime,
					UpdatedAt: mockTime,
					Version:   "42",
				},
				{
					ID:             "kiwiID",
					Name:           "kiwi-name",
					Description:    "kiwi description",
					AttachedOn:     "kiwi-node-a-id",
					AttachedOnName: "kiwi-node-a",
					Namespace:      "kiwi-namespace-id",
					NamespaceName:  "kiwi-namespace",
					Labels: labels.Set{
						"kiwi": "42",
						"pear": "42",
					},
					Filesystem: volume.FsTypeFromString("ext4"),
					SizeBytes:  humanize.GiByte,
					Master: &output.Deployment{
						ID:         "kiwiDeploymentID1",
						Node:       "kiwi-node1-id",
						NodeName:   "kiwi-node1",
						Health:     "ready",
						Promotable: true,
					},
					Replicas: []*output.Deployment{
						{
							ID:           "kiwiDeploymentID2",
							Node:         "kiwi-node2-id",
							NodeName:     "kiwi-node2",
							Health:       health.ReplicaReady,
							Promotable:   true,
							SyncProgress: nil,
						},
						{
							ID:           "kiwiDeploymentID3",
							Node:         "kiwi-node3-id",
							NodeName:     "kiwi-node3",
							Health:       health.ReplicaDeleted,
							Promotable:   false,
							SyncProgress: nil,
						},
						{
							ID:         "kiwiDeploymentID4",
							Node:       "kiwi-node4-id",
							NodeName:   "kiwi-node4",
							Health:     health.ReplicaSyncing,
							Promotable: false,
							SyncProgress: &output.SyncProgress{
								BytesRemaining:            256 * humanize.MiByte,
								ThroughputBytes:           0,
								EstimatedSecondsRemaining: 1200,
							},
						},
					},
					CreatedAt: mockTime,
					UpdatedAt: mockTime,
					Version:   "42",
				},
			},

			wantOutput: `ID               bananaID                                                                        
Name             banana-name                                                                     
Description      banana description                                                              
AttachedOn       banana-node-a (banana-node-a-id)                                                
Namespace        banana-namespace (banana-namespace-id)                                          
Labels           kiwi=42,pear=42                                                                 
Filesystem       ext4                                                                            
Size             1.0 GiB (1073741824 bytes)                                                      
Version          42                                                                              
Created at       2000-01-01T00:00:00Z (xx aeons ago)                                             
Updated at       2000-01-01T00:00:00Z (xx aeons ago)                                             
                                                                                                 
Master:        
  ID             bananaDeploymentID1                                                             
  Node           banana-node1 (banana-node1-id)                                                  
  Health         ready                                                                           
                                                                                                 
Replicas:      
  ID             bananaDeploymentID2                                                             
  Node           banana-node2 (banana-node2-id)                                                  
  Health         ready                                                                           
  Promotable     true                                                                            
                                                                                                 
  ID             bananaDeploymentID3                                                             
  Node           banana-node3 (banana-node3-id)                                                  
  Health         deleted                                                                         
  Promotable     false                                                                           
                                                                                                 
  ID             bananaDeploymentID4                                                             
  Node           banana-node4 (banana-node4-id)                                                  
  Health         syncing                                                                         
  Promotable     false                                                                           
  Sync Progress  768.00 MiB / 1.00 GiB [##########################........] 75.00%  -  ETA: 20m0s

ID               kiwiID                                                                          
Name             kiwi-name                                                                       
Description      kiwi description                                                                
AttachedOn       kiwi-node-a (kiwi-node-a-id)                                                    
Namespace        kiwi-namespace (kiwi-namespace-id)                                              
Labels           kiwi=42,pear=42                                                                 
Filesystem       ext4                                                                            
Size             1.0 GiB (1073741824 bytes)                                                      
Version          42                                                                              
Created at       2000-01-01T00:00:00Z (xx aeons ago)                                             
Updated at       2000-01-01T00:00:00Z (xx aeons ago)                                             
                                                                                                 
Master:        
  ID             kiwiDeploymentID1                                                               
  Node           kiwi-node1 (kiwi-node1-id)                                                      
  Health         ready                                                                           
                                                                                                 
Replicas:      
  ID             kiwiDeploymentID2                                                               
  Node           kiwi-node2 (kiwi-node2-id)                                                      
  Health         ready                                                                           
  Promotable     true                                                                            
                                                                                                 
  ID             kiwiDeploymentID3                                                               
  Node           kiwi-node3 (kiwi-node3-id)                                                      
  Health         deleted                                                                         
  Promotable     false                                                                           
                                                                                                 
  ID             kiwiDeploymentID4                                                               
  Node           kiwi-node4 (kiwi-node4-id)                                                      
  Health         syncing                                                                         
  Promotable     false                                                                           
  Sync Progress  768.00 MiB / 1.00 GiB [##########################........] 75.00%  -  ETA: 20m0s
`,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := NewDisplayer(&mockTimeFormatter{Str: "xx aeons ago"})
			w := &bytes.Buffer{}

			gotErr := d.DescribeListVolumes(context.Background(), w, tt.volumes)
			if gotErr != tt.wantErr {
				t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
				return
			}

			gotOutput := w.String()
			if gotOutput != tt.wantOutput {
				pretty.Ldiff(t, gotOutput, tt.wantOutput)
				t.Errorf("got output: \n%v\n\nwant: \n%v\n", gotOutput, tt.wantOutput)
			}
		})
	}
}
