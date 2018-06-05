package logs

import "testing"

func TestSkipNodeLog(t *testing.T) {
	testcases := map[string]struct {
		nodes    []string
		message  []byte
		wantSkip bool
	}{
		"empty message": {
			nodes:    []string{"storageos-1-80251"},
			message:  []byte(""),
			wantSkip: true,
		},
		"not skip message": {
			nodes:    []string{"storageos-1-80251", "storageos-2-80251"},
			message:  []byte(`{"category":"streamer","host":"storageos-2-80251","level":"debug","module":"logger","msg":"starting remote"}`),
			wantSkip: false,
		},
		"skip message": {
			nodes:    []string{"storageos-1.80251"},
			message:  []byte(`{"category":"streamer","host":"storageos-2-80251","level":"debug","module":"logger","msg":"starting remote"}`),
			wantSkip: true,
		},
	}

	for k, tc := range testcases {
		t.Run(k, func(t *testing.T) {
			gotSkip := skipNodeLog(tc.nodes, tc.message)
			if tc.wantSkip != gotSkip {
				t.Errorf("unexpected skipNodeLog result: \n\t(WNT): %t\n\t(GOT): %t", tc.wantSkip, gotSkip)
			}
		})
	}
}
