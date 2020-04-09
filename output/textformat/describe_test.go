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
	"code.storageos.net/storageos/c2-cli/policygroup"
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

func TestDisplayer_DescribeCluster(t *testing.T) {
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
				ID: "bananaCluster",
				Licence: &cluster.Licence{
					ClusterID:            "bananaCluster",
					ExpiresAt:            mockTime,
					ClusterCapacityBytes: 42 * humanize.GiByte,
					Kind:                 "bananaLicence",
					CustomerName:         "bananaCustomer",
				},
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
			wantW: `ID:               bananaCluster                      
Licence:                                             
  expiration:     2000-01-01T00:00:00Z (xx aeons ago)
  capacity:       42 GiB (45097156608)               
  kind:           bananaLicence                      
  customer name:  bananaCustomer                     
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
				t.Errorf("GetCluster() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("GetCluster() gotW = \n%v\n, want \n%v\n", gotW, tt.wantW)
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
				pretty.Ldiff(t, gotOutput, tt.wantOutput)
				// gotOutput = strings.ReplaceAll(w.String(), " ", "•")
				// gotOutput = strings.ReplaceAll(gotOutput, "\n", "%\n")
				t.Errorf("got output: \n%v\n\nwant: \n%v\n", gotOutput, tt.wantOutput)
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
