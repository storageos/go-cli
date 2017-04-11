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
