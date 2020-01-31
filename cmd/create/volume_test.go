package create

import (
	"reflect"
	"testing"

	"github.com/kr/pretty"

	"code.storageos.net/storageos/c2-cli/pkg/labels"
)

func TestSetKnownLabels(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string

		cmd         *volumeCommand
		inputLabels labels.Set

		wantLabels labels.Set
	}{
		{
			name: "sets storageos.com/nocache=true",

			cmd: &volumeCommand{
				useCaching: false,

				useCompression: true, // default of true
			},
			inputLabels: labels.Set{
				"arbitrary-label": "arbitrary-value",
			},

			wantLabels: labels.Set{
				"arbitrary-label":          "arbitrary-value",
				"storageos.com/nocache":    "true",
				"storageos.com/nocompress": "false",
				"storageos.com/throttle":   "false",
			},
		},
		{
			name: "sets storageos.com/nocompress=true",

			cmd: &volumeCommand{
				useCompression: false,

				useCaching: true, // default of true
			},
			inputLabels: labels.Set{
				"arbitrary-label": "arbitrary-value",
			},

			wantLabels: labels.Set{
				"arbitrary-label":          "arbitrary-value",
				"storageos.com/nocache":    "false",
				"storageos.com/nocompress": "true",
				"storageos.com/throttle":   "false",
			},
		},
		{
			name: "sets storageos.com/throttle=true",

			cmd: &volumeCommand{
				useThrottle: true,

				// defaults of true
				useCaching:     true,
				useCompression: true,
			},
			inputLabels: labels.Set{
				"arbitrary-label": "arbitrary-value",
			},

			wantLabels: labels.Set{
				"arbitrary-label":          "arbitrary-value",
				"storageos.com/nocache":    "false",
				"storageos.com/nocompress": "false",
				"storageos.com/throttle":   "true",
			},
		},
		{
			name: "sets storageos.com/replicas > 0",

			cmd: &volumeCommand{
				withReplicas: 5,

				// defaults of true
				useCaching:     true,
				useCompression: true,
			},
			inputLabels: labels.Set{
				"arbitrary-label": "arbitrary-value",
			},

			wantLabels: labels.Set{
				"arbitrary-label":          "arbitrary-value",
				"storageos.com/nocache":    "false",
				"storageos.com/nocompress": "false",
				"storageos.com/replicas":   "5",
				"storageos.com/throttle":   "false",
			},
		},
		{
			name: "sets storageos.com/hint.master when given",

			cmd: &volumeCommand{
				hintMaster: []string{"node-a", "node-b", "node-c"},

				// defaults of true
				useCaching:     true,
				useCompression: true,
			},
			inputLabels: labels.Set{
				"arbitrary-label": "arbitrary-value",
			},

			wantLabels: labels.Set{
				"arbitrary-label":           "arbitrary-value",
				"storageos.com/nocache":     "false",
				"storageos.com/nocompress":  "false",
				"storageos.com/hint.master": "node-a,node-b,node-c",
				"storageos.com/throttle":    "false",
			},
		},
		{
			name: "sets storageos.com/hint.replicas when given",

			cmd: &volumeCommand{
				hintReplicas: []string{"node-a", "node-b", "node-c"},

				// defaults of true
				useCaching:     true,
				useCompression: true,
			},
			inputLabels: labels.Set{
				"arbitrary-label": "arbitrary-value",
			},

			wantLabels: labels.Set{
				"arbitrary-label":             "arbitrary-value",
				"storageos.com/nocache":       "false",
				"storageos.com/nocompress":    "false",
				"storageos.com/hint.replicas": "node-a,node-b,node-c",
				"storageos.com/throttle":      "false",
			},
		},
	}

	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.cmd.setKnownLabels(tt.inputLabels)

			if !reflect.DeepEqual(tt.inputLabels, tt.wantLabels) {
				pretty.Ldiff(t, tt.inputLabels, tt.wantLabels)
				t.Errorf("got modified labels %v, want %v", tt.inputLabels, tt.wantLabels)
			}
		})
	}
}
