package formatter

import (
	"testing"
)

func TestWriteLabels(t *testing.T) {
	type args struct {
		labels map[string]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Test write format",
			args{
				labels: map[string]string{
					"iaas/failure-domain":    "europe-west-a",
					"iaas/region":            "europe-west",
					"iaas/size":              "big",
					"storageos.com/replicas": "2",
				},
			},
			"iaas/failure-domain=europe-west-a,iaas/region=europe-west,iaas/size=big,storageos.com/replicas=2",
		},
		{
			"Test label ordering",
			args{
				labels: map[string]string{
					"a/a": "v",
					"b/a": "v",
					"b/c": "v",
					"a/b": "v",
					"d/a": "v",
				},
			},
			"a/a=v,a/b=v,b/a=v,b/c=v,d/a=v",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := writeLabels(tt.args.labels); got != tt.want {
				t.Errorf("writeLabels() = %v, want %v", got, tt.want)
			}
		})
	}
}
