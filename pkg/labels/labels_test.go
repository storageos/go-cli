package labels

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/kr/pretty"
)

func TestNewSetFromPairs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string

		pairs []string

		wantSet Set
		wantErr error
	}{
		{
			name: "ok",

			pairs: []string{
				"storageos.com/replicas=1",
				"my-label=arbitrary-value",
			},

			wantSet: Set{
				"storageos.com/replicas": "1",
				"my-label":               "arbitrary-value",
			},
			wantErr: nil,
		},
		{
			name: "label key conflict returns error",

			pairs: []string{
				"my-repeated-label=value-a",
				"my-repeated-label=value-b",
			},

			wantSet: nil,
			wantErr: fmt.Errorf("%w: my-repeated-label", ErrLabelKeyConflict),
		},
		{
			name: "invalid pair - no key",

			pairs: []string{
				"=some-value",
			},

			wantSet: nil,
			wantErr: fmt.Errorf("%w: =some-value", ErrInvalidLabelFormat),
		},
		{
			name: "invalid pair - no value",

			pairs: []string{
				"some-key=",
			},

			wantSet: nil,
			wantErr: fmt.Errorf("%w: some-key=", ErrInvalidLabelFormat),
		},
	}

	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotSet, gotErr := NewSetFromPairs(tt.pairs)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
			}

			if !reflect.DeepEqual(gotSet, tt.wantSet) {
				pretty.Ldiff(t, gotSet, tt.wantSet)
				t.Errorf("got label set %v, want %v", pretty.Sprint(gotSet), pretty.Sprint(tt.wantSet))
			}
		})
	}
}

func TestSet_String(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		s    Set
		want string
	}{
		{
			name: "empty",
			s:    map[string]string{},
			want: "",
		},
		{
			name: "one elem",
			s: map[string]string{
				"bananaKey": "bananaValue",
			},
			want: "bananaKey=bananaValue",
		},
		{
			name: "two elems",
			s: map[string]string{
				"bananaKey": "bananaValue",
				"kiwiKey":   "kiwiValue",
			},
			want: `bananaKey=bananaValue,
kiwiKey=kiwiValue`,
		},
		{
			name: "three elems",
			s: map[string]string{
				"bananaKey": "bananaValue",
				"kiwiKey":   "kiwiValue",
				"peachKey":  "peachValue",
			},
			want: `bananaKey=bananaValue,
kiwiKey=kiwiValue,
peachKey=peachValue`,
		},
		{
			name: "three elems, different order",
			s: map[string]string{
				"peachKey":  "peachValue",
				"kiwiKey":   "kiwiValue",
				"bananaKey": "bananaValue",
			},
			want: `bananaKey=bananaValue,
kiwiKey=kiwiValue,
peachKey=peachValue`,
		},
	}
	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.s.String(); got != tt.want {
				t.Errorf("String() = \n%v,\n want \n%v", got, tt.want)
			}
		})
	}
}
