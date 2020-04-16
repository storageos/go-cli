// Package flagutil exports utility functions focused on command flag sets.
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

// SupportCAS registers a version constrained compare-and-set flag in flagSet,
// storing the value provided in p. This enables commands using API methods that
// operate on existing resources to take a version string to be used as an
// update constraint. The returned closure indicates whether the flag was
// given a value.
//
// It will replace the flag lookup for "cas".
func SupportCAS(flagSet *pflag.FlagSet, p *string) func() bool {
	const casName = "cas"

	flagSet.StringVar(p, casName, "", "make changes to a resource conditional upon matching the provided version")

	return func() bool {
		return flagSet.Changed(casName)
	}
}

// SupportAsync registers a boolean flag for enabling asynchronous command
// behaviour in flagSet, storing the value provided in p. The default for the
// flag is to not perform requests asynchronously.
//
// It will replace the flag lookup for "async".
func SupportAsync(flagSet *pflag.FlagSet, p *bool) {
	flagSet.BoolVar(p, "async", false, "perform the operation asynchronously, using the configured timeout duration")
}
