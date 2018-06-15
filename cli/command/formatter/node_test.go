package formatter

import (
	"bytes"
	"testing"
	"time"

	"github.com/storageos/go-api/types"
)

func TestNodeWrite(t *testing.T) {
	cases := []struct {
		context  Context
		expected string
	}{
		// Quiet format
		{
			Context{Format: NewNodeFormat(defaultNodeQuietFormat, true)},
			`storageos-1
storageos-2
storageos-3
`,
		},
		// Table format test cases
		{
			// Test default table format
			Context{Format: NewNodeFormat(defaultNodeTableFormat, false)},
			`NAME         ADDRESS    HEALTH                      SCHEDULER  VOLUMES     TOTAL  USED    VERSION
storageos-1  127.0.0.1  Alive Less than a second    true       M: 0, R: 2  5GB    80.00%  1.0.0
storageos-2  127.0.0.1  Alive Less than a second    false      M: 1, R: 0  5GB    52.00%  1.0.0
storageos-3  127.0.0.1  Unknown Less than a second  false      M: 1, R: 2  5GB    60.00%  1.0.0
`,
		},
	}

	aliveStatus := types.SubModuleStatus{Status: "alive"}
	unknownStatus := types.SubModuleStatus{Status: "unknown"}

	nodes := []*types.Node{
		{Name: "storageos-1",
			CapacityStats: types.CapacityStats{
				AvailableCapacityBytes: 1e9,
				TotalCapacityBytes:     5e9,
			},
			Health:          aliveStatus.Status,
			HealthUpdatedAt: time.Now(),
			NodeConfig: types.NodeConfig{
				Address: "127.0.0.1",
			},
			Scheduler: true,
			VolumeStats: types.VolumeStats{
				MasterVolumeCount:  0,
				ReplicaVolumeCount: 2,
			},
			VersionInfo: map[string]types.VersionInfo{
				"storageos": types.VersionInfo{
					Version: "1.0.0",
				},
			},
		},
		{Name: "storageos-2",
			CapacityStats: types.CapacityStats{
				AvailableCapacityBytes: 2.4e9,
				TotalCapacityBytes:     5e9,
			},
			Health:          aliveStatus.Status,
			HealthUpdatedAt: time.Now(),
			NodeConfig: types.NodeConfig{
				Address: "127.0.0.1",
			},
			Scheduler: false,
			VolumeStats: types.VolumeStats{
				MasterVolumeCount:  1,
				ReplicaVolumeCount: 0,
			},
			VersionInfo: map[string]types.VersionInfo{
				"storageos": types.VersionInfo{
					Version: "1.0.0",
				},
			},
		},
		{Name: "storageos-3",
			CapacityStats: types.CapacityStats{
				AvailableCapacityBytes: 2e9,
				TotalCapacityBytes:     5e9,
			},
			Health:          unknownStatus.Status,
			HealthUpdatedAt: time.Now(),
			NodeConfig: types.NodeConfig{
				Address: "127.0.0.1",
			},
			Scheduler: false,
			VolumeStats: types.VolumeStats{
				MasterVolumeCount:  1,
				ReplicaVolumeCount: 2,
			},
			VersionInfo: map[string]types.VersionInfo{
				"storageos": types.VersionInfo{
					Version: "1.0.0",
				},
			},
		},
	}

	for _, test := range cases {
		output := bytes.NewBufferString("")
		test.context.Output = output

		if err := NodeWrite(test.context, nodes); err != nil {
			t.Fatalf("unexpected error while writing node context: %s", err.Error())
		} else {
			if test.expected != output.String() {
				t.Errorf("unexpected result.\nexpected:\n%s\ngot:\n%s\n", test.expected, output)
			}
		}

	}
}
