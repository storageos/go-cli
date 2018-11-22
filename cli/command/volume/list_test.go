package volume

import (
	"reflect"
	"sort"
	"testing"

	"github.com/storageos/go-api/types"
)

func TestVolumesSort(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		vols []*types.Volume
		want []*types.Volume
	}{
		{
			name: "by namespace",
			vols: []*types.Volume{
				{
					Namespace: "b",
				},
				{
					Namespace: "a",
				},
			},
			want: []*types.Volume{
				{
					Namespace: "a",
				},
				{
					Namespace: "b",
				},
			},
		},
		{
			name: "by name",
			vols: []*types.Volume{
				{
					Name: "b",
				},
				{
					Name: "a",
				},
			},
			want: []*types.Volume{
				{
					Name: "a",
				},
				{
					Name: "b",
				},
			},
		},
		{
			name: "by namespace/name",
			vols: []*types.Volume{
				{
					Namespace: "b",
					Name:      "1",
				},
				{
					Namespace: "b",
					Name:      "2",
				},
				{
					Namespace: "a",
					Name:      "2",
				},
				{
					Namespace: "a",
					Name:      "1",
				},
			},
			want: []*types.Volume{
				{
					Namespace: "a",
					Name:      "1",
				},
				{
					Namespace: "a",
					Name:      "2",
				},
				{
					Namespace: "b",
					Name:      "1",
				},
				{
					Namespace: "b",
					Name:      "2",
				},
			},
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			vols := tc.vols
			sort.Sort(byNamespaceName(vols))
			if !reflect.DeepEqual(vols, tc.want) {
				t.Fatalf("expect %v got %v", tc.want, vols)
			}
		})
	}
}
