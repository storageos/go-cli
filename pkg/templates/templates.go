package templates

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"text/template"
)

// basicFunctions are the set of initial functions provided to every template.
//
// If changing these, be sure to update HelpText()
var basicFunctions = template.FuncMap{
	"json": func(v interface{}) string {
		a, err := json.Marshal(v)
		if err != nil {
			return "<unable to JSON encode this field>"
		}
		return string(a)
	},
	"prettyjson": func(v interface{}) string {
		a, err := json.MarshalIndent(v, "", "\t")
		if err != nil {
			return "<unable to JSON encode this field>"
		}
		return string(a)
	},
	"split":    strings.Split,
	"join":     joinWrap,
	"title":    strings.Title,
	"lower":    strings.ToLower,
	"upper":    strings.ToUpper,
	"pad":      padWith,
	"truncate": truncateWithLength,
}

// joinWrap makes a best attempt to pass arbitrary data to the strings.Join method
// this is needed as the templating system passes in slices of arbitrary data types
// to this method, causing it to fail in the case of anything other than a slice of
// strings.
//
// We try to reflectively iterate over the provided type and pass each element to a
// fmt.Sprint call. We use this method to build a slice of string types, then pass
// this to the originally intended strings.Join method.
func joinWrap(unknown interface{}, joinChar string) string {
	if reflect.TypeOf(unknown).Kind() == reflect.Slice {
		val := reflect.ValueOf(unknown)
		len := val.Len()

		formatted := []string{}
		for i := 0; i < len; i++ {
			asInterface := val.Index(i).Interface()

			// Assume fmt has the best chance of formatting the type
			formatted = append(formatted, fmt.Sprint(asInterface))

		}
		return strings.Join(formatted, joinChar)

	}
	return "<unable to process slice>"

}

// Parse creates a new anonymous template with the basic functions
// and parses the given format.
func Parse(format string) (*template.Template, error) {
	return NewParse("", format)
}

// NewParse creates a new tagged template with the basic functions
// and parses the given format.
func NewParse(tag, format string) (*template.Template, error) {
	return template.New(tag).Funcs(basicFunctions).Parse(format)
}

// padWith adds whitespace to the input if the input is non-empty
func padWith(source, padchar string, prefix, suffix int) string {
	if source == "" {
		return source
	}
	return strings.Repeat(padchar, prefix) + source + strings.Repeat(padchar, suffix)
}

// truncateWithLength truncates the source string up to the length provided by the input
func truncateWithLength(source string, length int) string {
	if len(source) < length {
		return source
	}
	return source[:length]
}
