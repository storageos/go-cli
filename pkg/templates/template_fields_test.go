package templates

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/storageos/go-api/types"
	cliTypes "github.com/storageos/go-cli/types"
)

type TestStruct struct {
	String  string `docs:"bananas"`
	Int     int
	Int64   int64
	Int8    int8
	Uint    uint
	Uint64  uint64
	Uint8   uint8
	Slice   []string `docs:"something something darkside"`
	private string
}

type TestContainer struct {
	TestStruct

	Nested            TestStruct
	NestedPtr         *TestStruct
	NestedSlice       []TestStruct
	NestedPtrSlice    []*TestStruct
	NestedPtrSlicePtr []*TestStruct
}

func TestDescribeFields(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   interface{}
		want []string
	}{
		{
			name: "simple value",
			in:   TestStruct{},
			want: []string{
				"{{ .String }}: bananas",
				"{{ .Int }}",
				"{{ .Int64 }}",
				"{{ .Int8 }}",
				"{{ .Uint }}",
				"{{ .Uint64 }}",
				"{{ .Uint8 }}",
				"{{ .[]Slice }}: something something darkside",
			},
		},
		{
			name: "simple slice",
			in:   []TestStruct{},
			want: []string{
				"{{ .String }}: bananas",
				"{{ .Int }}",
				"{{ .Int64 }}",
				"{{ .Int8 }}",
				"{{ .Uint }}",
				"{{ .Uint64 }}",
				"{{ .Uint8 }}",
				"{{ .[]Slice }}: something something darkside",
			},
		},
		{
			name: "simple pointer",
			in:   &TestStruct{},
			want: []string{
				"{{ .String }}: bananas",
				"{{ .Int }}",
				"{{ .Int64 }}",
				"{{ .Int8 }}",
				"{{ .Uint }}",
				"{{ .Uint64 }}",
				"{{ .Uint8 }}",
				"{{ .[]Slice }}: something something darkside",
			},
		},
		{
			name: "embedded",
			in: struct {
				TestStruct
			}{},
			want: []string{
				"{{ .TestStruct.String }}: bananas",
				"{{ .TestStruct.Int }}",
				"{{ .TestStruct.Int64 }}",
				"{{ .TestStruct.Int8 }}",
				"{{ .TestStruct.Uint }}",
				"{{ .TestStruct.Uint64 }}",
				"{{ .TestStruct.Uint8 }}",
				"{{ .TestStruct.[]Slice }}: something something darkside",
			},
		},
		{
			name: "last field description is used, intermediate ignored",
			in: struct {
				TestStruct `docs:"plátanos"`
			}{},
			want: []string{
				"{{ .TestStruct.String }}: bananas",
				"{{ .TestStruct.Int }}",
				"{{ .TestStruct.Int64 }}",
				"{{ .TestStruct.Int8 }}",
				"{{ .TestStruct.Uint }}",
				"{{ .TestStruct.Uint64 }}",
				"{{ .TestStruct.Uint8 }}",
				"{{ .TestStruct.[]Slice }}: something something darkside",
			},
		},
		{
			name: "embedded contains embedded",
			in: struct {
				TestContainer
			}{},
			want: []string{
				"{{ .TestContainer.TestStruct.String }}: bananas",
				"{{ .TestContainer.TestStruct.Int }}",
				"{{ .TestContainer.TestStruct.Int64 }}",
				"{{ .TestContainer.TestStruct.Int8 }}",
				"{{ .TestContainer.TestStruct.Uint }}",
				"{{ .TestContainer.TestStruct.Uint64 }}",
				"{{ .TestContainer.TestStruct.Uint8 }}",
				"{{ .TestContainer.TestStruct.[]Slice }}: something something darkside",
				"{{ .TestContainer.Nested.String }}: bananas",
				"{{ .TestContainer.Nested.Int }}",
				"{{ .TestContainer.Nested.Int64 }}",
				"{{ .TestContainer.Nested.Int8 }}",
				"{{ .TestContainer.Nested.Uint }}",
				"{{ .TestContainer.Nested.Uint64 }}",
				"{{ .TestContainer.Nested.Uint8 }}",
				"{{ .TestContainer.Nested.[]Slice }}: something something darkside",
				"{{ .TestContainer.NestedPtr.String }}: bananas",
				"{{ .TestContainer.NestedPtr.Int }}",
				"{{ .TestContainer.NestedPtr.Int64 }}",
				"{{ .TestContainer.NestedPtr.Int8 }}",
				"{{ .TestContainer.NestedPtr.Uint }}",
				"{{ .TestContainer.NestedPtr.Uint64 }}",
				"{{ .TestContainer.NestedPtr.Uint8 }}",
				"{{ .TestContainer.NestedPtr.[]Slice }}: something something darkside",
				"{{ .TestContainer.[]NestedSlice.String }}: bananas",
				"{{ .TestContainer.[]NestedSlice.Int }}",
				"{{ .TestContainer.[]NestedSlice.Int64 }}",
				"{{ .TestContainer.[]NestedSlice.Int8 }}",
				"{{ .TestContainer.[]NestedSlice.Uint }}",
				"{{ .TestContainer.[]NestedSlice.Uint64 }}",
				"{{ .TestContainer.[]NestedSlice.Uint8 }}",
				"{{ .TestContainer.[]NestedSlice.[]Slice }}: something something darkside",
				"{{ .TestContainer.[]NestedPtrSlice.String }}: bananas",
				"{{ .TestContainer.[]NestedPtrSlice.Int }}",
				"{{ .TestContainer.[]NestedPtrSlice.Int64 }}",
				"{{ .TestContainer.[]NestedPtrSlice.Int8 }}",
				"{{ .TestContainer.[]NestedPtrSlice.Uint }}",
				"{{ .TestContainer.[]NestedPtrSlice.Uint64 }}",
				"{{ .TestContainer.[]NestedPtrSlice.Uint8 }}",
				"{{ .TestContainer.[]NestedPtrSlice.[]Slice }}: something something darkside",
				"{{ .TestContainer.[]NestedPtrSlicePtr.String }}: bananas",
				"{{ .TestContainer.[]NestedPtrSlicePtr.Int }}",
				"{{ .TestContainer.[]NestedPtrSlicePtr.Int64 }}",
				"{{ .TestContainer.[]NestedPtrSlicePtr.Int8 }}",
				"{{ .TestContainer.[]NestedPtrSlicePtr.Uint }}",
				"{{ .TestContainer.[]NestedPtrSlicePtr.Uint64 }}",
				"{{ .TestContainer.[]NestedPtrSlicePtr.Uint8 }}",
				"{{ .TestContainer.[]NestedPtrSlicePtr.[]Slice }}: something something darkside",
			},
		},
		{
			name: "nested",
			in: struct {
				Child TestStruct
			}{},
			want: []string{
				"{{ .Child.String }}: bananas",
				"{{ .Child.Int }}",
				"{{ .Child.Int64 }}",
				"{{ .Child.Int8 }}",
				"{{ .Child.Uint }}",
				"{{ .Child.Uint64 }}",
				"{{ .Child.Uint8 }}",
				"{{ .Child.[]Slice }}: something something darkside",
			},
		},
		{
			name: "nested pointer",
			in: struct {
				Child *TestStruct
			}{},
			want: []string{
				"{{ .Child.String }}: bananas",
				"{{ .Child.Int }}",
				"{{ .Child.Int64 }}",
				"{{ .Child.Int8 }}",
				"{{ .Child.Uint }}",
				"{{ .Child.Uint64 }}",
				"{{ .Child.Uint8 }}",
				"{{ .Child.[]Slice }}: something something darkside",
			},
		},
		{
			name: "container value",
			in:   TestContainer{},
			want: []string{
				"{{ .TestStruct.String }}: bananas",
				"{{ .TestStruct.Int }}",
				"{{ .TestStruct.Int64 }}",
				"{{ .TestStruct.Int8 }}",
				"{{ .TestStruct.Uint }}",
				"{{ .TestStruct.Uint64 }}",
				"{{ .TestStruct.Uint8 }}",
				"{{ .TestStruct.[]Slice }}: something something darkside",
				"{{ .Nested.String }}: bananas",
				"{{ .Nested.Int }}",
				"{{ .Nested.Int64 }}",
				"{{ .Nested.Int8 }}",
				"{{ .Nested.Uint }}",
				"{{ .Nested.Uint64 }}",
				"{{ .Nested.Uint8 }}",
				"{{ .Nested.[]Slice }}: something something darkside",
				"{{ .NestedPtr.String }}: bananas",
				"{{ .NestedPtr.Int }}",
				"{{ .NestedPtr.Int64 }}",
				"{{ .NestedPtr.Int8 }}",
				"{{ .NestedPtr.Uint }}",
				"{{ .NestedPtr.Uint64 }}",
				"{{ .NestedPtr.Uint8 }}",
				"{{ .NestedPtr.[]Slice }}: something something darkside",
				"{{ .[]NestedSlice.String }}: bananas",
				"{{ .[]NestedSlice.Int }}",
				"{{ .[]NestedSlice.Int64 }}",
				"{{ .[]NestedSlice.Int8 }}",
				"{{ .[]NestedSlice.Uint }}",
				"{{ .[]NestedSlice.Uint64 }}",
				"{{ .[]NestedSlice.Uint8 }}",
				"{{ .[]NestedSlice.[]Slice }}: something something darkside",
				"{{ .[]NestedPtrSlice.String }}: bananas",
				"{{ .[]NestedPtrSlice.Int }}",
				"{{ .[]NestedPtrSlice.Int64 }}",
				"{{ .[]NestedPtrSlice.Int8 }}",
				"{{ .[]NestedPtrSlice.Uint }}",
				"{{ .[]NestedPtrSlice.Uint64 }}",
				"{{ .[]NestedPtrSlice.Uint8 }}",
				"{{ .[]NestedPtrSlice.[]Slice }}: something something darkside",
				"{{ .[]NestedPtrSlicePtr.String }}: bananas",
				"{{ .[]NestedPtrSlicePtr.Int }}",
				"{{ .[]NestedPtrSlicePtr.Int64 }}",
				"{{ .[]NestedPtrSlicePtr.Int8 }}",
				"{{ .[]NestedPtrSlicePtr.Uint }}",
				"{{ .[]NestedPtrSlicePtr.Uint64 }}",
				"{{ .[]NestedPtrSlicePtr.Uint8 }}",
				"{{ .[]NestedPtrSlicePtr.[]Slice }}: something something darkside",
			},
		},
		{
			name: "container pointer",
			in:   &TestContainer{},
			want: []string{
				"{{ .TestStruct.String }}: bananas",
				"{{ .TestStruct.Int }}",
				"{{ .TestStruct.Int64 }}",
				"{{ .TestStruct.Int8 }}",
				"{{ .TestStruct.Uint }}",
				"{{ .TestStruct.Uint64 }}",
				"{{ .TestStruct.Uint8 }}",
				"{{ .TestStruct.[]Slice }}: something something darkside",
				"{{ .Nested.String }}: bananas",
				"{{ .Nested.Int }}",
				"{{ .Nested.Int64 }}",
				"{{ .Nested.Int8 }}",
				"{{ .Nested.Uint }}",
				"{{ .Nested.Uint64 }}",
				"{{ .Nested.Uint8 }}",
				"{{ .Nested.[]Slice }}: something something darkside",
				"{{ .NestedPtr.String }}: bananas",
				"{{ .NestedPtr.Int }}",
				"{{ .NestedPtr.Int64 }}",
				"{{ .NestedPtr.Int8 }}",
				"{{ .NestedPtr.Uint }}",
				"{{ .NestedPtr.Uint64 }}",
				"{{ .NestedPtr.Uint8 }}",
				"{{ .NestedPtr.[]Slice }}: something something darkside",
				"{{ .[]NestedSlice.String }}: bananas",
				"{{ .[]NestedSlice.Int }}",
				"{{ .[]NestedSlice.Int64 }}",
				"{{ .[]NestedSlice.Int8 }}",
				"{{ .[]NestedSlice.Uint }}",
				"{{ .[]NestedSlice.Uint64 }}",
				"{{ .[]NestedSlice.Uint8 }}",
				"{{ .[]NestedSlice.[]Slice }}: something something darkside",
				"{{ .[]NestedPtrSlice.String }}: bananas",
				"{{ .[]NestedPtrSlice.Int }}",
				"{{ .[]NestedPtrSlice.Int64 }}",
				"{{ .[]NestedPtrSlice.Int8 }}",
				"{{ .[]NestedPtrSlice.Uint }}",
				"{{ .[]NestedPtrSlice.Uint64 }}",
				"{{ .[]NestedPtrSlice.Uint8 }}",
				"{{ .[]NestedPtrSlice.[]Slice }}: something something darkside",
				"{{ .[]NestedPtrSlicePtr.String }}: bananas",
				"{{ .[]NestedPtrSlicePtr.Int }}",
				"{{ .[]NestedPtrSlicePtr.Int64 }}",
				"{{ .[]NestedPtrSlicePtr.Int8 }}",
				"{{ .[]NestedPtrSlicePtr.Uint }}",
				"{{ .[]NestedPtrSlicePtr.Uint64 }}",
				"{{ .[]NestedPtrSlicePtr.Uint8 }}",
				"{{ .[]NestedPtrSlicePtr.[]Slice }}: something something darkside",
			},
		},
	}

	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := DescribeFields(tt.in)
			if !reflect.DeepEqual(got, tt.want) {
				t.Logf("got: (%d fields)", len(got))
				for _, f := range got {
					t.Logf("\t%v", f)
				}
				t.Log()

				t.Logf("want: (%d fields)", len(tt.want))
				for _, f := range tt.want {
					t.Logf("\t%v", f)
				}
				t.Log()

				t.Fail()
			}
		})
	}
}

// TestDescribeFields_CommonTypes runs DescribeFields over a selection of our
// types ensuring no panics occur.
//
// It does not validate correctness of the response, just that the user does not
// get a panic when using the --format help.
func TestDescribeFields_CommonTypes(t *testing.T) {
	t.Parallel()

	tests := []interface{}{
		cliTypes.Cluster{},
		&cliTypes.Cluster{},

		types.Node{},
		&types.Node{},
		[]types.Node{},
		[]*types.Node{},

		cliTypes.Node{},
		&cliTypes.Node{},
		[]cliTypes.Node{},
		[]*cliTypes.Node{},

		types.Volume{},
		&types.Volume{},
		[]types.Volume{},
		[]*types.Volume{},

		types.Namespace{},
		&types.Namespace{},
		[]types.Namespace{},
		[]*types.Namespace{},

		types.Policy{},
		&types.Policy{},
		[]types.Policy{},
		[]*types.Policy{},

		types.PolicyWithID{},
		&types.PolicyWithID{},
		[]types.PolicyWithID{},
		[]*types.PolicyWithID{},

		types.Pool{},
		&types.Pool{},
		[]types.Pool{},
		[]*types.Pool{},

		types.Pools{},
		&types.Pools{},
		[]types.Pools{},
		[]*types.Pools{},

		types.Rule{},
		&types.Rule{},
		[]types.Rule{},
		[]*types.Rule{},

		types.Rules{},
		&types.Rules{},
		[]types.Rules{},
		[]*types.Rules{},

		types.User{},
		&types.User{},
		[]types.User{},
		[]*types.User{},
	}

	for _, tt := range tests {
		var tt = tt
		t.Run(fmt.Sprintf("%T", tt), func(t *testing.T) {
			t.Parallel()

			DescribeFields(tt)
		})
	}
}
