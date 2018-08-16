package formatter

import (
	"bytes"
	"testing"

	apiTypes "github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/types"
)

func TestClusterHealthWrite(t *testing.T) {
	cases := []struct {
		context  Context
		expected string
	}{
		// Quiet format
		{
			Context{Format: NewClusterHealthFormat(defaultClusterHealthQuietFormat, true)},
			`Healthy
Not Ready
Not Ready
`,
		},
		// Table format test cases
		{
			// Test default table format
			Context{Format: NewClusterHealthFormat(defaultClusterHealthTableFormat, false)},
			`NODE         ADDRESS    CP_STATUS            DP_STATUS
storageos-1  127.0.0.1  Healthy              Healthy
storageos-2  127.0.0.1  Degraded (KV_WRITE)  Healthy
storageos-3  127.0.0.1  Healthy              Degraded (FS_DRIVER, FS)
`,
		},
		{
			// Test control plane table format
			Context{Format: NewClusterHealthFormat(cpClusterHealthTableFormat, false)},
			`NODE         ADDRESS    STATUS     KV     KV_WRITE  NATS
storageos-1  127.0.0.1  Healthy    alive  alive     alive
storageos-2  127.0.0.1  Not Ready  alive  unknown   alive
storageos-3  127.0.0.1  Not Ready  alive  alive     alive
`,
		},
		{
			// Test dataplane table format
			Context{Format: NewClusterHealthFormat(dpClusterHealthTableFormat, false)},
			`NODE         ADDRESS    STATUS     DFS_CLIENT  DIRECTOR  FS_DRIVER  FS
storageos-1  127.0.0.1  Healthy    alive       alive     alive      alive
storageos-2  127.0.0.1  Not Ready  alive       alive     alive      alive
storageos-3  127.0.0.1  Not Ready  alive       alive     unknown    unknown
`,
		},
		{
			// Test detailed table format
			Context{Format: NewClusterHealthFormat(detailedClusterHealthTableFormat, false)},
			`NODE         ADDRESS    STATUS     KV     NATS   DFS_CLIENT  DIRECTOR  FS_DRIVER  FS
storageos-1  127.0.0.1  Healthy    alive  alive  alive       alive     alive      alive
storageos-2  127.0.0.1  Not Ready  alive  alive  alive       alive     alive      alive
storageos-3  127.0.0.1  Not Ready  alive  alive  alive       alive     unknown    unknown
`,
		},
	}

	aliveStatus := apiTypes.SubModuleStatus{Status: "alive"}
	unknownStatus := apiTypes.SubModuleStatus{Status: "unknown"}

	nodes := []*types.Node{
		{Name: "storageos-1", AdvertiseAddress: "127.0.0.1",
			Health: struct {
				CP *apiTypes.CPHealthStatus
				DP *apiTypes.DPHealthStatus
			}{
				&apiTypes.CPHealthStatus{KV: aliveStatus, KVWrite: aliveStatus, NATS: aliveStatus, Scheduler: aliveStatus},
				&apiTypes.DPHealthStatus{DirectFSClient: aliveStatus, DirectFSServer: aliveStatus, Director: aliveStatus, FSDriver: aliveStatus, FS: aliveStatus},
			},
		},
		{Name: "storageos-2", AdvertiseAddress: "127.0.0.1",
			Health: struct {
				CP *apiTypes.CPHealthStatus
				DP *apiTypes.DPHealthStatus
			}{
				&apiTypes.CPHealthStatus{KV: aliveStatus, KVWrite: unknownStatus, NATS: aliveStatus, Scheduler: aliveStatus},
				&apiTypes.DPHealthStatus{DirectFSClient: aliveStatus, DirectFSServer: aliveStatus, Director: aliveStatus, FSDriver: aliveStatus, FS: aliveStatus},
			},
		},
		{Name: "storageos-3", AdvertiseAddress: "127.0.0.1",
			Health: struct {
				CP *apiTypes.CPHealthStatus
				DP *apiTypes.DPHealthStatus
			}{
				&apiTypes.CPHealthStatus{KV: aliveStatus, KVWrite: aliveStatus, NATS: aliveStatus, Scheduler: aliveStatus},
				&apiTypes.DPHealthStatus{DirectFSClient: aliveStatus, DirectFSServer: aliveStatus, Director: aliveStatus, FSDriver: unknownStatus, FS: unknownStatus},
			},
		},
	}

	for _, test := range cases {
		output := bytes.NewBufferString("")
		test.context.Output = output

		if err := ClusterHealthWrite(test.context, nodes); err != nil {
			t.Fatalf("unexpected error while writing cluster health context: %s", err.Error())
		} else {
			if test.expected != output.String() {
				t.Errorf("unexpected result.\nexpected:\n%s\ngot:\n%s\n", test.expected, output)
			}
		}

	}
}
