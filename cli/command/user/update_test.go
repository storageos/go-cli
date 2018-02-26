package user

import (
	"reflect"
	"testing"
)

func TestProcessGroups(t *testing.T) {
	testcases := []struct {
		name       string
		updateOpts updateOptions
		current    []string
		wantGroups []string
	}{
		{
			name: "opts with groups",
			updateOpts: updateOptions{
				groups: stringSlice{"foo", "bar"},
			},
			current:    []string{"aaa", "bbb"},
			wantGroups: []string{"foo", "bar"},
		},
		{
			name: "opts with removeGroups",
			updateOpts: updateOptions{
				removeGroups: []string{"group2", "group3"},
			},
			current:    []string{"group1", "group2", "group3", "group4"},
			wantGroups: []string{"group1", "group4"},
		},
		{
			name: "opts with remove non-existing group",
			updateOpts: updateOptions{
				removeGroups: []string{"groupA", "groupB"},
			},
			current:    []string{"group1", "groupA"},
			wantGroups: []string{"group1"},
		},
		{
			name: "opts with addGroups",
			updateOpts: updateOptions{
				addGroups: []string{"groupA", "groupB"},
			},
			current:    []string{"group1", "group2"},
			wantGroups: []string{"group1", "group2", "groupA", "groupB"},
		},
		{
			name: "opts with add existing group",
			updateOpts: updateOptions{
				addGroups: []string{"groupA", "groupB"},
			},
			current:    []string{"group1", "groupA"},
			wantGroups: []string{"group1", "groupA", "groupB"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotGroups := tc.updateOpts.processGroups(tc.current)
			if !reflect.DeepEqual(gotGroups, tc.wantGroups) {
				t.Fatalf("got unexpected groups while processing groups:\n\t(GOT): %v\n\t(WNT): %v", gotGroups, tc.wantGroups)
			}
		})
	}
}
