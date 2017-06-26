package discovery

import "testing"

func Test_extractTokenFromURL(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name      string
		args      args
		wantToken string
		wantErr   bool
	}{
		{
			name:      "expected url",
			args:      args{url: "https://discovery.etcd.io/bf8557cb56223b50ed778c3b0e127f3b"},
			wantToken: "bf8557cb56223b50ed778c3b0e127f3b",
			wantErr:   false,
		},
		{
			name:      "expected url",
			args:      args{url: "http://discovery.etcd.io/bf8557cb56223b50ed778c3b0e127f3b"},
			wantToken: "bf8557cb56223b50ed778c3b0e127f3b",
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken, err := extractTokenFromURL(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractTokenFromURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotToken != tt.wantToken {
				t.Errorf("extractTokenFromURL() = %v, want %v", gotToken, tt.wantToken)
			}
		})
	}
}
