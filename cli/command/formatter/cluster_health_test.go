package formatter

import (
	"bytes"
	"encoding/json"
	"testing"

	apiTypes "github.com/storageos/go-api/types"
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
			`NODE         CP_STATUS            DP_STATUS
storageos-1  Healthy              Healthy
storageos-2  Degraded (KV_WRITE)  Healthy
storageos-3  Healthy              Degraded (PRESENTATION, RDB)
`,
		},
		{
			// Test control plane table format
			Context{Format: NewClusterHealthFormat(cpClusterHealthTableFormat, false)},
			`NODE         STATUS     KV     KV_WRITE  NATS
storageos-1  Healthy    alive  alive     alive
storageos-2  Not Ready  alive  unknown   alive
storageos-3  Not Ready  alive  alive     alive
`,
		},
		{
			// Test dataplane table format
			Context{Format: NewClusterHealthFormat(dpClusterHealthTableFormat, false)},
			`NODE         STATUS     DFS_INITIATOR  DIRECTOR  PRESENTATION  RDB
storageos-1  Healthy    alive          alive     alive         alive
storageos-2  Not Ready  alive          alive     alive         alive
storageos-3  Not Ready  alive          alive     unknown       unknown
`,
		},
		{
			// Test detailed table format
			Context{Format: NewClusterHealthFormat(detailedClusterHealthTableFormat, false)},
			`NODE         STATUS     KV     KV_WRITE  NATS   DFS_INITIATOR  DIRECTOR  PRESENTATION  RDB
storageos-1  Healthy    alive  alive     alive  alive          alive     alive         alive
storageos-2  Not Ready  alive  unknown   alive  alive          alive     alive         alive
storageos-3  Not Ready  alive  alive     alive  alive          alive     unknown       unknown
`,
		},
	}

	nodes := []*apiTypes.ClusterHealthNode{}
	err := json.Unmarshal([]byte(`[
		{"nodeName": "storageos-1", "submodules": {
			"directfs_initiator": {"status": "alive"},
			"director": {"status": "alive"},
			"kv": {"status": "alive"},
			"kv_write": {"status": "alive"},
			"nats": {"status": "alive"},
			"presentation": {"status": "alive"},
			"rdb": {"status": "alive"}}
		},
		{"nodeName": "storageos-2", "submodules": {
			"directfs_initiator": {"status": "alive"},
			"director": {"status": "alive"},
			"kv": {"status": "alive"},
			"kv_write": {"status": "unknown"},
			"nats": {"status": "alive"},
			"presentation": {"status": "alive"},
			"rdb": {"status": "alive"}}
		},
		{"nodeName": "storageos-3", "submodules": {
			"directfs_initiator": {"status": "alive"},
			"director": {"status": "alive"},
			"kv": {"status": "alive"},
			"kv_write": {"status": "alive"},
			"nats": {"status": "alive"},
			"presentation": {"status": "unknown"},
			"rdb": {"status": "unknown"}}
		}
	]`), &nodes)
	if err != nil {
		t.Fatalf("unexpected json error: %s", err.Error())
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
