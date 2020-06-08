package textformat

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/dustin/go-humanize"

	"code.storageos.net/storageos/c2-cli/pkg/health"

	"code.storageos.net/storageos/c2-cli/licence"
	"code.storageos.net/storageos/c2-cli/pkg/labels"
	"code.storageos.net/storageos/c2-cli/pkg/version"
	"code.storageos.net/storageos/c2-cli/volume"

	"code.storageos.net/storageos/c2-cli/output"
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
				ClusterCapacityBytes: 42 * humanize.GiByte,
				UsedBytes:            42 / 2 * humanize.GiByte,
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

func TestDisplayer_UpdateVolume(t *testing.T) {
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

Volume banana-name (banana-id) updated.
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

			err := d.UpdateVolume(context.Background(), w, output.NewVolumeUpdate(tt.volume))
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
