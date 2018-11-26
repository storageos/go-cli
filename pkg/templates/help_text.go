package templates

import (
	"bytes"
	"text/template"
)

// helpTemplate defines the help text returned when --format help is passed to
// the CLI.
//
// The normal template variable delimitors are used as part of the help text, so
// to print a variable << and >> must be used instead.
var helpTemplate = `
The output of all commands can be fully customised by leveraging the
templating system. The available fields for this command are:

	<< range .Fields >><< . >>
	<< end >>
Fields that contain an array of objects are preceded with "[]", such as
"[]Devices" - arrays are 0-based, and can be indexed directly, or iterated
over. As an example, to access the first element in the array and print the
"Identifier" property:

	{{ (index .Devices 0).Identifier }}

Or to iterate over all the devices:

	{{range .Devices}} Identifier: {{.Identifier}} {{end}}

The templating engine also supports simple transformation functions for some
fields:
	<< range .Functions >>
	{{ << .Name >> << .Args >> }}
		<< .Description >><< end >>

These can be composed to really customise the output, for example to print a
list of replicas for a volume:

	{{ .DeviceDir}} has mounts: 
		{{ range .Devices }}
		{{ index (split .Identifier \"/\") 4}} - {{ upper .Status }}
	{{ end }}

Conditionals, additional functions, variable binding, and string formatting
are also supported for most fields, for more information refer to the Go
templating documentation at https://golang.org/pkg/text/template/

`

// brokenText is returned when HelpText goes wrong.
var brokenText = "Something went wrong! :(\n\nPlease file a bug at https://github.com/storageos/go-cli"

// HelpText generates the appropriate --format help text for in, dynamically
// generating a list of all available template fields.
func HelpText(in interface{}) string {
	// Parse the help text, overriding the normal delimiters as they're used in
	// the actual text
	t, err := template.New("help").Delims("<<", ">>").Parse(helpTemplate)
	if err != nil {
		return brokenText
	}

	type Func struct {
		Name        string
		Args        string
		Description string
	}

	// Define all the available functions, and some usage info
	data := struct {
		Fields    []string
		Functions []Func
	}{
		Fields: DescribeFields(in),
		Functions: []Func{
			{
				Name:        "json",
				Args:        ".Field",
				Description: "JSON encode .Field",
			},
			{
				Name:        "split",
				Args:        ".Field \"separator\"",
				Description: "Split .Field content at separator",
			},
			{
				Name:        "join",
				Args:        ".Field \",\"",
				Description: "Join the .Field array contents with a comma",
			},
			{
				Name:        "upper",
				Args:        ".Field",
				Description: "Convert .Field contents to uppercase",
			},
			{
				Name:        "lower",
				Args:        ".Field",
				Description: "Convert .Field contents to lowercase",
			},
			{
				Name:        "pad",
				Args:        ".Field \".\" <prefix> <suffix>",
				Description: "Pad .Field contents with a dot, with <prefix> and <suffix> number of dots",
			},
			{
				Name:        "truncate",
				Args:        ".Field <length>",
				Description: "Truncate .Field contents to <length> characters",
			},
		},
	}

	buf := &bytes.Buffer{}
	if err := t.Execute(buf, data); err != nil {
		return brokenText
	}

	return buf.String()
}
