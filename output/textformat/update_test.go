package textformat

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/dustin/go-humanize"

	"code.storageos.net/storageos/c2-cli/cluster"
	"code.storageos.net/storageos/c2-cli/output"
)

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
