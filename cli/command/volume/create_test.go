package volume

import (
	"fmt"
	"strings"
	"testing"
)

func TestParseNamespaceVolume(t *testing.T) {
	fixtures := []struct {
		namespace     string // provided as a flag
		volumename    string // provided as a flag
		positionalArg string // the positional argument

		// the expected outputs
		expectedNamespace  string
		expectedVolumeName string
		errorExpected      bool
	}{
		// Bad states
		{
			"", "", "",
			"", "", true,
		},
		{
			"foo", "", "",
			"", "", true,
		},
		{
			"", "foo", "foo",
			"", "", true,
		},
		{
			"foo", "foo", "foo",
			"", "", true,
		},
		{
			"foo", "", "foo/bar",
			"", "", true,
		},
		{
			"foo", "bar", "foo/bar",
			"", "", true,
		},

		// Good states
		{
			"", "", "foo",
			"default", "foo", false,
		},
		{
			"", "", "notdefault/foo",
			"notdefault", "foo", false,
		},
		{
			"", "", "notdefault/foo",
			"notdefault", "foo", false,
		},
		{
			"notdefault", "", "foo",
			"notdefault", "foo", false,
		},
		{
			"", "foo", "",
			"default", "foo", false,
		},
		{
			"notdefault", "foo", "",
			"notdefault", "foo", false,
		},
	}

	for _, fix := range fixtures {
		n, v, err := parseNamespaceVolume(fix.namespace, fix.volumename, fix.positionalArg)

		fails := []string{}

		if n != fix.expectedNamespace {
			fails = append(fails, fmt.Sprintf("Namespace (%v) doesn't match expected (%v)", n, fix.expectedNamespace))
		}

		if v != fix.expectedVolumeName {
			fails = append(fails, fmt.Sprintf("Volume name (%v) doesn't match expected (%v)", v, fix.expectedVolumeName))
		}

		if (err != nil) != fix.errorExpected {
			fails = append(fails, "Expected no error, got ("+err.Error()+")")
		}

		if len(fails) > 0 {
			t.Logf("Test fixture %+v failed, Reasons: \n\t%v", fix, strings.Join(fails, "\n\t"))
			t.Fail()
		}
	}
}
