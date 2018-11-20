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
storageos-1  127.0.0.1  Alive Less than a second    true       M: 0, R: 2  5GiB   80.00%  1.0.0
storageos-2  127.0.0.1  Alive Less than a second    false      M: 1, R: 0  5GiB   52.00%  1.0.0
storageos-3  127.0.0.1  Unknown Less than a second  false      M: 1, R: 2  5GiB   60.00%  1.0.0
`,
		},
		{
			// Custom IAAS table format
			Context{Format: NewNodeFormat("table {{.Name}}\t{{.Address}}\t{{.Health}}\t{{.CapacityUsed}}\t{{.Region}}\t{{.FailureDomain}}", false)},
			`NAME         ADDRESS    HEALTH                      USED    REGION  FAILURE_DOMAIN
storageos-1  127.0.0.1  Alive Less than a second    80.00%  euw     euw-1
storageos-2  127.0.0.1  Alive Less than a second    52.00%  euw     euw-2
storageos-3  127.0.0.1  Unknown Less than a second  60.00%  euw     euw-3
`,
		},
	}

	aliveStatus := types.SubModuleStatus{Status: "alive"}
	unknownStatus := types.SubModuleStatus{Status: "unknown"}

	nodes := []*types.Node{
		{Name: "storageos-1",
			CapacityStats: types.CapacityStats{
				AvailableCapacityBytes: 1 * GiB,
				TotalCapacityBytes:     5 * GiB,
			},
			Health:          aliveStatus.Status,
			HealthUpdatedAt: time.Now(),
			NodeConfig: types.NodeConfig{
				Address: "127.0.0.1",
				Labels: map[string]string{
					"iaas/failure-domain": "euw-1",
					"iaas/region":         "euw",
				},
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
				AvailableCapacityBytes: GiB * 12 / 5, // 2.4
				TotalCapacityBytes:     5 * GiB,
			},
			Health:          aliveStatus.Status,
			HealthUpdatedAt: time.Now(),
			NodeConfig: types.NodeConfig{
				Address: "127.0.0.1",
				Labels: map[string]string{
					"iaas/failure-domain": "euw-2",
					"iaas/region":         "euw",
				},
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
				AvailableCapacityBytes: 2 * GiB,
				TotalCapacityBytes:     5 * GiB,
			},
			Health:          unknownStatus.Status,
			HealthUpdatedAt: time.Now(),
			NodeConfig: types.NodeConfig{
				Address: "127.0.0.1",
				Labels: map[string]string{
					"iaas/failure-domain": "euw-3",
					"iaas/region":         "euw",
				},
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

func TestRegion(t *testing.T) {
	cases := []struct {
		context  nodeContext
		expected string
	}{
		{
			nodeContext{
				HeaderContext: HeaderContext{},
				v: types.Node{
					NodeConfig: types.NodeConfig{
						Labels: map[string]string{
							"iaas/region":         "lon",
							"iaas/failure-domain": "lon1",
						},
					},
				},
			},
			"lon",
		},
		{
			nodeContext{
				HeaderContext: HeaderContext{},
				v: types.Node{
					NodeConfig: types.NodeConfig{
						Labels: map[string]string{
							"env": "prod",
						},
					},
				},
			},
			"",
		},
	}

	for _, test := range cases {
		if region := test.context.Region(); region != test.expected {
			t.Errorf("unexpected result.\nexpected:\n%s\ngot:\n%s\n", test.expected, region)
		}
	}
}

func TestFailureDomain(t *testing.T) {
	cases := []struct {
		context  nodeContext
		expected string
	}{
		{
			nodeContext{
				HeaderContext: HeaderContext{},
				v: types.Node{
					NodeConfig: types.NodeConfig{
						Labels: map[string]string{
							"iaas/region":         "lon",
							"iaas/failure-domain": "lon1",
						},
					},
				},
			},
			"lon1",
		},
		{
			nodeContext{
				HeaderContext: HeaderContext{},
				v: types.Node{
					NodeConfig: types.NodeConfig{
						Labels: map[string]string{
							"env": "prod",
						},
					},
				},
			},
			"",
		},
	}

	for _, test := range cases {
		if fd := test.context.FailureDomain(); fd != test.expected {
			t.Errorf("unexpected result.\nexpected:\n%s\ngot:\n%s\n", test.expected, fd)
		}
	}
}
