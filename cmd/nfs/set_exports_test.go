package nfs

import (
	"reflect"
	"testing"

	"code.storageos.net/storageos/c2-cli/volume"
)

func Test_parseACLs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    []volume.NFSExportConfigACL
		wantErr error
	}{
		{
			name:  "ok",
			input: "cidr;10.0.0.0/8;0;1;root;rw",
			want: []volume.NFSExportConfigACL{
				{
					Identity: volume.NFSExportConfigACLIdentity{
						IdentityType: "cidr",
						Matcher:      "10.0.0.0/8",
					},
					SquashConfig: volume.NFSExportConfigACLSquashConfig{
						GID:    1,
						UID:    0,
						Squash: "root",
					},
					AccessLevel: "rw",
				},
			},
			wantErr: nil,
		},
		{
			name:  "ok 2",
			input: "hostname;*.prod.storageos.com;1000;1001;all;ro",
			want: []volume.NFSExportConfigACL{
				{
					Identity: volume.NFSExportConfigACLIdentity{
						IdentityType: "hostname",
						Matcher:      "*.prod.storageos.com",
					},
					SquashConfig: volume.NFSExportConfigACLSquashConfig{
						GID:    1001,
						UID:    1000,
						Squash: "all",
					},
					AccessLevel: "ro",
				},
			},
			wantErr: nil,
		},

		{
			name:  "ok many",
			input: "cidr;10.0.0.0/8;0;1;root;rw+hostname;*.prod.storageos.com;1000;1001;all;ro",
			want: []volume.NFSExportConfigACL{
				{
					Identity: volume.NFSExportConfigACLIdentity{
						IdentityType: "cidr",
						Matcher:      "10.0.0.0/8",
					},
					SquashConfig: volume.NFSExportConfigACLSquashConfig{
						GID:    1,
						UID:    0,
						Squash: "root",
					},
					AccessLevel: "rw",
				},
				{
					Identity: volume.NFSExportConfigACLIdentity{
						IdentityType: "hostname",
						Matcher:      "*.prod.storageos.com",
					},
					SquashConfig: volume.NFSExportConfigACLSquashConfig{
						GID:    1001,
						UID:    1000,
						Squash: "all",
					},
					AccessLevel: "ro",
				},
			},
			wantErr: nil,
		},

		{
			name:    "wrong identity type",
			input:   "hostn;*.prod.storageos.com;1000;1000;all;ro",
			want:    nil,
			wantErr: errWrongIdentityType,
		},
		{
			name:    "wrong uid",
			input:   "hostname;*.prod.storageos.com;abc;1000;all;ro",
			want:    nil,
			wantErr: errWrongACLUID,
		},
		{
			name:    "wrong gid",
			input:   "hostname;*.prod.storageos.com;1000;abc;all;ro",
			want:    nil,
			wantErr: errWrongACLGID,
		},
		{
			name:    "wrong squash",
			input:   "hostname;*.prod.storageos.com;1000;100;abc;ro",
			want:    nil,
			wantErr: errWrongSquash,
		},
		{
			name:    "wrong access level",
			input:   "hostname;*.prod.storageos.com;1000;100;all;abc",
			want:    nil,
			wantErr: errWrongSquashAccessLevel,
		},
	}
	for _, tt := range tests {
		var tt = tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := parseACLs(tt.input)
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("parseACLs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseACLs() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseExportString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    volume.NFSExportConfig
		wantErr error
	}{
		{
			name:  "ok without ACL",
			input: "1,/,/,",
			want: volume.NFSExportConfig{
				ExportID:   1,
				Path:       "/",
				PseudoPath: "/",
				ACLs:       []volume.NFSExportConfigACL{},
			},
			wantErr: nil,
		},
		{
			name:  "ok with ACL",
			input: "1,/a,/b,cidr;10.0.0.0/8;0;1;root;rw",
			want: volume.NFSExportConfig{
				ExportID:   1,
				Path:       "/a",
				PseudoPath: "/b",
				ACLs: []volume.NFSExportConfigACL{
					{
						Identity: volume.NFSExportConfigACLIdentity{
							IdentityType: "cidr",
							Matcher:      "10.0.0.0/8",
						},
						SquashConfig: volume.NFSExportConfigACLSquashConfig{
							GID:    1,
							UID:    0,
							Squash: "root",
						},
						AccessLevel: "rw",
					},
				},
			},
			wantErr: nil,
		},
		{
			name:    "wrong ID",
			input:   "abc,/a,/b,cidr;10.0.0.0/8;0;0;root;rw",
			want:    volume.NFSExportConfig{},
			wantErr: errWrongExportID,
		},
		{
			name:    "wrong ACL",
			input:   "1,/a,/b,cidr10.0.0.0/800rootrw",
			want:    volume.NFSExportConfig{},
			wantErr: newErrInvalidExportConfigArg("cidr10.0.0.0/800rootrw"),
		},
	}
	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := parseExportString(tt.input)
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("parseExportString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseExportString() got = %v, want %v", got, tt.want)
			}
		})
	}
}
