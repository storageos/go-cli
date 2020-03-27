package selectors

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/kr/pretty"

	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/pkg/labels"
	"code.storageos.net/storageos/c2-cli/volume"
)

func TestFilterNodes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string

		set        *Set
		inputNodes []*node.Resource

		wantNodes []*node.Resource
	}{
		{
			name: "Set with no selectors returns input as is",

			set: &Set{},
			inputNodes: []*node.Resource{
				{
					Name: "node-a",
				},
				{
					Name: "node-b",
				},
				{
					Name: "node-c",
				},
			},

			wantNodes: []*node.Resource{
				{
					Name: "node-a",
				},
				{
					Name: "node-b",
				},
				{
					Name: "node-c",
				},
			},
		},
		{
			name: "no input matching all selectors",

			set: &Set{
				selectors: []selector{
					func(_ labels.Set) bool {
						return false // Won't match anything
					},
				},
			},
			inputNodes: []*node.Resource{
				{
					Name: "node-a",
				},
				{
					Name: "node-b",
				},
				{
					Name: "node-c",
				},
			},

			wantNodes: []*node.Resource{},
		},
		{
			name: "some input matching all selectors",

			set: &Set{
				selectors: []selector{
					func(set labels.Set) bool {
						v := set["first"]
						return v == "true"
					},
					func(set labels.Set) bool {
						v := set["second"]
						return v == "true"
					},
				},
			},
			inputNodes: []*node.Resource{
				{
					Name: "node-a",
					Labels: labels.Set{
						"first":  "false",
						"second": "true",
					},
				},
				{
					Name: "node-b",
					Labels: labels.Set{
						"first":  "true",
						"second": "true",
					},
				},
				{
					Name: "node-c",
					Labels: labels.Set{
						"first":  "true",
						"second": "true",
					},
				},
				{
					Name: "node-d",
					Labels: labels.Set{
						"first":  "true",
						"second": "false",
					},
				},
			},

			wantNodes: []*node.Resource{
				{
					Name: "node-b",
					Labels: labels.Set{
						"first":  "true",
						"second": "true",
					},
				},
				{
					Name: "node-c",
					Labels: labels.Set{
						"first":  "true",
						"second": "true",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotNodes := tt.set.FilterNodes(tt.inputNodes)
			if !reflect.DeepEqual(gotNodes, tt.wantNodes) {
				pretty.Ldiff(t, gotNodes, tt.wantNodes)
				t.Logf("using selectors %v", pretty.Sprint(tt.set.selectors))
				t.Errorf("got filtered nodes %v, want %v", pretty.Sprint(gotNodes), pretty.Sprint(tt.wantNodes))
			}
		})
	}
}

func TestFilterNamespaces(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string

		set             *Set
		inputNamespaces []*namespace.Resource

		wantNamespaces []*namespace.Resource
	}{
		{
			name: "Set with no selectors returns input as is",

			set: &Set{},
			inputNamespaces: []*namespace.Resource{
				{
					Name: "namespace-a",
				},
				{
					Name: "namespace-b",
				},
				{
					Name: "namespace-c",
				},
			},

			wantNamespaces: []*namespace.Resource{
				{
					Name: "namespace-a",
				},
				{
					Name: "namespace-b",
				},
				{
					Name: "namespace-c",
				},
			},
		},
		{
			name: "no input matching all selectors",

			set: &Set{
				selectors: []selector{
					func(_ labels.Set) bool {
						return false // Won't match anything
					},
				},
			},
			inputNamespaces: []*namespace.Resource{
				{
					Name: "namespace-a",
				},
				{
					Name: "namespace-b",
				},
				{
					Name: "namespace-c",
				},
			},

			wantNamespaces: []*namespace.Resource{},
		},
		{
			name: "some input matching all selectors",

			set: &Set{
				selectors: []selector{
					func(set labels.Set) bool {
						v := set["first"]
						return v == "true"
					},
					func(set labels.Set) bool {
						v := set["second"]
						return v == "true"
					},
				},
			},
			inputNamespaces: []*namespace.Resource{
				{
					Name: "namespace-a",
					Labels: labels.Set{
						"first":  "false",
						"second": "true",
					},
				},
				{
					Name: "namespace-b",
					Labels: labels.Set{
						"first":  "true",
						"second": "true",
					},
				},
				{
					Name: "namespace-c",
					Labels: labels.Set{
						"first":  "true",
						"second": "true",
					},
				},
				{
					Name: "namespace-d",
					Labels: labels.Set{
						"first":  "true",
						"second": "false",
					},
				},
			},

			wantNamespaces: []*namespace.Resource{
				{
					Name: "namespace-b",
					Labels: labels.Set{
						"first":  "true",
						"second": "true",
					},
				},
				{
					Name: "namespace-c",
					Labels: labels.Set{
						"first":  "true",
						"second": "true",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotNamespaces := tt.set.FilterNamespaces(tt.inputNamespaces)
			if !reflect.DeepEqual(gotNamespaces, tt.wantNamespaces) {
				pretty.Ldiff(t, gotNamespaces, tt.wantNamespaces)
				t.Logf("using selectors %v", pretty.Sprint(tt.set.selectors))
				t.Errorf("got filtered namespaces %v, want %v", pretty.Sprint(gotNamespaces), pretty.Sprint(tt.wantNamespaces))
			}
		})
	}
}

func TestFilterVolumes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string

		set          *Set
		inputVolumes []*volume.Resource

		wantVolumes []*volume.Resource
	}{
		{
			name: "Set with no selectors returns input as is",

			set: &Set{},
			inputVolumes: []*volume.Resource{
				{
					Name: "volume-a",
				},
				{
					Name: "volume-b",
				},
				{
					Name: "volume-c",
				},
			},

			wantVolumes: []*volume.Resource{
				{
					Name: "volume-a",
				},
				{
					Name: "volume-b",
				},
				{
					Name: "volume-c",
				},
			},
		},
		{
			name: "no input matching all selectors",

			set: &Set{
				selectors: []selector{
					func(_ labels.Set) bool {
						return false // Won't match anything
					},
				},
			},
			inputVolumes: []*volume.Resource{
				{
					Name: "volume-a",
				},
				{
					Name: "volume-b",
				},
				{
					Name: "volume-c",
				},
			},

			wantVolumes: []*volume.Resource{},
		},
		{
			name: "some input matching all selectors",

			set: &Set{
				selectors: []selector{
					func(set labels.Set) bool {
						v := set["first"]
						return v == "true"
					},
					func(set labels.Set) bool {
						v := set["second"]
						return v == "true"
					},
				},
			},
			inputVolumes: []*volume.Resource{
				{
					Name: "volume-a",
					Labels: labels.Set{
						"first":  "false",
						"second": "true",
					},
				},
				{
					Name: "volume-b",
					Labels: labels.Set{
						"first":  "true",
						"second": "true",
					},
				},
				{
					Name: "volume-c",
					Labels: labels.Set{
						"first":  "true",
						"second": "true",
					},
				},
				{
					Name: "volume-d",
					Labels: labels.Set{
						"first":  "true",
						"second": "false",
					},
				},
			},

			wantVolumes: []*volume.Resource{
				{
					Name: "volume-b",
					Labels: labels.Set{
						"first":  "true",
						"second": "true",
					},
				},
				{
					Name: "volume-c",
					Labels: labels.Set{
						"first":  "true",
						"second": "true",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotVolumes := tt.set.FilterVolumes(tt.inputVolumes)
			if !reflect.DeepEqual(gotVolumes, tt.wantVolumes) {
				pretty.Ldiff(t, gotVolumes, tt.wantVolumes)
				t.Logf("using selectors %v", pretty.Sprint(tt.set.selectors))
				t.Errorf("got filtered volumes %v, want %v", pretty.Sprint(gotVolumes), pretty.Sprint(tt.wantVolumes))
			}
		})
	}
}

func TestNewSetFromStrings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string

		selectors []string

		wantErr error
	}{
		{
			name: "ok - key",

			selectors: []string{"some-key"},

			wantErr: nil,
		},
		{
			name: "ok - pair",

			selectors: []string{"some-key=some-value"},

			wantErr: nil,
		},
		{
			name: "ok - multiple selectors",

			selectors: []string{"some-key", "other-key=another-value"},

			wantErr: nil,
		},
		{
			name: "invalid - empty selector",

			selectors: []string{""},

			wantErr: fmt.Errorf("%w: ", ErrInvalidSelectorFormat),
		},
		{
			name: "invalid - pair with no value",

			selectors: []string{"some-key="},

			wantErr: fmt.Errorf("%w: some-key=", ErrInvalidSelectorFormat),
		},
		{
			name: "invalid - pair with no key",

			selectors: []string{"=some-value"},

			wantErr: fmt.Errorf("%w: =some-value", ErrInvalidSelectorFormat),
		},
		{
			name: "invalid - too many parts",

			selectors: []string{"=invalid-label-format-selector="},

			wantErr: fmt.Errorf("%w: =invalid-label-format-selector=", ErrInvalidSelectorFormat),
		},
		{
			name: "invalid - one of multiple selectors not ok",

			selectors: []string{"some-key", "other-key="},

			wantErr: fmt.Errorf("%w: other-key=", ErrInvalidSelectorFormat),
		},
	}

	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, gotErr := NewSetFromStrings(tt.selectors...)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Logf("selectors: %v", tt.selectors)
				t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestNewSelector(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string

		set   labels.Set
		parts []string

		wantMatched bool
		wantErr     error
	}{
		{
			name: "matches valid label key selector",

			set: labels.Set{
				"some-key": "some-value",
			},
			// some-key
			parts: []string{"some-key"},

			wantMatched: true,
			wantErr:     nil,
		},
		{
			name: "matches valid label pair selector",

			set: labels.Set{
				"some-key": "correct-value",
			},
			// some-key=correct-value
			parts: []string{"some-key", "correct-value"},

			wantMatched: true,
			wantErr:     nil,
		},
		{
			name: "no match for valid label key selector",

			set: labels.Set{
				"some-key": "some-value",
			},
			// missing-key
			parts: []string{"missing-key"},

			wantMatched: false,
			wantErr:     nil,
		},
		{
			name: "no match for valid label pair selector",

			set: labels.Set{
				"some-key": "some-value",
			},
			// some-key=wrong-value
			parts: []string{"some-key", "wrong-value"},

			wantMatched: false,
			wantErr:     nil,
		},
		{
			name: "returns error on invalid selector",

			set: labels.Set{},
			// =invalid-selector=
			parts: []string{"", "invalid-selector", ""},

			wantMatched: false,
			wantErr:     ErrInvalidSelectorFormat,
		},
	}

	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotSelector, gotErr := newSelector(tt.parts)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				// if the construction error does not match up then fail early
				t.Fatalf("got error %v, want %v", gotErr, tt.wantErr)
			}

			if gotSelector == nil && tt.wantErr != nil {
				return
			}

			gotMatched := gotSelector(tt.set)

			if gotMatched != tt.wantMatched {
				t.Logf("label set: %v", pretty.Sprint(tt.set))
				t.Logf("selector: %v", pretty.Sprint(tt.parts))
				t.Errorf("got match %v, want %v", gotMatched, tt.wantMatched)
			}
		})
	}
}
