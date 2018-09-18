package validation

import (
	"testing"
)

func TestGetDeprecations(t *testing.T) {
	fixtures := []struct {
		name          string
		labels        map[string]string
		expectWarning bool
	}{
		{
			name:          "single label",
			labels:        map[string]string{"foo": "bar"},
			expectWarning: false,
		},
		{
			name:          "multiple labels",
			labels:        map[string]string{"foo": "bar", "baz": "bang"},
			expectWarning: false,
		},
		{
			name:          "single label, with special meaning",
			labels:        map[string]string{"storageos.com/replication": "true"},
			expectWarning: false,
		},
		{
			name: "multiple labels, with special meaning",
			labels: map[string]string{
				"storageos.com/replication":   "true",
				"storageos.com/deduplication": "true",
			},
			expectWarning: false,
		},
		{
			name: "multiple labels, some with special meaning",
			labels: map[string]string{
				"foo": "bar",
				"baz": "bang",
				"storageos.com/replication":   "true",
				"storageos.com/deduplication": "true",
			},
			expectWarning: false,
		},
		{
			name:          "single deprecated label",
			labels:        map[string]string{"storageos.feature.nocompress": "true"},
			expectWarning: true,
		},
		{
			name: "multiple deprecated labels",
			labels: map[string]string{
				"storageos.feature.nocompress": "true",
				"storageos.feature.replicas":   "5",
			},
			expectWarning: true,
		},
		{
			name: "multiple labels, some deprecated",
			labels: map[string]string{
				"foo": "bar",
				"baz": "bang",
				"storageos.com/replication":    "true",
				"storageos.com/deduplication":  "true",
				"storageos.feature.nocompress": "true",
				"storageos.feature.replicas":   "5",
			},
			expectWarning: true,
		},
	}

	for _, fix := range fixtures {
		t.Run(fix.name, func(t *testing.T) {
			notices := GetDeprecations(fix.labels)

			if (len(notices) > 0) != fix.expectWarning {
				if len(notices) > 0 {
					t.Errorf("GetDeprecations() returned warnings %+#v for input %+#v", notices, fix.labels)
				} else {
					t.Error("GetDeprecations() didn't return warnings")
				}
			}
		})
	}
}

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
