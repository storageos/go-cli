package cmdcontext

import (
	"bytes"
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestMinimumTimeoutProvider(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string

		inner   *mockTimeoutProvider
		minimum time.Duration

		wantDuration time.Duration
		wantErr      error
		wantOutput   string
	}{
		{
			name: "ok, use inner when larger than minimum",

			inner: &mockTimeoutProvider{
				CommandTimeoutReturnDuration: 2 * time.Second,
			},
			minimum: time.Second,

			wantDuration: 2 * time.Second,
			wantErr:      nil,
			wantOutput:   "",
		},
		{
			name: "ok, use inner when equals minimum",

			inner: &mockTimeoutProvider{
				CommandTimeoutReturnDuration: time.Second,
			},
			minimum: time.Second,

			wantDuration: time.Second,
			wantErr:      nil,
			wantOutput:   "",
		},
		{
			name: "ok, use minimum when inner is smaller",

			inner: &mockTimeoutProvider{
				CommandTimeoutReturnDuration: time.Second,
			},
			minimum: 2 * time.Second,

			wantDuration: 2 * time.Second,
			wantErr:      nil,
			wantOutput:   "increasing command timeout to 2s\n",
		},
		{
			name: "returns error when inner cannot be retrieved",

			inner: &mockTimeoutProvider{
				CommandTimeoutErr: errors.New("erroneous maximus"),
			},
			minimum: time.Second,

			wantDuration: 0,
			wantErr:      errors.New("erroneous maximus"),
			wantOutput:   "",
		},
	}

	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer

			wrapped := NewMinimumTimeoutProvider(tt.inner, tt.minimum, &buf)

			gotDuration, gotErr := wrapped.CommandTimeout()
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
			}

			if gotDuration != tt.wantDuration {
				t.Errorf("got duration %v, want %v", gotDuration, tt.wantDuration)
			}

			gotOutput := buf.String()

			if gotOutput != tt.wantOutput {
				t.Errorf("got output %q, want %q", gotOutput, tt.wantOutput)
			}
		})
	}
}
