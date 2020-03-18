package textformat

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/dustin/go-humanize"

	"code.storageos.net/storageos/c2-cli/cluster"
	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/output"
	"code.storageos.net/storageos/c2-cli/pkg/labels"
	"code.storageos.net/storageos/c2-cli/volume"
)

type mockTimeFormatter struct {
	Str string
}

func (m *mockTimeFormatter) TimeToHuman(t time.Time) string {
	return m.Str
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
Created At:       2000-01-01T00:00:00Z (xx aeons ago)
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

			err := d.GetCluster(context.Background(), w, tt.resource)
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
			err := d.GetNode(context.Background(), w, tt.resource)
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
			err := d.GetListNodes(context.Background(), w, tt.resources)
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
			err := d.GetNamespace(context.Background(), w, tt.resource)
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
			err := d.GetListNamespaces(context.Background(), w, tt.resources)
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
