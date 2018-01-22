package jointools_test

import (
	"testing"

	"github.com/storageos/go-cli/pkg/jointools"
)

func TestVerifyJOINSingleHost(t *testing.T) {
	fixtures := []struct {
		input       string
		expectError bool
	}{
		// FQDN w/wo host & port
		{"google.com", false},
		{"http://google.com", false},
		{"http://google.com:5705", false},
		{"google.com:5705", false},
		{"google.com:123456789", true},      // port too high
		{"google.com:-1", true},             // port too low
		{"google.com:notaportnumber", true}, // port not numeric

		// IP addr w/wo host & port
		{"8.8.8.8", false},
		{"tcp://8.8.8.8:5705", false},
		{"8.8.8.8:5705", false},
		{"tcp://8.8.8.8", false},
		{"8.8.8.8:123456789", true},      // port too high
		{"8.8.8.8:-1", true},             // port too low
		{"8.8.8.8:notaportnumber", true}, // port not numeric

		// local domain name w/wo host&port
		{"localhost", false},
		{"http://localhost", false},
		{"http://localhost:5705", false},
		{"localhost:5705", false},
		{"localhost:123456789", true},      // port too high
		{"localhost:-1", true},             // port too low
		{"localhost:notaportnumber", true}, // port not numeric

		// unresolveable domain name w/wo host&port
		{"unresolveableJunk", true},
		{"http://unresolveableJunk", true},
		{"http://unresolveableJunk:5705", true},
		{"unresolveableJunk:5705", true},
	}

	for _, f := range fixtures {
		errs := jointools.VerifyJOIN(f.input)
		if (errs != nil) != f.expectError {
			t.Errorf("unexpected result. input: %v, errors: %+v (expecting error? %v)", f.input, errs, f.expectError)
		}
	}
}

func TestVerifyJOINMultiHost(t *testing.T) {
	fixtures := []struct {
		input       string
		expectError bool
	}{
		{"8.8.8.8,8.8.4.4", false},
		{"8.8.8.8:5705,8.8.4.4:5705", false},
		{"tcp://8.8.8.8:5705,tcp://8.8.4.4:5705", false},
		{"tcp://8.8.8.8,tcp://8.8.4.4", false},
		{"tcp://8.8.8.8:5705,8.8.4.4", false},
		{"8.8.8.8:5705,8.8.4.4:5705", false},

		{"google.com,facebook.com", false},
		{"google.com,8.8.8.8", false},
		{"8.8.8.8,facebook.com", false},

		{"http://google.com:5705,http://facebook.com:5705", false},
		{"http://google.com:5705,http://8.8.8.8:5705", false},
		{"http://google.com,facebook.com:5705", false},
		{"http://google.com,8.8.8.8:5705", false},

		{"google.com,localhost", false},
		{"localhost,google.com", false},
		{"8.8.8.8,localhost", false},
		{"localhost,8.8.8.8", false},
	}

	for _, f := range fixtures {
		errs := jointools.VerifyJOIN(f.input)
		if (errs != nil) != f.expectError {
			t.Errorf("unexpected result. input: %v, errors: %+v (expecting error? %v)", f.input, errs, f.expectError)
		}
	}
}

func TestVerifyJOINFragment(t *testing.T) {
	fixtures := []struct {
		input       string
		expectError bool
	}{
		// FQDN w/wo host & port
		{"google.com", false},
		{"http://google.com", false},
		{"http://google.com:5705", false},
		{"google.com:5705", false},
		{"google.com:123456789", true},      // port too high
		{"google.com:-1", true},             // port too low
		{"google.com:notaportnumber", true}, // port not numeric

		// IP addr w/wo host & port
		{"8.8.8.8", false},
		{"tcp://8.8.8.8:5705", false},
		{"8.8.8.8:5705", false},
		{"tcp://8.8.8.8", false},
		{"8.8.8.8:123456789", true},      // port too high
		{"8.8.8.8:-1", true},             // port too low
		{"8.8.8.8:notaportnumber", true}, // port not numeric

		// local domain name w/wo host&port
		{"localhost", false},
		{"http://localhost", false},
		{"http://localhost:5705", false},
		{"localhost:5705", false},
		{"localhost:123456789", true},      // port too high
		{"localhost:-1", true},             // port too low
		{"localhost:notaportnumber", true}, // port not numeric

		// unresolveable domain name w/wo host&port
		{"unresolveableJunk", true},
		{"http://unresolveableJunk", true},
		{"http://unresolveableJunk:5705", true},
		{"unresolveableJunk:5705", true},
	}

	for _, f := range fixtures {
		errs := jointools.VerifyJOINFragment(f.input)
		if (errs != nil) != f.expectError {
			t.Errorf("unexpected result. input: %v, errors: %+v (expecting error? %v)", f.input, errs, f.expectError)
		}
	}
}
