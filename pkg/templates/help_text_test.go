package templates

import (
	"testing"

	"github.com/storageos/go-api/types"
)

func TestHelpText(t *testing.T) {
	t.Parallel()

	out := HelpText(&types.Node{})

	if out == brokenText {
		t.Error("got 'broken text' help text")
	}
}
