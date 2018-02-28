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

func TestVerifyUpdate(t *testing.T) {
	testcases := []struct {
		name       string
		updateOpts updateOptions
		wantErr    error
	}{
		{
			name: "invalid username",
			updateOpts: updateOptions{
				username: "%$#",
			},
			wantErr: errors.New(`Username doesn't follow format "[a-zA-Z0-9]+"`),
		},
		{
			name: "valid username",
			updateOpts: updateOptions{
				username: "foo",
			},
		},
		{
			name: "empty username",
			updateOpts: updateOptions{
				username: "",
			},
		},
		{
			name: "invalid groups",
			updateOpts: updateOptions{
				groups: stringSlice{"grp1", "!@#%", "grp2"},
			},
			wantErr: errors.New(`Group element 1 doesn't follow format "[a-zA-Z0-9]+"`),
		},
		{
			name: "valid groups",
			updateOpts: updateOptions{
				groups: stringSlice{"grp1", "grp2"},
			},
		},
		{
			name: "invalid addGroups",
			updateOpts: updateOptions{
				addGroups: stringSlice{"grp1", "!@#%", "grp2"},
			},
			wantErr: errors.New(`add-group element 1 doesn't follow format "[a-zA-Z0-9]+"`),
		},
		{
			name: "valid addGroups",
			updateOpts: updateOptions{
				addGroups: stringSlice{"grp1", "grp2"},
			},
		},
		{
			name: "invalid removeGroups",
			updateOpts: updateOptions{
				removeGroups: stringSlice{"grp1", "!@#%", "grp2"},
			},
			wantErr: errors.New(`remove-group element 1 doesn't follow format "[a-zA-Z0-9]+"`),
		},
		{
			name: "valid removeGroups",
			updateOpts: updateOptions{
				removeGroups: stringSlice{"grp1", "grp2"},
			},
		},
		{
			name: "invalid role",
			updateOpts: updateOptions{
				role: "foo",
			},
			wantErr: errors.New(`Role must be either "user" or "admin", not "foo"`),
		},
		{
			name: "valid role",
			updateOpts: updateOptions{
				role: "user",
			},
		},
		{
			name: "valid role",
			updateOpts: updateOptions{
				role: "admin",
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := verifyUpdate(tc.updateOpts)
			if err != nil {
				if tc.wantErr != nil {
					if err.Error() != tc.wantErr.Error() {
						t.Errorf("unexpected error while verifying update:\n\t(GOT): %v\n\t(WNT): %v", err, tc.wantErr)
					}
				} else {
					t.Errorf("unexpected error while verifying update:\n\t(GOT): %v\n\t(WNT): %v", err, tc.wantErr)
				}
			} else {
				if tc.wantErr != nil {
					t.Errorf("unexpected error while verifying update:\n\t(GOT): %v\n\t(WNT): %v", err, tc.wantErr)
				}
			}
		})
	}
}
