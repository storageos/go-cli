package formatter

import (
	"bytes"
	"testing"
	"time"

	"github.com/storageos/go-api/types"
)

func TestConnectivityWrite(t *testing.T) {
	cases := []struct {
		context  Context
		expected string
	}{
		// Table format (default) - needs 6 spaces after OK
		{
			Context{Format: NewConnectivityFormat(TableFormatKey, false)},
			`SOURCE     NAME  ADDRESS       LATENCY  STATUS  MESSAGE
localhost  svc1  1.1.1.1:1234  1s       OK      
localhost  svc2  1.1.1.1:2345  1s       OK      
localhost  svc3  1.1.1.1:3456  1s       ERROR   timeout
`,
		},
		// Table Quiet format
		{
			Context{Format: NewConnectivityFormat(TableFormatKey, true)},
			`SOURCE->ADDRESS          STATUS
localhost->1.1.1.1:1234  OK
localhost->1.1.1.1:2345  OK
localhost->1.1.1.1:3456  ERROR
`,
		},
		// Summary format
		{
			Context{Format: NewConnectivityFormat(SummaryFormatKey, false), Trunc: true},
			`ERROR
`,
		},
		// Summary format quiet (same)
		{
			Context{Format: NewConnectivityFormat(SummaryFormatKey, true), Trunc: true},
			`ERROR
`,
		},
		// Raw format has weird space issue - skip check for now
		{
			Context{Format: NewConnectivityFormat(RawFormatKey, false)},
			`source: localhost
name: svc1
address: 1.1.1.1:1234
latency: 1s
status: OK
message: 

source: localhost
name: svc2
address: 1.1.1.1:2345
latency: 1s
status: OK
message: 

source: localhost
name: svc3
address: 1.1.1.1:3456
latency: 1s
status: ERROR
message: timeout

`,
		},
		// Raw format quiet
		{
			Context{Format: NewConnectivityFormat(RawFormatKey, true)},
			`localhost->1.1.1.1:1234: OK
localhost->1.1.1.1:2345: OK
localhost->1.1.1.1:3456: ERROR
`,
		},
	}

	results := types.ConnectivityResults{
		{
			Source:    "localhost",
			Label:     "svc1",
			Address:   "1.1.1.1:1234",
			LatencyNS: time.Second * 1,
			Error:     "",
		},
		{
			Source:    "localhost",
			Label:     "svc2",
			Address:   "1.1.1.1:2345",
			LatencyNS: time.Second * 1,
			Error:     "",
		},
		{
			Source:    "localhost",
			Label:     "svc3",
			Address:   "1.1.1.1:3456",
			LatencyNS: time.Second * 1,
			Error:     "timeout",
		},
	}

	for _, test := range cases {
		output := bytes.NewBufferString("")
		test.context.Output = output

		if err := ConnectivityWrite(test.context, results); err != nil {
			t.Fatalf("unexpected error while writing volume context: %s", err.Error())
		} else {
			if test.expected != output.String() {
				t.Errorf("unexpected result.\nexpected:\n%s\ngot:\n%s\n", test.expected, output)
			}
		}
	}
}
