package inspect

import (
	"fmt"

	"github.com/storageos/go-cli/pkg/templates"
)

// formatHelp implements the Inspect interface to print the templates.HelpText
// for the first object passed to Inspect.
//
// Calling flush, or subsequent calls to Inspect are NOPs.
type formatHelp struct {
	printed bool
}

func (f *formatHelp) Inspect(in interface{}, _ []byte) error {
	if f.printed {
		return nil
	}

	fmt.Println(templates.HelpText(in))
	f.printed = true

	return nil
}

func (f *formatHelp) Flush() error {
	return nil
}
