package apiclient

import (
	"reflect"
	"testing"

	"github.com/kr/pretty"

	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/policygroup"
)

func TestFilterPolicyGroupsForNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string

		policyGroups []*policygroup.Resource
		names        []string

		wantPolicyGroups []*policygroup.Resource
		wantErr          error
	}{
		{
			name: "dont filter when no names provided",

			policyGroups: []*policygroup.Resource{
				&policygroup.Resource{
					Name: "policyGroup-a",
				},
				&policygroup.Resource{
					Name: "policyGroup-b",
				},
				&policygroup.Resource{
					Name: "policyGroup-c",
				},
			},
			names: nil,

			wantPolicyGroups: []*policygroup.Resource{
				&policygroup.Resource{
					Name: "policyGroup-a",
				},
				&policygroup.Resource{
					Name: "policyGroup-b",
				},
				&policygroup.Resource{
					Name: "policyGroup-c",
				},
			},
			wantErr: nil,
		},
		{
			name: "filters for provided names",

			policyGroups: []*policygroup.Resource{
				&policygroup.Resource{
					Name: "policyGroup-a",
				},
				&policygroup.Resource{
					Name: "policyGroup-b",
				},
				&policygroup.Resource{
					Name: "policyGroup-c",
				},
			},
			names: []string{"policyGroup-a", "policyGroup-c"},

			wantPolicyGroups: []*policygroup.Resource{
				&policygroup.Resource{
					Name: "policyGroup-a",
				},
				&policygroup.Resource{
					Name: "policyGroup-c",
				},
			},
			wantErr: nil,
		},
		{
			name: "errors when a provided name is not present",

			policyGroups: []*policygroup.Resource{
				&policygroup.Resource{
					Name: "policyGroup-a",
				},
				&policygroup.Resource{
					Name: "policyGroup-b",
				},
				&policygroup.Resource{
					Name: "policyGroup-c",
				},
			},
			names: []string{"policyGroup-a", "certainly-a-dave"},

			wantPolicyGroups: nil,
			wantErr: PolicyGroupNotFoundError{
				name: "certainly-a-dave",
			},
		},
	}

	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			t.Log(tt.wantErr)
			gotPolicyGroups, gotErr := filterPolicyGroupsForNames(tt.policyGroups, tt.names...)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
			}

			if !reflect.DeepEqual(gotPolicyGroups, tt.wantPolicyGroups) {
				pretty.Ldiff(t, gotPolicyGroups, tt.wantPolicyGroups)
				t.Errorf("got policyGroups %v, want %v", pretty.Sprint(gotPolicyGroups), pretty.Sprint(tt.wantPolicyGroups))
			}
		})
	}
}

func TestFilterPolicyGroupsForGIDs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string

		policyGroups []*policygroup.Resource
		gids         []id.PolicyGroup

		wantPolicyGroups []*policygroup.Resource
		wantErr          error
	}{
		{
			name: "dont filter when no gids provided",

			policyGroups: []*policygroup.Resource{
				&policygroup.Resource{
					ID: "policyGroup-1",
				},
				&policygroup.Resource{
					ID: "policyGroup-2",
				},
				&policygroup.Resource{
					ID: "policyGroup-3",
				},
			},
			gids: nil,

			wantPolicyGroups: []*policygroup.Resource{
				&policygroup.Resource{
					ID: "policyGroup-1",
				},
				&policygroup.Resource{
					ID: "policyGroup-2",
				},
				&policygroup.Resource{
					ID: "policyGroup-3",
				},
			},
			wantErr: nil,
		},
		{
			name: "filters for provided gids",

			policyGroups: []*policygroup.Resource{
				&policygroup.Resource{
					ID: "policyGroup-1",
				},
				&policygroup.Resource{
					ID: "policyGroup-2",
				},
				&policygroup.Resource{
					ID: "policyGroup-3",
				},
			},
			gids: []id.PolicyGroup{"policyGroup-1", "policyGroup-3"},

			wantPolicyGroups: []*policygroup.Resource{
				&policygroup.Resource{
					ID: "policyGroup-1",
				},
				&policygroup.Resource{
					ID: "policyGroup-3",
				},
			},
			wantErr: nil,
		},
		{
			name: "errors when a provided gid is not present",

			policyGroups: []*policygroup.Resource{
				&policygroup.Resource{
					ID: "policyGroup-1",
				},
				&policygroup.Resource{
					ID: "policyGroup-2",
				},
				&policygroup.Resource{
					ID: "policyGroup-3",
				},
			},
			gids: []id.PolicyGroup{"policyGroup-1", "policyGroup-42"},

			wantPolicyGroups: nil,
			wantErr: PolicyGroupNotFoundError{
				gid: "policyGroup-42",
			},
		},
	}

	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			t.Log(tt.wantErr)
			gotPolicyGroups, gotErr := filterPolicyGroupsForIDs(tt.policyGroups, tt.gids...)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
			}

			if !reflect.DeepEqual(gotPolicyGroups, tt.wantPolicyGroups) {
				pretty.Ldiff(t, gotPolicyGroups, tt.wantPolicyGroups)
				t.Errorf("got policyGroups %v, want %v", pretty.Sprint(gotPolicyGroups), pretty.Sprint(tt.wantPolicyGroups))
			}
		})
	}
}
