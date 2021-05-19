package textformat

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/kr/pretty"

	"code.storageos.net/storageos/c2-cli/cluster"
	"code.storageos.net/storageos/c2-cli/licence"
	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/output"
	"code.storageos.net/storageos/c2-cli/pkg/capacity"
	"code.storageos.net/storageos/c2-cli/pkg/health"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/labels"
	"code.storageos.net/storageos/c2-cli/pkg/size"
	"code.storageos.net/storageos/c2-cli/policygroup"
	"code.storageos.net/storageos/c2-cli/user"
	"code.storageos.net/storageos/c2-cli/volume"
)

func TestDisplayer_DescribeNode(t *testing.T) {
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
						Total: size.GiB,
						Free:  550 * size.MiB,
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
						SizeBytes:     size.GiB,
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
						SizeBytes:     2 * size.GiB,
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
Labels                     1=2,                                  
                           a=b                                   
Created at                 2000-01-01T00:00:00Z (a long time ago)
Updated at                 2001-01-01T00:00:00Z (a long time ago)
Version                    some-version                          
Available capacity         1.0 GiB (474 MiB in use)              

Local volume deployments:
  NAMESPACE                VOLUME                                  DEPLOYMENT ID  HEALTH  TYPE     SIZE   
  namespaceName            volumeName                              deployID       online  master   1.0 GiB
  namespaceName2           volumeName2                             deployID2      ready   replica  2.0 GiB
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
						Total: size.GiB,
						Free:  550 * size.MiB,
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
Labels                   1=2,                                  
                         a=b                                   
Created at               2000-01-01T00:00:00Z (a long time ago)
Updated at               2001-01-01T00:00:00Z (a long time ago)
Version                  some-version                          
Available capacity       1.0 GiB (474 MiB in use)              
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
				compareOutput(t, gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestDisplayer_DescribeCluster(t *testing.T) {
	t.Parallel()

	var mockTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		timeHumanizer output.TimeHumanizer
		resource      *cluster.Resource
		nodes         []*node.Resource
		wantOutput    string
		wantErr       bool
	}{
		{
			name: "describe cluster",
			resource: &cluster.Resource{
				ID:                    "bananaCluster",
				DisableTelemetry:      false,
				DisableCrashReporting: true,
				DisableVersionCheck:   false,
				LogLevel:              "debug",
				LogFormat:             "default",
				CreatedAt:             mockTime,
				UpdatedAt:             mockTime,
				Version:               "42",
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
			wantOutput: `ID:               bananaCluster                      
Version:          42                                 
Created at:       2000-01-01T00:00:00Z (xx aeons ago)
Updated at:       2000-01-01T00:00:00Z (xx aeons ago)
Telemetry:        Enabled                            
Crash Reporting:  Disabled                           
Version Check:    Enabled                            
Log Level:        debug                              
Log Format:       default                            
Nodes:                                               
  ID:             bananaNodeID                       
  Name:           bananaNodeName                     
  Health:         online                             
  Address:        127.0.0.1                          
                                                     
  ID:             kiwiNodeID                         
  Name:           kiwiNodeName                       
  Health:         online                             
  Address:        127.0.0.2                          
                                                     
  ID:             pearNodeID                         
  Name:           pearNodeName                       
  Health:         offline                            
  Address:        127.0.0.3                          
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

			err := d.DescribeCluster(context.Background(), w, output.NewCluster(tt.resource, tt.nodes))
			if (err != nil) != tt.wantErr {
				t.Errorf("DescribeCluster() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			gotOutput := w.String()
			if gotOutput != tt.wantOutput {
				compareOutput(t, gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestDisplayer_DescribeLicence(t *testing.T) {
	t.Parallel()

	var mockTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		timeHumanizer output.TimeHumanizer
		resource      *licence.Resource
		wantOutput    string
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
			wantOutput: `ClusterID:      bananaID                           
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

			err := d.DescribeLicence(context.Background(), w, output.NewLicence(tt.resource))
			if (err != nil) != tt.wantErr {
				t.Errorf("DescribeLicence() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotOutput := w.String()
			if gotOutput != tt.wantOutput {
				compareOutput(t, gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestDisplayer_DescribeNamespace(t *testing.T) {
	t.Parallel()

	var mockTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		timeHumanizer output.TimeHumanizer
		resource      *namespace.Resource
		wantOutput    string
		wantErr       bool
	}{
		{
			name: "describe namespace",
			resource: &namespace.Resource{
				ID:   "bananaNamespaceID",
				Name: "bananaNamespaceName",
				Labels: labels.Set{
					"bananaKey": "bananaValue",
					"kiwiKey":   "kiwiValue",
				},
				CreatedAt: mockTime,
				UpdatedAt: mockTime,
				Version:   "42",
			},
			wantOutput: `ID:          bananaNamespaceID                  
Name:        bananaNamespaceName                
Labels:      bananaKey=bananaValue,             
             kiwiKey=kiwiValue                  
Version:     42                                 
Created at:  2000-01-01T00:00:00Z (xx aeons ago)
Updated at:  2000-01-01T00:00:00Z (xx aeons ago)
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

			err := d.DescribeNamespace(context.Background(), w, output.NewNamespace(tt.resource))
			if (err != nil) != tt.wantErr {
				t.Errorf("DescribeNamespace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			gotOutput := w.String()
			if gotOutput != tt.wantOutput {
				compareOutput(t, gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestDisplayer_DescribeListNamespaces(t *testing.T) {
	t.Parallel()

	var mockTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		timeHumanizer output.TimeHumanizer
		resource      []*namespace.Resource
		wantOutput    string
		wantErr       bool
	}{
		{
			name: "describe namespace",
			resource: []*namespace.Resource{
				{
					ID:   "bananaNamespaceID",
					Name: "bananaNamespaceName",
					Labels: labels.Set{
						"bananaKey": "bananaValue",
						"kiwiKey":   "kiwiValue",
					},
					CreatedAt: mockTime,
					UpdatedAt: mockTime,
					Version:   "42",
				},
				{
					ID:        "kiwiNamespaceID",
					Name:      "kiwiNamespaceName",
					Labels:    labels.Set{},
					CreatedAt: mockTime,
					UpdatedAt: mockTime,
					Version:   "43",
				},
				{
					ID:        "pineappleNamespaceID",
					Name:      "pineappleNamespaceName",
					Labels:    labels.Set{},
					CreatedAt: mockTime,
					UpdatedAt: mockTime,
					Version:   "44",
				},
			},
			wantOutput: `ID:          bananaNamespaceID                  
Name:        bananaNamespaceName                
Labels:      bananaKey=bananaValue,             
             kiwiKey=kiwiValue                  
Version:     42                                 
Created at:  2000-01-01T00:00:00Z (xx aeons ago)
Updated at:  2000-01-01T00:00:00Z (xx aeons ago)

ID:          kiwiNamespaceID                    
Name:        kiwiNamespaceName                  
Labels:      -                                  
Version:     43                                 
Created at:  2000-01-01T00:00:00Z (xx aeons ago)
Updated at:  2000-01-01T00:00:00Z (xx aeons ago)

ID:          pineappleNamespaceID               
Name:        pineappleNamespaceName             
Labels:      -                                  
Version:     44                                 
Created at:  2000-01-01T00:00:00Z (xx aeons ago)
Updated at:  2000-01-01T00:00:00Z (xx aeons ago)
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

			err := d.DescribeListNamespaces(context.Background(), w, output.NewNamespaces(tt.resource))
			if (err != nil) != tt.wantErr {
				t.Errorf("DescribeListNamespaces() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			gotOutput := w.String()
			if gotOutput != tt.wantOutput {
				compareOutput(t, gotOutput, tt.wantOutput)
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
				AttachType:     volume.AttachTypeHost,
				NFS: output.NFSConfig{
					Exports: []output.NFSExportConfig{
						{
							ExportID:   1,
							Path:       "/",
							PseudoPath: "/",
							ACLs: []output.NFSExportConfigACL{
								{
									Identity: output.NFSExportConfigACLIdentity{
										IdentityType: "cidr",
										Matcher:      "10.0.0.0/8",
									},
									SquashConfig: output.NFSExportConfigACLSquashConfig{
										GID:    0,
										UID:    0,
										Squash: "root",
									},
									AccessLevel: "rw",
								},
								{
									Identity: output.NFSExportConfigACLIdentity{
										IdentityType: "hostname",
										Matcher:      "*.storageos.com",
									},
									SquashConfig: output.NFSExportConfigACLSquashConfig{
										GID:    1000,
										UID:    1000,
										Squash: "rootuid",
									},
									AccessLevel: "ro",
								},
							},
						},
						{
							ExportID:   2,
							Path:       "/path",
							PseudoPath: "/psuedo",
							ACLs:       []output.NFSExportConfigACL{},
						},
					},
					ServiceEndpoint: "10.0.0.1:/",
				},
				Namespace:     "banana-namespace-id",
				NamespaceName: "banana-namespace",
				Labels: labels.Set{
					"kiwi": "42",
					"pear": "42",
				},
				Filesystem: volume.FsTypeFromString("ext4"),
				SizeBytes:  size.GiB,
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
							BytesRemaining:            256 * size.MiB,
							ThroughputBytes:           0,
							EstimatedSecondsRemaining: 1200,
						},
					},
				},
				CreatedAt: mockTime,
				UpdatedAt: mockTime,
				Version:   "42",
			},

			wantOutput: `ID                      bananaID                                                                        
Name                    banana-name                                                                     
Description             banana description                                                              
AttachedOn              banana-node-a (banana-node-a-id)                                                
Attachment Type         host                                                                            
NFS                                                                                                     
  Service Endpoint      10.0.0.1:/                                                                      
  Exports:                                                                                              
  - ID                  1                                                                               
    Path                /                                                                               
    Pseudo Path         /                                                                               
    ACLs                                                                                                
    - Identity Type     cidr                                                                            
      Identity Matcher  10.0.0.0/8                                                                      
      Squash            root                                                                            
      Squash UID        0                                                                               
      Squash GUID       0                                                                               
    - Identity Type     hostname                                                                        
      Identity Matcher  *.storageos.com                                                                 
      Squash            rootuid                                                                         
      Squash UID        1000                                                                            
      Squash GUID       1000                                                                            
  - ID                  2                                                                               
    Path                /path                                                                           
    Pseudo Path         /psuedo                                                                         
    ACLs                                                                                                
Namespace               banana-namespace (banana-namespace-id)                                          
Labels                  kiwi=42,                                                                        
                        pear=42                                                                         
Filesystem              ext4                                                                            
Size                    1.0 GiB (1073741824 bytes)                                                      
Version                 42                                                                              
Created at              2000-01-01T00:00:00Z (xx aeons ago)                                             
Updated at              2000-01-01T00:00:00Z (xx aeons ago)                                             
                                                                                                        
Master:               
  ID                    bananaDeploymentID1                                                             
  Node                  banana-node1 (banana-node1-id)                                                  
  Health                ready                                                                           
                                                                                                        
Replicas:             
  ID                    bananaDeploymentID2                                                             
  Node                  banana-node2 (banana-node2-id)                                                  
  Health                ready                                                                           
  Promotable            true                                                                            
                                                                                                        
  ID                    bananaDeploymentID3                                                             
  Node                  banana-node3 (banana-node3-id)                                                  
  Health                deleted                                                                         
  Promotable            false                                                                           
                                                                                                        
  ID                    bananaDeploymentID4                                                             
  Node                  banana-node4 (banana-node4-id)                                                  
  Health                syncing                                                                         
  Promotable            false                                                                           
  Sync Progress         768.00 MiB / 1.00 GiB [##########################........] 75.00%  -  ETA: 20m0s
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
				AttachType:     volume.AttachTypeHost,
				NFS: output.NFSConfig{
					Exports: []output.NFSExportConfig{
						{
							ExportID:   1,
							Path:       "/",
							PseudoPath: "/",
							ACLs: []output.NFSExportConfigACL{
								{
									Identity: output.NFSExportConfigACLIdentity{
										IdentityType: "cidr",
										Matcher:      "10.0.0.0/8",
									},
									SquashConfig: output.NFSExportConfigACLSquashConfig{
										GID:    0,
										UID:    0,
										Squash: "root",
									},
									AccessLevel: "rw",
								},
								{
									Identity: output.NFSExportConfigACLIdentity{
										IdentityType: "hostname",
										Matcher:      "*.storageos.com",
									},
									SquashConfig: output.NFSExportConfigACLSquashConfig{
										GID:    1000,
										UID:    1000,
										Squash: "rootuid",
									},
									AccessLevel: "ro",
								},
							},
						},
						{
							ExportID:   2,
							Path:       "/path",
							PseudoPath: "/psuedo",
							ACLs:       []output.NFSExportConfigACL{},
						},
					},
					ServiceEndpoint: "10.0.0.1:/",
				},
				Namespace:     "banana-namespace-id",
				NamespaceName: "banana-namespace",
				Labels: labels.Set{
					"kiwi": "42",
					"pear": "42",
				},
				Filesystem: volume.FsTypeFromString("ext4"),
				SizeBytes:  size.GiB,
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

			wantOutput: `ID                      bananaID                              
Name                    banana-name                           
Description             banana description                    
AttachedOn              banana-node-a (banana-node-a-id)      
Attachment Type         host                                  
NFS                                                           
  Service Endpoint      10.0.0.1:/                            
  Exports:                                                    
  - ID                  1                                     
    Path                /                                     
    Pseudo Path         /                                     
    ACLs                                                      
    - Identity Type     cidr                                  
      Identity Matcher  10.0.0.0/8                            
      Squash            root                                  
      Squash UID        0                                     
      Squash GUID       0                                     
    - Identity Type     hostname                              
      Identity Matcher  *.storageos.com                       
      Squash            rootuid                               
      Squash UID        1000                                  
      Squash GUID       1000                                  
  - ID                  2                                     
    Path                /path                                 
    Pseudo Path         /psuedo                               
    ACLs                                                      
Namespace               banana-namespace (banana-namespace-id)
Labels                  kiwi=42,                              
                        pear=42                               
Filesystem              ext4                                  
Size                    1.0 GiB (1073741824 bytes)            
Version                 42                                    
Created at              2000-01-01T00:00:00Z (xx aeons ago)   
Updated at              2000-01-01T00:00:00Z (xx aeons ago)   
                                                              
Master:               
  ID                    bananaDeploymentID1                   
  Node                  banana-node1 (banana-node1-id)        
  Health                ready                                 
`,
			wantErr: nil,
		},
		{
			name: "describe volume ok with missing sync progress",

			volume: &output.Volume{
				ID:             "bananaID",
				Name:           "banana-name",
				Description:    "banana description",
				AttachedOn:     "banana-node-a-id",
				AttachedOnName: "banana-node-a",
				AttachType:     volume.AttachTypeHost,
				NFS: output.NFSConfig{
					Exports: []output.NFSExportConfig{
						{
							ExportID:   1,
							Path:       "/",
							PseudoPath: "/",
							ACLs: []output.NFSExportConfigACL{
								{
									Identity: output.NFSExportConfigACLIdentity{
										IdentityType: "cidr",
										Matcher:      "10.0.0.0/8",
									},
									SquashConfig: output.NFSExportConfigACLSquashConfig{
										GID:    0,
										UID:    0,
										Squash: "root",
									},
									AccessLevel: "rw",
								},
								{
									Identity: output.NFSExportConfigACLIdentity{
										IdentityType: "hostname",
										Matcher:      "*.storageos.com",
									},
									SquashConfig: output.NFSExportConfigACLSquashConfig{
										GID:    1000,
										UID:    1000,
										Squash: "rootuid",
									},
									AccessLevel: "ro",
								},
							},
						},
						{
							ExportID:   2,
							Path:       "/path",
							PseudoPath: "/psuedo",
							ACLs:       []output.NFSExportConfigACL{},
						},
					},
					ServiceEndpoint: "10.0.0.1:/",
				},
				Namespace:     "banana-namespace-id",
				NamespaceName: "banana-namespace",
				Labels: labels.Set{
					"kiwi": "42",
					"pear": "42",
				},
				Filesystem: volume.FsTypeFromString("ext4"),
				SizeBytes:  size.GiB,
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
						Health:       health.ReplicaSyncing,
						Promotable:   false,
						SyncProgress: nil,
					},
				},
				CreatedAt: mockTime,
				UpdatedAt: mockTime,
				Version:   "42",
			},

			wantOutput: `ID                      bananaID                              
Name                    banana-name                           
Description             banana description                    
AttachedOn              banana-node-a (banana-node-a-id)      
Attachment Type         host                                  
NFS                                                           
  Service Endpoint      10.0.0.1:/                            
  Exports:                                                    
  - ID                  1                                     
    Path                /                                     
    Pseudo Path         /                                     
    ACLs                                                      
    - Identity Type     cidr                                  
      Identity Matcher  10.0.0.0/8                            
      Squash            root                                  
      Squash UID        0                                     
      Squash GUID       0                                     
    - Identity Type     hostname                              
      Identity Matcher  *.storageos.com                       
      Squash            rootuid                               
      Squash UID        1000                                  
      Squash GUID       1000                                  
  - ID                  2                                     
    Path                /path                                 
    Pseudo Path         /psuedo                               
    ACLs                                                      
Namespace               banana-namespace (banana-namespace-id)
Labels                  kiwi=42,                              
                        pear=42                               
Filesystem              ext4                                  
Size                    1.0 GiB (1073741824 bytes)            
Version                 42                                    
Created at              2000-01-01T00:00:00Z (xx aeons ago)   
Updated at              2000-01-01T00:00:00Z (xx aeons ago)   
                                                              
Master:               
  ID                    bananaDeploymentID1                   
  Node                  banana-node1 (banana-node1-id)        
  Health                ready                                 
                                                              
Replicas:             
  ID                    bananaDeploymentID2                   
  Node                  banana-node2 (banana-node2-id)        
  Health                syncing                               
  Promotable            false                                 
  Sync Progress         n/a                                   
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
				compareOutput(t, gotOutput, tt.wantOutput)
			}
		})
	}
}

func compareOutput(t *testing.T, gotOutput string, wantOutput string) {
	t.Helper()

	wantLines := strings.Split(wantOutput, "\n")
	gotLines := strings.Split(gotOutput, "\n")

	if len(wantLines) != len(gotLines) {
		t.Errorf("different number of lines. Got %d, want %d", len(gotLines), len(wantLines))
		return
	}

	for i := range gotLines {
		if gotLines[i] != wantLines[i] {
			t.Errorf("Line %d is different.\nGOT:\n%q\nWANT:\n%q\n", i, gotLines[i], wantLines[i])
			return
		}
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
					AttachType:     volume.AttachTypeHost,
					NFS: output.NFSConfig{
						Exports: []output.NFSExportConfig{
							{
								ExportID:   1,
								Path:       "/",
								PseudoPath: "/",
								ACLs: []output.NFSExportConfigACL{
									{
										Identity: output.NFSExportConfigACLIdentity{
											IdentityType: "cidr",
											Matcher:      "10.0.0.0/8",
										},
										SquashConfig: output.NFSExportConfigACLSquashConfig{
											GID:    0,
											UID:    0,
											Squash: "root",
										},
										AccessLevel: "rw",
									},
									{
										Identity: output.NFSExportConfigACLIdentity{
											IdentityType: "hostname",
											Matcher:      "*.storageos.com",
										},
										SquashConfig: output.NFSExportConfigACLSquashConfig{
											GID:    1000,
											UID:    1000,
											Squash: "rootuid",
										},
										AccessLevel: "ro",
									},
								},
							},
							{
								ExportID:   2,
								Path:       "/path",
								PseudoPath: "/psuedo",
								ACLs:       []output.NFSExportConfigACL{},
							},
						},
						ServiceEndpoint: "10.0.0.1:/",
					},
					Namespace:     "banana-namespace-id",
					NamespaceName: "banana-namespace",
					Labels: labels.Set{
						"kiwi": "42",
						"pear": "42",
					},
					Filesystem: volume.FsTypeFromString("ext4"),
					SizeBytes:  size.GiB,
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
								BytesRemaining:            256 * size.MiB,
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
					AttachType:     volume.AttachTypeHost,
					NFS: output.NFSConfig{
						Exports: []output.NFSExportConfig{
							{
								ExportID:   1,
								Path:       "/",
								PseudoPath: "/",
								ACLs: []output.NFSExportConfigACL{
									{
										Identity: output.NFSExportConfigACLIdentity{
											IdentityType: "cidr",
											Matcher:      "10.0.0.0/8",
										},
										SquashConfig: output.NFSExportConfigACLSquashConfig{
											GID:    0,
											UID:    0,
											Squash: "root",
										},
										AccessLevel: "rw",
									},
									{
										Identity: output.NFSExportConfigACLIdentity{
											IdentityType: "hostname",
											Matcher:      "*.storageos.com",
										},
										SquashConfig: output.NFSExportConfigACLSquashConfig{
											GID:    1000,
											UID:    1000,
											Squash: "rootuid",
										},
										AccessLevel: "ro",
									},
								},
							},
							{
								ExportID:   2,
								Path:       "/path",
								PseudoPath: "/psuedo",
								ACLs:       []output.NFSExportConfigACL{},
							},
						},
						ServiceEndpoint: "10.0.0.1:/",
					},
					Namespace:     "kiwi-namespace-id",
					NamespaceName: "kiwi-namespace",
					Labels: labels.Set{
						"kiwi": "42",
						"pear": "42",
					},
					Filesystem: volume.FsTypeFromString("ext4"),
					SizeBytes:  size.GiB,
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
								BytesRemaining:            256 * size.MiB,
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

			wantOutput: `ID                      bananaID                                                                        
Name                    banana-name                                                                     
Description             banana description                                                              
AttachedOn              banana-node-a (banana-node-a-id)                                                
Attachment Type         host                                                                            
NFS                                                                                                     
  Service Endpoint      10.0.0.1:/                                                                      
  Exports:                                                                                              
  - ID                  1                                                                               
    Path                /                                                                               
    Pseudo Path         /                                                                               
    ACLs                                                                                                
    - Identity Type     cidr                                                                            
      Identity Matcher  10.0.0.0/8                                                                      
      Squash            root                                                                            
      Squash UID        0                                                                               
      Squash GUID       0                                                                               
    - Identity Type     hostname                                                                        
      Identity Matcher  *.storageos.com                                                                 
      Squash            rootuid                                                                         
      Squash UID        1000                                                                            
      Squash GUID       1000                                                                            
  - ID                  2                                                                               
    Path                /path                                                                           
    Pseudo Path         /psuedo                                                                         
    ACLs                                                                                                
Namespace               banana-namespace (banana-namespace-id)                                          
Labels                  kiwi=42,                                                                        
                        pear=42                                                                         
Filesystem              ext4                                                                            
Size                    1.0 GiB (1073741824 bytes)                                                      
Version                 42                                                                              
Created at              2000-01-01T00:00:00Z (xx aeons ago)                                             
Updated at              2000-01-01T00:00:00Z (xx aeons ago)                                             
                                                                                                        
Master:               
  ID                    bananaDeploymentID1                                                             
  Node                  banana-node1 (banana-node1-id)                                                  
  Health                ready                                                                           
                                                                                                        
Replicas:             
  ID                    bananaDeploymentID2                                                             
  Node                  banana-node2 (banana-node2-id)                                                  
  Health                ready                                                                           
  Promotable            true                                                                            
                                                                                                        
  ID                    bananaDeploymentID3                                                             
  Node                  banana-node3 (banana-node3-id)                                                  
  Health                deleted                                                                         
  Promotable            false                                                                           
                                                                                                        
  ID                    bananaDeploymentID4                                                             
  Node                  banana-node4 (banana-node4-id)                                                  
  Health                syncing                                                                         
  Promotable            false                                                                           
  Sync Progress         768.00 MiB / 1.00 GiB [##########################........] 75.00%  -  ETA: 20m0s

ID                      kiwiID                                                                          
Name                    kiwi-name                                                                       
Description             kiwi description                                                                
AttachedOn              kiwi-node-a (kiwi-node-a-id)                                                    
Attachment Type         host                                                                            
NFS                                                                                                     
  Service Endpoint      10.0.0.1:/                                                                      
  Exports:                                                                                              
  - ID                  1                                                                               
    Path                /                                                                               
    Pseudo Path         /                                                                               
    ACLs                                                                                                
    - Identity Type     cidr                                                                            
      Identity Matcher  10.0.0.0/8                                                                      
      Squash            root                                                                            
      Squash UID        0                                                                               
      Squash GUID       0                                                                               
    - Identity Type     hostname                                                                        
      Identity Matcher  *.storageos.com                                                                 
      Squash            rootuid                                                                         
      Squash UID        1000                                                                            
      Squash GUID       1000                                                                            
  - ID                  2                                                                               
    Path                /path                                                                           
    Pseudo Path         /psuedo                                                                         
    ACLs                                                                                                
Namespace               kiwi-namespace (kiwi-namespace-id)                                              
Labels                  kiwi=42,                                                                        
                        pear=42                                                                         
Filesystem              ext4                                                                            
Size                    1.0 GiB (1073741824 bytes)                                                      
Version                 42                                                                              
Created at              2000-01-01T00:00:00Z (xx aeons ago)                                             
Updated at              2000-01-01T00:00:00Z (xx aeons ago)                                             
                                                                                                        
Master:               
  ID                    kiwiDeploymentID1                                                               
  Node                  kiwi-node1 (kiwi-node1-id)                                                      
  Health                ready                                                                           
                                                                                                        
Replicas:             
  ID                    kiwiDeploymentID2                                                               
  Node                  kiwi-node2 (kiwi-node2-id)                                                      
  Health                ready                                                                           
  Promotable            true                                                                            
                                                                                                        
  ID                    kiwiDeploymentID3                                                               
  Node                  kiwi-node3 (kiwi-node3-id)                                                      
  Health                deleted                                                                         
  Promotable            false                                                                           
                                                                                                        
  ID                    kiwiDeploymentID4                                                               
  Node                  kiwi-node4 (kiwi-node4-id)                                                      
  Health                syncing                                                                         
  Promotable            false                                                                           
  Sync Progress         768.00 MiB / 1.00 GiB [##########################........] 75.00%  -  ETA: 20m0s
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
				compareOutput(t, gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestDisplayer_DescribePolicyGroup(t *testing.T) {
	t.Parallel()

	var mockTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name string

		group      *policygroup.Resource
		namespaces []*namespace.Resource

		wantOutput string
		wantErr    error
	}{
		{
			name: "describe policy group ok",
			group: &policygroup.Resource{
				ID:   "bananaID",
				Name: "banana-name",
				Users: []*policygroup.Member{
					{
						ID:       "banana-user-1",
						Username: "banana-username-1",
					},
					{
						ID:       "banana-user-2",
						Username: "banana-username-2",
					},
					{
						ID:       "banana-user-3",
						Username: "banana-username-3",
					},
					{
						ID:       "banana-user-4",
						Username: "banana-username-4",
					},
				},
				Specs: []*policygroup.Spec{
					{
						NamespaceID:  "banana-namespace-id",
						ResourceType: "volume",
						ReadOnly:     false,
					},
					{
						NamespaceID:  "banana-namespace-id",
						ResourceType: "*",
						ReadOnly:     true,
					},
				},
				CreatedAt: mockTime,
				UpdatedAt: mockTime,
				Version:   "42",
			},
			namespaces: []*namespace.Resource{
				{
					ID:        "banana-namespace-id",
					Name:      "banana-namespace-name",
					Labels:    labels.Set{},
					CreatedAt: mockTime,
					UpdatedAt: mockTime,
					Version:   "42",
				},
			},
			wantOutput: `ID          bananaID                             
Name        banana-name                          
Specs:                                           
            write volume on banana-namespace-name
             read      * on banana-namespace-name
Members:                                         
            banana-username-1                    
            banana-username-2                    
            banana-username-3                    
            banana-username-4                    
Created at  2000-01-01T00:00:00Z (xx aeons ago)  
Updated at  2000-01-01T00:00:00Z (xx aeons ago)  
Version     42                                   
`,
			wantErr: nil,
		},
		{
			name: "describe policy group with no users ok",

			group: &policygroup.Resource{
				ID:    "bananaID",
				Name:  "banana-name",
				Users: []*policygroup.Member{},
				Specs: []*policygroup.Spec{
					{
						NamespaceID:  "banana-namespace-id",
						ResourceType: "volume",
						ReadOnly:     false,
					},
					{
						NamespaceID:  "banana-namespace-id",
						ResourceType: "*",
						ReadOnly:     true,
					},
				},
				CreatedAt: mockTime,
				UpdatedAt: mockTime,
				Version:   "42",
			},
			namespaces: []*namespace.Resource{
				{
					ID:        "banana-namespace-id",
					Name:      "banana-namespace-name",
					Labels:    labels.Set{},
					CreatedAt: mockTime,
					UpdatedAt: mockTime,
					Version:   "42",
				},
			},

			wantOutput: `ID          bananaID                             
Name        banana-name                          
Specs:                                           
            write volume on banana-namespace-name
             read      * on banana-namespace-name
Members:    []                                   
Created at  2000-01-01T00:00:00Z (xx aeons ago)  
Updated at  2000-01-01T00:00:00Z (xx aeons ago)  
Version     42                                   
`,
			wantErr: nil,
		},

		{
			name: "describe policy group with no specs ok",

			group: &policygroup.Resource{
				ID:   "bananaID",
				Name: "banana-name",
				Users: []*policygroup.Member{
					{
						ID:       "banana-user-1",
						Username: "banana-username-1",
					},
					{
						ID:       "banana-user-2",
						Username: "banana-username-2",
					},
					{
						ID:       "banana-user-3",
						Username: "banana-username-3",
					},
					{
						ID:       "banana-user-4",
						Username: "banana-username-4",
					},
				},
				Specs:     []*policygroup.Spec{},
				CreatedAt: mockTime,
				UpdatedAt: mockTime,
				Version:   "42",
			},
			namespaces: []*namespace.Resource{},

			wantOutput: `ID          bananaID                           
Name        banana-name                        
Specs:      []                                 
Members:                                       
            banana-username-1                  
            banana-username-2                  
            banana-username-3                  
            banana-username-4                  
Created at  2000-01-01T00:00:00Z (xx aeons ago)
Updated at  2000-01-01T00:00:00Z (xx aeons ago)
Version     42                                 
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

			gotErr := d.DescribePolicyGroup(context.Background(), w, output.NewPolicyGroup(tt.group, tt.namespaces))
			if gotErr != tt.wantErr {
				t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
				return
			}

			gotOutput := w.String()
			if gotOutput != tt.wantOutput {
				compareOutput(t, gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestDisplayer_DescribeListPolicyGroups(t *testing.T) {
	t.Parallel()

	var mockTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name string

		groups     []*policygroup.Resource
		namespaces []*namespace.Resource

		wantOutput string
		wantErr    error
	}{
		{
			name: "describe policy group ok",

			groups: []*policygroup.Resource{
				{
					ID:   "bananaID",
					Name: "banana-name",
					Users: []*policygroup.Member{
						{
							ID:       "banana-user-1",
							Username: "banana-username-1",
						},
						{
							ID:       "banana-user-2",
							Username: "banana-username-2",
						},
						{
							ID:       "banana-user-3",
							Username: "banana-username-3",
						},
						{
							ID:       "banana-user-4",
							Username: "banana-username-4",
						},
					},
					Specs: []*policygroup.Spec{
						{
							NamespaceID:  "banana-namespace-id",
							ResourceType: "volume",
							ReadOnly:     false,
						},
						{
							NamespaceID:  "banana-namespace-id",
							ResourceType: "*",
							ReadOnly:     true,
						},
					},
					CreatedAt: mockTime,
					UpdatedAt: mockTime,
					Version:   "42",
				},
				{
					ID:   "pineappleID",
					Name: "pineapple-name",
					Users: []*policygroup.Member{
						{
							ID:       "pineapple-user-1",
							Username: "pineapple-username-1",
						},
						{
							ID:       "pineapple-user-2",
							Username: "pineapple-username-2",
						},
						{
							ID:       "pineapple-user-3",
							Username: "pineapple-username-3",
						},
						{
							ID:       "pineapple-user-4",
							Username: "pineapple-username-4",
						},
					},
					Specs: []*policygroup.Spec{
						{
							NamespaceID:  "pineapple-namespace-id",
							ResourceType: "volume",
							ReadOnly:     false,
						},
						{
							NamespaceID:  "pineapple-namespace-id",
							ResourceType: "*",
							ReadOnly:     true,
						},
					},
					CreatedAt: mockTime,
					UpdatedAt: mockTime,
					Version:   "43",
				},
				{
					ID:   "kiwiID",
					Name: "kiwi-name",
					Users: []*policygroup.Member{
						{
							ID:       "kiwi-user-1",
							Username: "kiwi-username-1",
						},
						{
							ID:       "kiwi-user-2",
							Username: "kiwi-username-2",
						},
						{
							ID:       "kiwi-user-3",
							Username: "kiwi-username-3",
						},
						{
							ID:       "kiwi-user-4",
							Username: "kiwi-username-4",
						},
					},
					Specs: []*policygroup.Spec{
						{
							NamespaceID:  "kiwi-namespace-id",
							ResourceType: "volume",
							ReadOnly:     false,
						},
						{
							NamespaceID:  "kiwi-namespace-id",
							ResourceType: "*",
							ReadOnly:     true,
						},
					},
					CreatedAt: mockTime,
					UpdatedAt: mockTime,
					Version:   "44",
				},
			},

			namespaces: []*namespace.Resource{
				{
					ID:        "banana-namespace-id",
					Name:      "banana-namespace-name",
					Labels:    labels.Set{},
					CreatedAt: mockTime,
					UpdatedAt: mockTime,
					Version:   "42",
				},
				{
					ID:        "kiwi-namespace-id",
					Name:      "kiwi-namespace-name",
					Labels:    labels.Set{},
					CreatedAt: mockTime,
					UpdatedAt: mockTime,
					Version:   "42",
				},
			},

			wantOutput: `ID          bananaID                             
Name        banana-name                          
Specs:                                           
            write volume on banana-namespace-name
             read      * on banana-namespace-name
Members:                                         
            banana-username-1                    
            banana-username-2                    
            banana-username-3                    
            banana-username-4                    
Created at  2000-01-01T00:00:00Z (xx aeons ago)  
Updated at  2000-01-01T00:00:00Z (xx aeons ago)  
Version     42                                   

ID          pineappleID                           
Name        pineapple-name                        
Specs:                                            
            write volume on pineapple-namespace-id
             read      * on pineapple-namespace-id
Members:                                          
            pineapple-username-1                  
            pineapple-username-2                  
            pineapple-username-3                  
            pineapple-username-4                  
Created at  2000-01-01T00:00:00Z (xx aeons ago)   
Updated at  2000-01-01T00:00:00Z (xx aeons ago)   
Version     43                                    

ID          kiwiID                             
Name        kiwi-name                          
Specs:                                         
            write volume on kiwi-namespace-name
             read      * on kiwi-namespace-name
Members:                                       
            kiwi-username-1                    
            kiwi-username-2                    
            kiwi-username-3                    
            kiwi-username-4                    
Created at  2000-01-01T00:00:00Z (xx aeons ago)
Updated at  2000-01-01T00:00:00Z (xx aeons ago)
Version     44                                 
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

			gotErr := d.DescribeListPolicyGroups(context.Background(), w, output.NewPolicyGroups(tt.groups, tt.namespaces))
			if gotErr != tt.wantErr {
				t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
				return
			}

			gotOutput := w.String()
			if gotOutput != tt.wantOutput {
				pretty.Ldiff(t, gotOutput, tt.wantOutput)
				// gotOutput = strings.ReplaceAll(w.String(), " ", "•")
				// gotOutput = strings.ReplaceAll(gotOutput, "\n", "%\n")
				t.Errorf("got output: \n%v\n\nwant: \n%v\n", gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestDisplayer_DescribeUser(t *testing.T) {
	t.Parallel()

	var mockTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name string

		usr          *user.Resource
		policyGroups []*policygroup.Resource

		wantOutput string
		wantErr    error
	}{
		{
			name: "describe user ok",
			usr: &user.Resource{
				ID:       "bananaID",
				Username: "bananaUsername",
				IsAdmin:  false,
				Groups: []id.PolicyGroup{
					"banana-policy-id",
					"pineapple-policy-id",
				},
				CreatedAt: mockTime,
				UpdatedAt: mockTime,
				Version:   "42",
			},
			policyGroups: []*policygroup.Resource{
				{
					ID:   "banana-policy-id",
					Name: "banana-policy-name",
				},
				{
					ID:   "kiwi-policy-id",
					Name: "kiwi-policy-name",
				},
				{
					ID:   "pineapple-policy-id",
					Name: "pineapple-policy-name",
				},
			},
			wantOutput: `ID          bananaID                           
Username    bananaUsername                     
Admin       false                              
Version     42                                 
Created at  2000-01-01T00:00:00Z (xx aeons ago)
Updated at  2000-01-01T00:00:00Z (xx aeons ago)
Policies:                                      
         -  banana-policy-name                 
         -  pineapple-policy-name              
`,
			wantErr: nil,
		},
		{
			name: "describe user with no policy ok",

			usr: &user.Resource{
				ID:        "bananaID",
				Username:  "bananaUsername",
				IsAdmin:   false,
				Groups:    []id.PolicyGroup{},
				CreatedAt: mockTime,
				UpdatedAt: mockTime,
				Version:   "42",
			},
			policyGroups: []*policygroup.Resource{
				{
					ID:   "banana-policy-id",
					Name: "banana-policy-name",
				},
				{
					ID:   "kiwi-policy-id",
					Name: "kiwi-policy-name",
				},
				{
					ID:   "pineapple-policy-id",
					Name: "pineapple-policy-name",
				},
			},
			wantOutput: `ID          bananaID                           
Username    bananaUsername                     
Admin       false                              
Version     42                                 
Created at  2000-01-01T00:00:00Z (xx aeons ago)
Updated at  2000-01-01T00:00:00Z (xx aeons ago)
Policies:   []                                 
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

			gotErr := d.DescribeUser(context.Background(), w, output.NewUser(tt.usr, tt.policyGroups))
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

func TestDisplayer_DescribeListUsers(t *testing.T) {
	t.Parallel()

	var mockTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name string

		users        []*user.Resource
		policyGroups []*policygroup.Resource

		wantOutput string
		wantErr    error
	}{
		{
			name: "describe list user ok",
			users: []*user.Resource{
				{
					ID:       "bananaID",
					Username: "bananaUsername",
					IsAdmin:  false,
					Groups: []id.PolicyGroup{
						"banana-policy-id",
						"pineapple-policy-id",
					},
					CreatedAt: mockTime,
					UpdatedAt: mockTime,
					Version:   "42",
				},
				{
					ID:        "kiwiID",
					Username:  "kiwiUsername",
					IsAdmin:   true,
					Groups:    nil,
					CreatedAt: mockTime,
					UpdatedAt: mockTime,
					Version:   "42",
				},
			},
			policyGroups: []*policygroup.Resource{
				{
					ID:   "banana-policy-id",
					Name: "banana-policy-name",
				},
				{
					ID:   "kiwi-policy-id",
					Name: "kiwi-policy-name",
				},
				{
					ID:   "pineapple-policy-id",
					Name: "pineapple-policy-name",
				},
			},
			wantOutput: `ID          bananaID                           
Username    bananaUsername                     
Admin       false                              
Version     42                                 
Created at  2000-01-01T00:00:00Z (xx aeons ago)
Updated at  2000-01-01T00:00:00Z (xx aeons ago)
Policies:                                      
         -  banana-policy-name                 
         -  pineapple-policy-name              

ID          kiwiID                             
Username    kiwiUsername                       
Admin       true                               
Version     42                                 
Created at  2000-01-01T00:00:00Z (xx aeons ago)
Updated at  2000-01-01T00:00:00Z (xx aeons ago)
Policies:   []                                 
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

			gotErr := d.DescribeListUsers(context.Background(), w, output.NewUsers(tt.users, tt.policyGroups))
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
