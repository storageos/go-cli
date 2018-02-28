package user

import (
	"errors"
	"reflect"
	"testing"
)

func TestVerifyGroupLogic(t *testing.T) {
	testcases := []struct {
		name       string
		updateOpts updateOptions
		wantErr    error
	}{
		{
			name: "groups only",
			updateOpts: updateOptions{
				groups: stringSlice{"grp1", "grp2", "grp3"},
			},
		},
		{
			name: "addGroups and removeGroups unique",
			updateOpts: updateOptions{
				addGroups:    stringSlice{"grp4", "grp5"},
				removeGroups: stringSlice{"grp1", "grp2"},
			},
		},
		{
			name: "addGroups and removeGroups with duplicate",
			updateOpts: updateOptions{
				addGroups:    stringSlice{"grp4", "grp5"},
				removeGroups: stringSlice{"grp5", "grp2"},
			},
			wantErr: errors.New("Cannot add and remove the same group at a time"),
		},
		{
			name: "groups and addGroups",
			updateOpts: updateOptions{
				groups:    stringSlice{"grp1", "grp2", "grp3"},
				addGroups: stringSlice{"grp4", "grp5"},
			},
			wantErr: errors.New("Cannot set both groups and add/remove groups"),
		},
		{
			name: "groups and removeGroups",
			updateOpts: updateOptions{
				groups:       stringSlice{"grp1", "grp2", "grp3"},
				removeGroups: stringSlice{"grp4", "grp5"},
			},
			wantErr: errors.New("Cannot set both groups and add/remove groups"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := verifyGroupLogic(tc.updateOpts)
			if err != nil {
				if tc.wantErr != nil {
					if err.Error() != tc.wantErr.Error() {
						t.Errorf("unexpected error while verifying group logic:\n\t(GOT): %v\n\t(WNT): %v", err, tc.wantErr)
					}
				} else {
					t.Errorf("unexpected error while verifying group logic:\n\t(GOT): %v\n\t(WNT): %v", err, tc.wantErr)
				}
			} else {
				if tc.wantErr != nil {
					t.Errorf("unexpected error while verifying group logic:\n\t(GOT): %v\n\t(WNT): %v", err, tc.wantErr)
				}
			}
		})
	}
}

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
