package formatter

import (
	"bytes"
	"testing"

	"github.com/storageos/go-api/types"
	cliconfig "github.com/storageos/go-cli/cli/config"
)

func TestVolumeWrite(t *testing.T) {
	cases := []struct {
		context  Context
		expected string
	}{
		// Quiet format
		{
			Context{Format: NewVolumeFormat(defaultVolumeQuietFormat, true)},
			`default/myVol
production/prodVol
chaos/unknown
`,
		},
		// Table format
		{
			Context{Format: NewVolumeFormat(defaultVolumeTableFormat, false)},
			`NAMESPACE/NAME      SIZE    MOUNT  SELECTOR  STATUS       REPLICAS  LOCATION
default/myVol       100GiB                   active       0/0       storageos-1 (healthy)
production/prodVol  50GiB                    active       1/1       storageos-2 (healthy)
chaos/unknown       5GiB                     unavailable  0/1       storageos-1 (unknown)
`,
		},
	}

	volumes := []*types.Volume{
		{
			Health: "healthy",
			Master: &types.Deployment{
				Node:     "1",
				NodeName: "storageos-1",
			},
			Name:      "myVol",
			Namespace: "default",
			Size:      100,
			Status:    "active",
		},
		{
			Health: "healthy",
			Labels: map[string]string{
				cliconfig.FeatureReplicas: "1",
			},
			Master: &types.Deployment{
				Node:     "2",
				NodeName: "storageos-2",
			},
			Name:      "prodVol",
			Namespace: "production",
			Replicas: []*types.Deployment{
				&types.Deployment{
					Health:   "healthy",
					Node:     "1",
					NodeName: "storageos-1",
					Status:   "active",
				},
			},
			Size:   50,
			Status: "active",
		},
		{
			Labels: map[string]string{
				cliconfig.FeatureReplicas: "1",
			},
			Master: &types.Deployment{
				Node:     "1",
				NodeName: "storageos-1",
			},
			Name:      "unknown",
			Namespace: "chaos",
			Size:      5,
			Status:    "unavailable",
		},
	}

	nodes := []*types.Node{
		{
			NodeConfig: types.NodeConfig{
				ID: "1",
			},
			Name: "storageos-1",
		},
		{
			NodeConfig: types.NodeConfig{
				ID: "2",
			},
			Name: "storageos-2",
		},
	}

	for _, test := range cases {
		output := bytes.NewBufferString("")
		test.context.Output = output

		if err := VolumeWrite(test.context, volumes, nodes); err != nil {
			t.Fatalf("unexpected error while writing volume context: %s", err.Error())
		} else {
			if test.expected != output.String() {
				t.Errorf("unexpected result.\nexpected:\n%s\ngot:\n%s\n", test.expected, output)
			}
		}
	}
}
