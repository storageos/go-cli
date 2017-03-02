package opts

import (
	"reflect"
	"testing"
)

func TestConvertKVStringsToMap(t *testing.T) {
	input := []string{
		"foo=bar",
		"env=prod",
		"_foobar=foobaz",
		"with.dots=working",
		"and_underscore=working too",
	}

	var expected = map[string]string{
		"foo":            "bar",
		"env":            "prod",
		"_foobar":        "foobaz",
		"with.dots":      "working",
		"and_underscore": "working too",
	}

	output := ConvertKVStringsToMap(input)

	if !reflect.DeepEqual(output, expected) {
		t.Fatal("output not equal to expected")
	}
}
