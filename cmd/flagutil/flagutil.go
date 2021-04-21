// Package flagutil exports utility functions focused on command flag sets.
// Its purpose is to provide consistency for commands across different packages
// that share a given flag.
package flagutil

import (
	"fmt"

	"github.com/spf13/cobra"
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

// SupportCASWithRestrictions has the same behaviour of SupportCAS but specifies
// restriction on the helper description.
func SupportCASWithRestrictions(flagSet *pflag.FlagSet, p *string, restrictions map[string]string) func() bool {
	const casName = "cas"

	description := "make changes to a resource conditional upon matching the provided version"

	if len(restrictions) > 0 {
		description += ". Valid only if "
		for k, v := range restrictions {
			description += fmt.Sprintf("%s=%s ", k, v)
		}
	}

	flagSet.StringVar(p, casName, "", description)

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

// WarnAboutValueBeingOverwrittenByK8s adds a required flag that makes sure that the user
// undertstand that the values being set will be overwrtitten by the kubernetes controller
// label synchronisation mechanism
func WarnAboutValueBeingOverwrittenByK8s(flagSet *pflag.FlagSet) {
	// ignore the returned boolean pointer - we don't need to check what the value of this flag is
	_ = flagSet.Bool("i-understand-that-this-value-will-be-overwritten-by-kubernetes", false,
		"this flag must be set to indicate that the user understands the label synchronisation behaviour "+
			"of the kubernetes controller. The flag value (true/false) does not matter.")

	// we know that the flag exist, so we can safely ingore the returned error
	_ = cobra.MarkFlagRequired(flagSet, "i-understand-that-this-value-will-be-overwritten-by-kubernetes")
}
