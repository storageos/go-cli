package openapi

import (
	"errors"
	"net/http"
	"testing"
)

func TestGetFilenameFromHeader(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string

		header http.Header

		wantName string
		wantErr  error
	}{
		{
			name: "ok, gets name from correct quoted format amid other values",

			header: http.Header{
				"Content-Disposition": []string{
					"bogus",
					"attachment; filename=\"bundle.bin\"",
				},
			},

			wantName: "bundle.bin",
			wantErr:  nil,
		},
		{
			name: "ok, gets name from badly spaced format amid other values",

			header: http.Header{
				"Content-Disposition": []string{
					"a",
					"attachment; filename= \" bundle.gz \" ",
					"b",
				},
			},

			wantName: "bundle.gz",
			wantErr:  nil,
		},
		{
			name: "err, no attachment header value",

			header: http.Header{
				"Content-Disposition": []string{
					"a",
					"b",
				},
			},

			wantName: "",
			wantErr:  errExtractingFilename,
		},
		{
			name: "err, attachment missing filename=x",

			header: http.Header{
				"Content-Disposition": []string{
					"a",
					"attachment; filename",
					"b",
				},
			},

			wantName: "",
			wantErr:  errExtractingFilename,
		},
	}

	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotName, gotErr := getFilenameFromHeader(tt.header)
			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
			}

			if gotName != tt.wantName {
				t.Errorf("got name %v, want %v", gotName, tt.wantName)
			}
		})
	}
}
