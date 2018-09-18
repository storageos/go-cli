package formatter

import (
	"bytes"
	"testing"

	"github.com/storageos/go-api/types"
)

func TestPolicyWrite(t *testing.T) {
	cases := []struct {
		context  Context
		expected string
	}{
		// Quiet format
		{
			Context{Format: NewPolicyFormat(defaultPolicyQuietFormat, true)},
			`policy1ID
policy2ID
policy3ID
`,
		},
		// Table format
		{
			Context{Format: NewPolicyFormat(defaultPolicyTableFormat, false)},
			`ID         USER  GROUP  NAMESPACE
policy1ID  foo1         ns1
policy2ID        grp1   ns4
policy3ID        foo2   ns1
`,
		},
	}

	type policyType struct {
		User            string `json:"user,omitempty"`
		Group           string `json:"group,omitempty"`
		Readonly        bool   `json:"readonly,omitempty"`
		APIGroup        string `json:"apiGroup,omitempty"`
		Resource        string `json:"resource,omitempty"`
		Namespace       string `json:"namespace,omitempty"`
		NonResourcePath string `json:"nonResourcePath,omitempty"`
	}

	policies := []*types.PolicyWithID{
		{
			ID: "policy1ID",
			Policy: types.Policy{
				Spec: policyType{
					User:      "foo1",
					Namespace: "ns1",
				},
			},
		},
		{
			ID: "policy2ID",
			Policy: types.Policy{
				Spec: policyType{
					Group:     "grp1",
					Namespace: "ns4",
				},
			},
		},
		{
			ID: "policy3ID",
			Policy: types.Policy{
				Spec: policyType{
					Group:     "foo2",
					Namespace: "ns1",
					Readonly:  true,
				},
			},
		},
	}

	for _, test := range cases {
		output := bytes.NewBufferString("")
		test.context.Output = output

		if err := PolicyWrite(test.context, policies); err != nil {
			t.Fatalf("unexpected error while writing policy context: %v", err)
		} else {
			if test.expected != output.String() {
				t.Errorf("unexpected result.\n\t(GOT): %v\n\t(WNT): %v", output.String(), test.expected)
			}
		}
	}
}
