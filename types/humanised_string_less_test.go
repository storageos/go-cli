package types_test

import "testing"

func TestHumanisedStringLess(t *testing.T) {
	tests := []struct {
		name string
		a, b string
		less bool // expecting a<b
	}{
		// Sort simple integers
		{name: "sort integers 1", a: "0", b: "0", less: true},
		{name: "sort integers 2", a: "0", b: "1", less: true},
		{name: "sort integers 3", a: "1", b: "2", less: true},
		{name: "sort integers 4", a: "0", b: "1", less: false},
		{name: "sort integers 5", a: "345876234", b: "346098235", less: true},
		{name: "sort integers 6", a: "5678972347", b: "346098235", less: false},

		// Sort integer postfixes
		{name: "sort integer postfix 1", a: "foo0", b: "foo0", less: true},
		{name: "sort integer postfix 2", a: "bar0", b: "bar1", less: true},
		{name: "sort integer postfix 3", a: "Baz1", b: "Baz2", less: true},
		{name: "sort integer postfix 4", a: "bAng0", b: "bAng1", less: false},
		{name: "sort integer postfix 5", a: "the red robbin 345876234", b: "the red robbin 346098235", less: true},
		{name: "sort integer postfix 6", a: "sat on the branch 5678972347", b: "sat on the branch 346098235", less: false},

		// Lexicographical sort
		{name: "lexicographical 1", a: "a", b: "a", less: true},
		{name: "lexicographical 2", a: "a", b: "b", less: true},
		{name: "lexicographical 3", a: "b", b: "a", less: true},
		{name: "lexicographical 4", a: "bbbb", b: "bbbb", less: true},
		{name: "lexicographical 5", a: "aaa", b: "aab", less: true},
		{name: "lexicographical 6", a: "aab", b: "aaa", less: true},
		{name: "lexicographical 7", a: "aaaa", b: "aaaaz", less: true},
	}

	for _, test := range tests {
		var tt = test
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
		})
	}

}
