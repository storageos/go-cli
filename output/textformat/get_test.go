package textformat

import (
	"bytes"
	"context"
	"testing"
	"time"

	"code.storageos.net/storageos/c2-cli/cluster"
	"code.storageos.net/storageos/c2-cli/licence"
	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/output"
	"code.storageos.net/storageos/c2-cli/pkg/health"
	"code.storageos.net/storageos/c2-cli/pkg/labels"
	"code.storageos.net/storageos/c2-cli/pkg/size"
	"code.storageos.net/storageos/c2-cli/policygroup"
	"code.storageos.net/storageos/c2-cli/volume"
)

func TestDisplayer_GetCluster(t *testing.T) {
	t.Parallel()

	var mockTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		timeHumanizer output.TimeHumanizer
		resource      *cluster.Resource
		nodes         []*node.Resource
		wantW         string
		wantErr       bool
	}{
		{
			name: "print cluster",
			resource: &cluster.Resource{
				ID:        "bananaCluster",
				CreatedAt: mockTime,
			},
			nodes: []*node.Resource{
				{
					ID:     "bananaNodeID",
					Name:   "bananaNodeName",
					IOAddr: "127.0.0.1",
					Health: health.NodeOnline,
				},
				{
					ID:     "kiwiNodeID",
					Name:   "kiwiNodeName",
					IOAddr: "127.0.0.2",
					Health: health.NodeOnline,
				},
				{
					ID:     "pearNodeID",
					Name:   "pearNodeName",
					IOAddr: "127.0.0.3",
					Health: health.NodeOffline,
				},
			},
			wantW: `ID:           bananaCluster                      
Created at:   2000-01-01T00:00:00Z (xx aeons ago)
Updated at:   0001-01-01T00:00:00Z (xx aeons ago)
Nodes:        3                                  
  Healthy:    2                                  
  Unhealthy:  1                                  
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

			err := d.GetCluster(context.Background(), w, output.NewCluster(tt.resource, tt.nodes))
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

func TestDisplayer_GetLicence(t *testing.T) {
	t.Parallel()

	var mockTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		timeHumanizer output.TimeHumanizer
		resource      *licence.Resource
		wantW         string
		wantErr       bool
	}{
		{
			name: "print licence",
			resource: &licence.Resource{
				ClusterID:            "bananaID",
				ExpiresAt:            mockTime,
				ClusterCapacityBytes: 42 * size.GiB,
				UsedBytes:            42 / 2 * size.GiB,
				Kind:                 "bananaKind",
				Features:             []string{"nfs", "banana"},
				CustomerName:         "bananaCustomer",
			},
			wantW: `ClusterID:      bananaID                           
Expiration:     2000-01-01T00:00:00Z (xx aeons ago)
Capacity:       42 GiB (45097156608)               
Used:           21 GiB (22548578304)               
Kind:           bananaKind                         
Features:       [banana nfs]                       
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

			err := d.GetLicence(context.Background(), w, output.NewLicence(tt.resource))
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLicence() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("GetLicence() gotW = \n%v\n, want \n%v\n", gotW, tt.wantW)
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
					ID:        "kiwiID",
					Name:      "kiwiName",
					Labels:    map[string]string{"kiwiKey": "kiwiValue"},
					CreatedAt: mockTime,
					UpdatedAt: mockTime,
					Version:   "42",
				},
			},
			wantW: `NAME           AGE         
bananaName     xx aeons ago
pineappleName  xx aeons ago
kiwiName       xx aeons ago
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

func TestDisplayer_GetPolicyGroup(t *testing.T) {
	t.Parallel()

	var mockTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		resource   *policygroup.Resource
		namespaces []*namespace.Resource
		wantW      string
		wantErr    bool
	}{
		{
			name: "print policy group",
			resource: &policygroup.Resource{
				ID:   "bananaID",
				Name: "bananaName",
				Users: []*policygroup.Member{
					{
						ID:       "member1",
						Username: "memberName1",
					},
					{
						ID:       "member2",
						Username: "memberName2",
					},
				},
				Specs: []*policygroup.Spec{
					{
						NamespaceID:  "bananaNamespace1",
						ResourceType: "*",
						ReadOnly:     false,
					},
					{
						NamespaceID:  "bananaNamespace2",
						ResourceType: "volume",
						ReadOnly:     true,
					},
				},
				CreatedAt: mockTime,
				UpdatedAt: mockTime,
				Version:   "42",
			},
			namespaces: []*namespace.Resource{},
			wantW: `NAME        USERS  SPECS  AGE         
bananaName  2      2      xx aeons ago
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
			err := d.GetPolicyGroup(context.Background(), w, output.NewPolicyGroup(tt.resource, tt.namespaces))
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPolicyGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("GetPolicyGroup() gotW = \n%v\n, want \n%v\n", gotW, tt.wantW)
			}
		})
	}
}

func TestDisplayer_GetListPolicyGroups(t *testing.T) {
	t.Parallel()

	var mockTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		resources  []*policygroup.Resource
		namespaces []*namespace.Resource
		wantW      string
		wantErr    bool
	}{
		{
			name: "print policy groups",
			resources: []*policygroup.Resource{
				{
					ID:   "bananaID",
					Name: "bananaName",
					Users: []*policygroup.Member{
						{
							ID:       "member1",
							Username: "memberName1",
						},
						{
							ID:       "member2",
							Username: "memberName2",
						},
					},
					Specs: []*policygroup.Spec{
						{
							NamespaceID:  "bananaNamespace1",
							ResourceType: "*",
							ReadOnly:     false,
						},
						{
							NamespaceID:  "bananaNamespace2",
							ResourceType: "volume",
							ReadOnly:     true,
						},
					},
					CreatedAt: mockTime,
					UpdatedAt: mockTime,
					Version:   "42",
				},
				{
					ID:   "pineappleID",
					Name: "pineappleName",
					Users: []*policygroup.Member{
						{
							ID:       "member1",
							Username: "memberName1",
						},
					},
					Specs: []*policygroup.Spec{
						{
							NamespaceID:  "pineappleNamespace1",
							ResourceType: "*",
							ReadOnly:     false,
						},
						{
							NamespaceID:  "pineappleNamespace2",
							ResourceType: "volume",
							ReadOnly:     true,
						},
						{
							NamespaceID:  "pineappleNamespace3",
							ResourceType: "volume",
							ReadOnly:     true,
						},
					},
					CreatedAt: mockTime,
					UpdatedAt: mockTime,
					Version:   "42",
				},
				{
					ID:   "kiwiID",
					Name: "kiwiName",
					Users: []*policygroup.Member{
						{
							ID:       "member1",
							Username: "memberName1",
						},
						{
							ID:       "member2",
							Username: "memberName2",
						},
					},
					Specs: []*policygroup.Spec{
						{
							NamespaceID:  "kiwiNamespace1",
							ResourceType: "*",
							ReadOnly:     false,
						},
						{
							NamespaceID:  "kiwiNamespace2",
							ResourceType: "volume",
							ReadOnly:     true,
						},
					},
					CreatedAt: mockTime,
					UpdatedAt: mockTime,
					Version:   "42",
				},
			},
			namespaces: []*namespace.Resource{},
			wantW: `NAME           USERS  SPECS  AGE         
bananaName     2      2      xx aeons ago
pineappleName  1      3      xx aeons ago
kiwiName       2      2      xx aeons ago
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
			err := d.GetListPolicyGroups(context.Background(), w, output.NewPolicyGroups(tt.resources, tt.namespaces))
			if (err != nil) != tt.wantErr {
				t.Errorf("GetListPolicyGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("GetListPolicyGroups() gotW = \n%v\n, want \n%v\n", gotW, tt.wantW)
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
			name: "print volume uses replica length for total count if bad label",
			volume: &output.Volume{
				ID:             "bananaID",
				Name:           "banana-name",
				Description:    "banana description",
				AttachedOnName: "banana-node-a",
				Namespace:      "banana-namespace",
				NamespaceName:  "kiwi",
				Labels: labels.Set{
					"kiwi":               "42",
					"pear":               "42",
					volume.LabelReplicas: "NaN",
				},
				Filesystem: volume.FsTypeFromString("ext4"),
				SizeBytes:  size.GiB,
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
		{
			name: "print volume uses replica label number for total count",
			volume: &output.Volume{
				ID:             "bananaID",
				Name:           "banana-name",
				Description:    "banana description",
				AttachedOnName: "banana-node-a",
				Namespace:      "banana-namespace",
				NamespaceName:  "kiwi",
				Labels: labels.Set{
					"kiwi":               "42",
					"pear":               "42",
					volume.LabelReplicas: "3",
				},
				Filesystem: volume.FsTypeFromString("ext4"),
				SizeBytes:  size.GiB,
				Master: &output.Deployment{
					ID:         "bananaDeploymentID1",
					NodeName:   "banana-node1",
					Health:     "ready",
					Promotable: true,
				},
				Replicas:  []*output.Deployment{},
				CreatedAt: mockTime,
				UpdatedAt: mockTime,
				Version:   "42",
			},
			wantW: `NAMESPACE  NAME         SIZE     LOCATION              ATTACHED ON    REPLICAS  AGE         
kiwi       banana-name  1.0 GiB  banana-node1 (ready)  banana-node-a  0/3       xx aeons ago
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
					SizeBytes:  size.GiB,
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
					SizeBytes:  2 * size.GiB,
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
