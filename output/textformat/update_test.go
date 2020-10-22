package textformat

import (
	"bytes"
	"context"
	"testing"
	"time"

	"code.storageos.net/storageos/c2-cli/licence"
	"code.storageos.net/storageos/c2-cli/output"
	"code.storageos.net/storageos/c2-cli/pkg/health"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/labels"
	"code.storageos.net/storageos/c2-cli/pkg/size"
	"code.storageos.net/storageos/c2-cli/pkg/version"
	"code.storageos.net/storageos/c2-cli/volume"
)

func TestDisplayer_UpdateLicence(t *testing.T) {
	t.Parallel()

	var mockTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		timeHumanizer output.TimeHumanizer
		licence       *licence.Resource
		wantW         string
		wantErr       bool
	}{
		{
			name: "print licence update success",
			licence: &licence.Resource{
				ClusterID:            "bananaCluster",
				ExpiresAt:            mockTime,
				ClusterCapacityBytes: 42 * size.GiB,
				UsedBytes:            42 / 2 * size.GiB,
				Kind:                 "bananaLicence",
				CustomerName:         "bananaCustomer",
			},
			wantW: `Licence applied to cluster bananaCluster.

Expiration:     2000-01-01T00:00:00Z (xx aeons ago)
Capacity:       42 GiB                             
Used:           21 GiB                             
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
				t.Errorf("UpdateLicence() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("UpdateLicence() gotW = \n%v\n, want \n%v\n", gotW, tt.wantW)
			}
		})
	}
}

func TestDisplayer_UpdateVolumeDescription(t *testing.T) {
	t.Parallel()

	var mockTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		volume  *volume.Resource
		wantW   string
		wantErr bool
	}{
		{
			name: "print volume update success",
			volume: &volume.Resource{
				ID:          "banana-id",
				Name:        "banana-name",
				Description: "ceciestunvolume",
				AttachedOn:  "barnacle",
				Namespace:   "spaaaaace",
				Labels: labels.Set{
					"mymarvelouslabels": "sogood",
					"woohoo":            "woohooo",
				},
				Filesystem: "ext4",
				SizeBytes:  42,
				Master:     nil,
				Replicas: []*volume.Deployment{
					{
						ID:           "id1",
						Node:         "node1",
						Health:       health.ReplicaReady,
						Promotable:   true,
						SyncProgress: nil,
					},
					{
						ID:           "id2",
						Node:         "node2",
						Health:       health.ReplicaFailed,
						Promotable:   true,
						SyncProgress: nil,
					},
				},
				CreatedAt: mockTime,
				UpdatedAt: mockTime,
				Version:   version.Version("MQ"),
			},
			wantW: `Name:                  banana-name        
ID:                    banana-id          
Size:                  42 B               
Description:           ceciestunvolume    
AttachedOn:            barnacle           
Replicas:              1x ready, 1x failed
Labels:                                   
  - mymarvelouslabels  sogood             
  - woohoo             woohooo            

Volume banana-name (banana-id) updated. Description changed.
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

			err := d.UpdateVolumeDescription(context.Background(), w, output.NewVolumeUpdate(tt.volume))
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateVolume() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("UpdateVolume() gotW = \n%v\n, want \n%v\n", gotW, tt.wantW)
			}
		})
	}
}

func TestDisplayer_UpdateNFSVolumeExports(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		volumeID id.Volume
		exports  []volume.NFSExportConfig
		wantW    string
		wantErr  bool
	}{
		{
			name:     "print licence update success",
			volumeID: "banana",
			exports: []volume.NFSExportConfig{
				{
					ExportID:   42,
					Path:       "/banana",
					PseudoPath: "/kiwi",
					ACLs: []volume.NFSExportConfigACL{
						{
							Identity: volume.NFSExportConfigACLIdentity{
								IdentityType: "bananaIdentity",
								Matcher:      "bananaMatcher",
							},
							SquashConfig: volume.NFSExportConfigACLSquashConfig{
								GID:    42,
								UID:    43,
								Squash: "bananaSquash",
							},
							AccessLevel: "bananaAccessLevel",
						},
						{
							Identity: volume.NFSExportConfigACLIdentity{
								IdentityType: "kiwiIdentity",
								Matcher:      "kiwiMatcher",
							},
							SquashConfig: volume.NFSExportConfigACLSquashConfig{
								GID:    42,
								UID:    43,
								Squash: "kiwiSquash",
							},
							AccessLevel: "kiwiAccessLevel",
						},
					},
				},
				{
					ExportID:   43,
					Path:       "/pineapple",
					PseudoPath: "/orange",
					ACLs: []volume.NFSExportConfigACL{
						{
							Identity: volume.NFSExportConfigACLIdentity{
								IdentityType: "pineappleIdentity",
								Matcher:      "pineappleMatcher",
							},
							SquashConfig: volume.NFSExportConfigACLSquashConfig{
								GID:    42,
								UID:    43,
								Squash: "pineappleSquash",
							},
							AccessLevel: "pineappleAccessLevel",
						},
						{
							Identity: volume.NFSExportConfigACLIdentity{
								IdentityType: "orangeIdentity",
								Matcher:      "orangeMatcher",
							},
							SquashConfig: volume.NFSExportConfigACLSquashConfig{
								GID:    42,
								UID:    43,
								Squash: "orangeSquash",
							},
							AccessLevel: "orangeAccessLevel",
						},
					},
				},
			},
			wantW: `Volume banana updated. NFS export configs changed with: 
                                      
---                                   
ID:               42                  
Path              /banana             
PseudoPath:       /kiwi               
ACLs:                                 
- Access Level:   bananaAccessLevel   
  Identity:                           
    Type:         bananaIdentity      
    Matcher:      bananaMatcher       
  Squash Config:                      
    GID:          42                  
    UID:          43                  
    Squash:       bananaSquash        
- Access Level:   kiwiAccessLevel     
  Identity:                           
    Type:         kiwiIdentity        
    Matcher:      kiwiMatcher         
  Squash Config:                      
    GID:          42                  
    UID:          43                  
    Squash:       kiwiSquash          
                                      
---                                   
ID:               43                  
Path              /pineapple          
PseudoPath:       /orange             
ACLs:                                 
- Access Level:   pineappleAccessLevel
  Identity:                           
    Type:         pineappleIdentity   
    Matcher:      pineappleMatcher    
  Squash Config:                      
    GID:          42                  
    UID:          43                  
    Squash:       pineappleSquash     
- Access Level:   orangeAccessLevel   
  Identity:                           
    Type:         orangeIdentity      
    Matcher:      orangeMatcher       
  Squash Config:                      
    GID:          42                  
    UID:          43                  
    Squash:       orangeSquash        
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

			err := d.UpdateNFSVolumeExports(context.Background(), w, tt.volumeID, output.NewNFSExportConfigs(tt.exports))
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateNFSVolumeExports() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("UpdateNFSVolumeExports() gotW = \n%+q\n, want \n%+q\n", gotW, tt.wantW)
			}
		})
	}
}

func TestDisplayer_UpdateNFSVolumeMountEndpoint(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		volumeID id.Volume
		endpoint string
		wantW    string
		wantErr  bool
	}{
		{
			name:     "print NFS volume mount endpoint update success",
			volumeID: "banana",
			endpoint: "10.0.0.1:/",
			wantW: `Volume banana updated. NFS mount endpoint changed with: 10.0.0.1:/
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

			err := d.UpdateNFSVolumeMountEndpoint(context.Background(), w, tt.volumeID, tt.endpoint)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateNFSVolumeMountEndpoint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("UpdateNFSVolumeMountEndpoint() gotW = \n%v\n, want \n%v\n", gotW, tt.wantW)
			}
		})
	}
}
