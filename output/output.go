// Package output provides a collection of Displayer types. Displayer types
// are used to abstract away the implementation details of how output is
// written for CLI app commands. TODO: Add yaml decoder
package output

// DefaultEncodingIndent is the encoding indent string which consumers of the
// output package can default to when initialising Displayer types.
const DefaultEncodingIndent = "\t"
