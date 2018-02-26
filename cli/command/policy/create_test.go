package policy

import (
	"errors"
	"testing"
)

func TestJsonlValidate(t *testing.T) {
	testcases := []struct {
		name    string
		data    []byte
		wantErr error
	}{
		{
			name: "single object jsonl",
			data: []byte(`{"spec":{"group":"baz"}}`),
		},
		{
			name: "single object jsonl with newline",
			data: []byte(`{"spec":{"group":"baz", "namespace": "restricted"}}
`),
		},
		{
			name: "multiple object jsonl",
			data: []byte(`{"spec":{"user":"foo","namespace":"*"}}
				{"spec":{"user":"bar","namespace":"testing"}}
`),
		},
		{
			name:    "invalid object",
			data:    []byte(`{`),
			wantErr: errors.New("unexpected end of JSON input"),
		},
		{
			name: "non JSON line",
			data: []byte(`{"spec": {
			"include":"nested",
			"objects":[
				"and","arrays"
			]
		}}`),
			wantErr: errors.New("unexpected end of JSON input"),
		},
		{
			name:    "empty input",
			data:    []byte(``),
			wantErr: errors.New("empty JSON line input"),
		},
		{
			name:    "array input",
			data:    []byte(`[{"spec":{"user":"foo","namespace":"*"}}]`),
			wantErr: errors.New("expected a json object per line, got a json array"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := jsonlValidate(tc.data)
			if err == nil {
				if tc.wantErr != nil {
					t.Fatalf("unexpected error while validating jsonl content:\n\t(GOT): %v\n\t(WNT): %v", err, tc.wantErr)
				}
			} else {
				if tc.wantErr != nil {
					if err.Error() != tc.wantErr.Error() {
						t.Fatalf("unexpected error while validating jsonl content:\n\t(GOT): %v\n\t(WNT): %v", err, tc.wantErr)
					}
				} else {
					t.Fatalf("unexpected error while validating jsonl content:\n\t(GOT): %v\n\t(WNT): %v", err, tc.wantErr)
				}
			}
		})
	}
}
