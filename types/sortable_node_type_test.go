package types_test

import (
	"testing"

	apiTypes "github.com/storageos/go-api/types"
	cliTypes "github.com/storageos/go-cli/types"
)

func TestSortableNodeTypeCLINodes(t *testing.T) {
	// Utility functions to instansiate a fake node object(s) with only the name
	// field populated
	n := func(name string) *cliTypes.Node { return &cliTypes.Node{Name: name} }
	ns := func(names ...string) []*cliTypes.Node {
		rtn := []*cliTypes.Node{}
		for _, name := range names {
			rtn = append(rtn, n(name))
		}
		return rtn
	}

	fixtures := []struct {
		name          string
		nodes         []*cliTypes.Node
		expectError   bool
		expectedOrder []string
	}{
		{
			name:          "just ints, ordered",
			nodes:         ns("0", "1", "2", "3"),
			expectError:   false,
			expectedOrder: []string{"0", "1", "2", "3"},
		},
		{
			name:          "just ints, random",
			nodes:         ns("1", "3", "2", "0"),
			expectError:   false,
			expectedOrder: []string{"0", "1", "2", "3"},
		},
		{
			name:          "just ints, reversed",
			nodes:         ns("3", "2", "1", "0"),
			expectError:   false,
			expectedOrder: []string{"0", "1", "2", "3"},
		},
		{
			name:          "with prefix, ordered",
			nodes:         ns("host-prefix0", "host-prefix1", "host-prefix2", "host-prefix3"),
			expectError:   false,
			expectedOrder: []string{"host-prefix0", "host-prefix1", "host-prefix2", "host-prefix3"},
		},
		{
			name:          "with prefix, random",
			nodes:         ns("host-prefix2", "host-prefix1", "host-prefix0", "host-prefix3"),
			expectError:   false,
			expectedOrder: []string{"host-prefix0", "host-prefix1", "host-prefix2", "host-prefix3"},
		},
		{
			name:          "with prefix, reversed",
			nodes:         ns("host-prefix3", "host-prefix2", "host-prefix1", "host-prefix0"),
			expectError:   false,
			expectedOrder: []string{"host-prefix0", "host-prefix1", "host-prefix2", "host-prefix3"},
		},
	}

	for _, fixture := range fixtures {
		if err := cliTypes.SortCLINodes(cliTypes.ByNodeName, fixture.nodes); (err != nil) != fixture.expectError {
			t.Fatalf("Got error %e, expected error? %t", err, fixture.expectError)
		}

		// Slice should now be sorted
		// lets check if the order is right
		for i, node := range fixture.nodes {
			if given, expected := node.Name, fixture.expectedOrder[i]; given != expected {
				t.Fatalf("Fixture %s failed at index %d: %s != %s", fixture.name, i, given, expected)
			}
		}
	}
}

func TestSortableNodeTypeAPINodes(t *testing.T) {
	// Utility functions to instansiate a fake node object(s) with only the name
	// field populated
	n := func(name string) *apiTypes.Node { return &apiTypes.Node{Name: name} }
	ns := func(names ...string) []*apiTypes.Node {
		rtn := []*apiTypes.Node{}
		for _, name := range names {
			rtn = append(rtn, n(name))
		}
		return rtn
	}

	fixtures := []struct {
		name          string
		nodes         []*apiTypes.Node
		expectError   bool
		expectedOrder []string
	}{
		{
			name:          "just ints, ordered",
			nodes:         ns("0", "1", "2", "3"),
			expectError:   false,
			expectedOrder: []string{"0", "1", "2", "3"},
		},
		{
			name:          "just ints, random",
			nodes:         ns("1", "3", "2", "0"),
			expectError:   false,
			expectedOrder: []string{"0", "1", "2", "3"},
		},
		{
			name:          "just ints, reversed",
			nodes:         ns("3", "2", "1", "0"),
			expectError:   false,
			expectedOrder: []string{"0", "1", "2", "3"},
		},
		{
			name:          "with prefix, ordered",
			nodes:         ns("host-prefix0", "host-prefix1", "host-prefix2", "host-prefix3"),
			expectError:   false,
			expectedOrder: []string{"host-prefix0", "host-prefix1", "host-prefix2", "host-prefix3"},
		},
		{
			name:          "with prefix, random",
			nodes:         ns("host-prefix2", "host-prefix1", "host-prefix0", "host-prefix3"),
			expectError:   false,
			expectedOrder: []string{"host-prefix0", "host-prefix1", "host-prefix2", "host-prefix3"},
		},
		{
			name:          "with prefix, reversed",
			nodes:         ns("host-prefix3", "host-prefix2", "host-prefix1", "host-prefix0"),
			expectError:   false,
			expectedOrder: []string{"host-prefix0", "host-prefix1", "host-prefix2", "host-prefix3"},
		},
		{
			name:          "caused panic",
			nodes:         ns("kind-v1.14.3-worker", "kind-v1.14.3-worker2"),
			expectError:   false,
			expectedOrder: []string{"kind-v1.14.3-worker", "kind-v1.14.3-worker2"},
		},
		{
			name:          "caused panic reversed",
			nodes:         ns("kind-v1.14.3-worker2", "kind-v1.14.3-worker"),
			expectError:   false,
			expectedOrder: []string{"kind-v1.14.3-worker", "kind-v1.14.3-worker2"},
		},
	}

	for _, fixture := range fixtures {
		if err := cliTypes.SortAPINodes(cliTypes.ByNodeName, fixture.nodes); (err != nil) != fixture.expectError {
			t.Fatalf("Got error %e, expected error? %t", err, fixture.expectError)
		}

		// Slice should now be sorted
		// lets check if the order is right
		for i, node := range fixture.nodes {
			if given, expected := node.Name, fixture.expectedOrder[i]; given != expected {
				t.Fatalf("Fixture %s failed at index %d: %s != %s", fixture.name, i, given, expected)
			}
		}
	}
}
