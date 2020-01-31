// Package flagutil exports utility functions focussed on command flag sets.
// Its purpose is to provide consistency for commands across different packages
// that share a given flag.
package flagutil

import (
	"github.com/spf13/pflag"
)

// SupportSelectors registers a label selector flag in flagSet, storing the
// value provided in p.
//
// It will replace the flag lookup for "selector" and the shorthand lookup for
// "l".
func SupportSelectors(flagSet *pflag.FlagSet, p *[]string) {
	flagSet.StringArrayVarP(p, "selector", "l", []string{}, "filter returned results by a set of comma-separated label selectors")
}
