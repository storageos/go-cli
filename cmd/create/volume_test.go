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

				useCompression: false,
			},
			inputLabels: labels.Set{
				"arbitrary-label": "arbitrary-value",
			},

			wantLabels: labels.Set{
				"arbitrary-label":          "arbitrary-value",
				"storageos.com/nocache":    "true",
				"storageos.com/nocompress": "true",
				"storageos.com/throttle":   "false",
			},
		},
		{
			name: "sets storageos.com/nocompress=false",

			cmd: &volumeCommand{
				useCompression: true,

				useCaching: true, // default of true
			},
			inputLabels: labels.Set{
				"arbitrary-label": "arbitrary-value",
			},

			wantLabels: labels.Set{
				"arbitrary-label":          "arbitrary-value",
				"storageos.com/nocache":    "false",
				"storageos.com/nocompress": "false",
				"storageos.com/throttle":   "false",
			},
		},
		{
			name: "sets storageos.com/throttle=true",

			cmd: &volumeCommand{
				useThrottle: true,

				// defaults of true
				useCaching: true,
			},
			inputLabels: labels.Set{
				"arbitrary-label": "arbitrary-value",
			},

			wantLabels: labels.Set{
				"arbitrary-label":          "arbitrary-value",
				"storageos.com/nocache":    "false",
				"storageos.com/nocompress": "true",
				"storageos.com/throttle":   "true",
			},
		},
		{
			name: "sets storageos.com/replicas > 0",

			cmd: &volumeCommand{
				withReplicas: 5,

				// defaults of true
				useCaching: true,
			},
			inputLabels: labels.Set{
				"arbitrary-label": "arbitrary-value",
			},

			wantLabels: labels.Set{
				"arbitrary-label":          "arbitrary-value",
				"storageos.com/nocache":    "false",
				"storageos.com/nocompress": "true",
				"storageos.com/replicas":   "5",
				"storageos.com/throttle":   "false",
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
