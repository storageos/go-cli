package formatter

import "testing"

func TestIsFormat(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name   string
		format Format
		want   bool
	}{
		{
			name:   "empty format",
			format: "",
			want:   true,
		},
		{
			name:   "table format",
			format: "table",
			want:   true,
		},
		{
			name:   "template",
			format: "{{.Name}}",
			want:   true,
		},
		{
			name:   "help with -h",
			format: "-h",
			want:   false,
		},
		{
			name:   "help with --help",
			format: "--help",
			want:   false,
		},
		{
			name:   "help with wrong format",
			format: "wrong_format",
			want:   false,
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if is := IsFormat(tc.format); is != tc.want {
				t.Fatalf("expect %v got %v", tc.want, is)
			}
		})
	}
}
