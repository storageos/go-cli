package validation

import (
	"testing"
)

func TestIsValidFSType(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "valid ext2",
			args:    args{value: "ext2"},
			wantErr: false,
		},
		{
			name:    "valid ext3",
			args:    args{value: "ext3"},
			wantErr: false,
		},
		{
			name:    "valid ext4",
			args:    args{value: "ext4"},
			wantErr: false,
		},
		{
			name:    "valid xfs",
			args:    args{value: "xfs"},
			wantErr: false,
		},
		{
			name:    "valid btrfs",
			args:    args{value: "btrfs"},
			wantErr: false,
		},
		{
			name:    "invalid foo",
			args:    args{value: "foo"},
			wantErr: true,
		},
		{
			name:    "invalid random symbol ext2",
			args:    args{value: "ext2%"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := IsValidFSType(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("IsValidFSType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseHostPort(t *testing.T) {
	type args struct {
		host        string
		defaultPort string
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "just IP",
			args:    args{host: "1.1.1.1"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "just IP - default port",
			args:    args{host: "1.1.1.1", defaultPort: "80"},
			want:    "1.1.1.1:80",
			wantErr: false,
		},
		{
			name:    "IP with port",
			args:    args{host: "1.1.1.1:5000", defaultPort: "80"},
			want:    "1.1.1.1:5000",
			wantErr: false,
		},
		{
			name:    "http scheme with IP",
			args:    args{host: "http://1.1.1.1"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "http scheme with IP - default port",
			args:    args{host: "http://1.1.1.1", defaultPort: "80"},
			want:    "1.1.1.1:80",
			wantErr: false,
		},
		{
			name:    "http scheme with IP and port",
			args:    args{host: "http://1.1.1.1:6000", defaultPort: "80"},
			want:    "1.1.1.1:6000",
			wantErr: false,
		},
		{
			name:    "http scheme with IP and port - trailing slash",
			args:    args{host: "http://1.1.1.1:6000/", defaultPort: "80"},
			want:    "1.1.1.1:6000",
			wantErr: false,
		},
		{
			name:    "just hostname",
			args:    args{host: "foo.bar"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "just hostname - default port",
			args:    args{host: "foo.bar", defaultPort: "80"},
			want:    "foo.bar:80",
			wantErr: false,
		},
		{
			name:    "hostname with port",
			args:    args{host: "foo.bar:5000", defaultPort: "80"},
			want:    "foo.bar:5000",
			wantErr: false,
		},
		{
			name:    "http scheme with hostname",
			args:    args{host: "http://foo.bar"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "http scheme with hostname - default port",
			args:    args{host: "http://foo.bar", defaultPort: "80"},
			want:    "foo.bar:80",
			wantErr: false,
		},
		{
			name:    "http scheme with hostname and port",
			args:    args{host: "http://foo.bar:6000", defaultPort: "80"},
			want:    "foo.bar:6000",
			wantErr: false,
		},
		{
			name:    "http scheme with hostname and port - trailing slash",
			args:    args{host: "http://foo.bar:6000/", defaultPort: "80"},
			want:    "foo.bar:6000",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseHostPort(tt.args.host, tt.args.defaultPort)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseHostPort() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseHostPort() = %v, want %v", got, tt.want)
			}
		})
	}
}
