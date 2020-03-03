package output

import (
	"testing"
)

func TestFormatFromString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		s       string
		want    Format
		wantErr error
	}{
		{
			name:    "json",
			s:       "json",
			want:    JSON,
			wantErr: nil,
		},
		{
			name:    "yaml",
			s:       "yaml",
			want:    YAML,
			wantErr: nil,
		},
		{
			name:    "text",
			s:       "text",
			want:    Text,
			wantErr: nil,
		},
		{
			name:    "text",
			s:       "text",
			want:    Text,
			wantErr: nil,
		},
		{
			name:    "toml",
			s:       "toml",
			want:    Unknown,
			wantErr: errInvalidFormat,
		},
		{
			name:    "xml",
			s:       "xml",
			want:    Unknown,
			wantErr: errInvalidFormat,
		},
	}
	for _, tt := range tests {
		var tt = tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := FormatFromString(tt.s)

			if err != tt.wantErr {
				t.Errorf("FormatFromString() error = %+q, wantErr %+q", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FormatFromString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidFormats(t *testing.T) {
	for _, f := range ValidFormats {
		format, err := FormatFromString(f)
		if err != nil {
			t.Errorf("ValidFormats contains a non valid format: %s", f)
		}

		if format.String() != f {
			t.Errorf("String doesn't return a valid format: %s", f)
		}
	}
}
