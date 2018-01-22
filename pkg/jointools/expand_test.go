package jointools_test

import (
	"testing"

	"github.com/storageos/go-cli/pkg/jointools"
)

func TestExpandJOINFragment(t *testing.T) {
	fixtures := []struct {
		input  string
		output string
	}{
		// FQDN w/wo host & port
		{"google-public-dns-a.google.com", "http://8.8.8.8:5705"},
		{"http://google-public-dns-a.google.com", "http://8.8.8.8:5705"},
		{"http://google-public-dns-a.google.com:5705", "http://8.8.8.8:5705"},
		{"google-public-dns-a.google.com:5705", "http://8.8.8.8:5705"},

		// IP addr w/wo host & port
		{"8.8.8.8", "http://8.8.8.8:5705"},
		{"tcp://8.8.8.8:5705", "tcp://8.8.8.8:5705"},
		{"8.8.8.8:5705", "http://8.8.8.8:5705"},
		{"tcp://8.8.8.8", "tcp://8.8.8.8:5705"},

		// local domain name w/wo host&port
		{"localhost", "http://127.0.0.1:5705"},
		{"http://localhost", "http://127.0.0.1:5705"},
		{"http://localhost:5705", "http://127.0.0.1:5705"},
		{"localhost:5705", "http://127.0.0.1:5705"},
	}

	for _, f := range fixtures {
		frags := jointools.ExpandJOINFragment(f.input)
		if len(frags) != 1 {
			t.Errorf("unexpected number of endpoints (output: %+v), cluster token?", frags)
		} else {
			if frags[0] != f.output {
				t.Errorf("unexpected result. input: %v, output: %+v (expected: %v)", f.input, frags, f.output)
			}
		}
	}
}

func TestExpandJOINSingleHost(t *testing.T) {
	fixtures := []struct {
		input  string
		output string
	}{
		// FQDN w/wo host & port
		{"google-public-dns-a.google.com", "http://8.8.8.8:5705"},
		{"http://google-public-dns-a.google.com", "http://8.8.8.8:5705"},
		{"http://google-public-dns-a.google.com:5705", "http://8.8.8.8:5705"},
		{"google-public-dns-a.google.com:5705", "http://8.8.8.8:5705"},

		// IP addr w/wo host & port
		{"8.8.8.8", "http://8.8.8.8:5705"},
		{"tcp://8.8.8.8:5705", "tcp://8.8.8.8:5705"},
		{"8.8.8.8:5705", "http://8.8.8.8:5705"},
		{"tcp://8.8.8.8", "tcp://8.8.8.8:5705"},

		// local domain name w/wo host&port
		{"localhost", "http://127.0.0.1:5705"},
		{"http://localhost", "http://127.0.0.1:5705"},
		{"http://localhost:5705", "http://127.0.0.1:5705"},
		{"localhost:5705", "http://127.0.0.1:5705"},
	}

	for _, f := range fixtures {
		out := jointools.ExpandJOIN(f.input)
		if out != f.output {
			t.Errorf("unexpected result. input: %v, output: %v (expected %v)", f.input, out, f.output)
		}
	}
}

func TestExpandJOINMultiHost(t *testing.T) {
	fixtures := []struct {
		input  string
		output string
	}{
		{"8.8.8.8,8.8.4.4", "http://8.8.8.8:5705,http://8.8.4.4:5705"},
		{"8.8.8.8:5705,8.8.4.4:5705", "http://8.8.8.8:5705,http://8.8.4.4:5705"},
		{"tcp://8.8.8.8:5705,tcp://8.8.4.4:5705", "tcp://8.8.8.8:5705,tcp://8.8.4.4:5705"},
		{"tcp://8.8.8.8,tcp://8.8.4.4", "tcp://8.8.8.8:5705,tcp://8.8.4.4:5705"},
		{"tcp://8.8.8.8:5705,8.8.4.4", "tcp://8.8.8.8:5705,http://8.8.4.4:5705"},
		{"8.8.8.8:5705,8.8.4.4:5705", "http://8.8.8.8:5705,http://8.8.4.4:5705"},

		{"google-public-dns-a.google.com,google-public-dns-b.google.com", "http://8.8.8.8:5705,http://8.8.4.4:5705"},
		{"google-public-dns-a.google.com,8.8.4.4", "http://8.8.8.8:5705,http://8.8.4.4:5705"},
		{"8.8.8.8,google-public-dns-b.google.com", "http://8.8.8.8:5705,http://8.8.4.4:5705"},

		{"http://google-public-dns-a.google.com:5705,http://google-public-dns-b.google.com:5705", "http://8.8.8.8:5705,http://8.8.4.4:5705"},
		{"http://google-public-dns-a.google.com:5705,http://8.8.4.4:5705", "http://8.8.8.8:5705,http://8.8.4.4:5705"},
		{"http://google-public-dns-a.google.com,google-public-dns-b.google.com:5705", "http://8.8.8.8:5705,http://8.8.4.4:5705"},
		{"http://google-public-dns-a.google.com,8.8.4.4:5705", "http://8.8.8.8:5705,http://8.8.4.4:5705"},

		{"google-public-dns-a.google.com,localhost", "http://8.8.8.8:5705,http://127.0.0.1:5705"},
		{"localhost,google-public-dns-a.google.com", "http://127.0.0.1:5705,http://8.8.8.8:5705"},
		{"8.8.8.8,localhost", "http://8.8.8.8:5705,http://127.0.0.1:5705"},
		{"localhost,8.8.8.8", "http://127.0.0.1:5705,http://8.8.8.8:5705"},
	}

	for _, f := range fixtures {
		out := jointools.ExpandJOIN(f.input)
		if out != f.output {
			t.Errorf("unexpected result. input: %v, output: %v (expected %v)", f.input, out, f.output)
		}
	}
}
