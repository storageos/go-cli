package templates

import (
	"reflect"
	"strings"
)

// A field describes the full path to a singular templatable field and its
// optional description.
//
// A field description is set by placing a struct tag "docs" on the field.
// Descriptions should be less than 80 characters, preferably shorter for purely
// aesthetic reasons.
type field struct {
	Path        []string
	Description string
}

// Copy returns a deep copy of f.
func (f *field) Copy() *field {
	if f == nil {
		return &field{}
	}

	path := make([]string, len(f.Path))
	copy(path, f.Path)

	return &field{
		Path:        path,
		Description: f.Description,
	}
}

// DescribeFields generates a list of fields for in that are available for using
// in a format template. The structure of in is analysed, not the contents - it
// is not safe to assume the returned fields are non-nil.
//
// The returned fields are in the format:
//
//      {{ .Field1.Child }}
//
// If a field is a slice, it is returned as:
//
//      {{ .Field1.[]Child }}
//
// Where "[]" indicates an indexable/rangable slice of elements.
func DescribeFields(in interface{}) []string {
	rt := reflect.ValueOf(in)

	out := []string{}
	for _, field := range walkFields(rt, nil) {
		desc := ""
		if field.Description != "" {
			desc = ": " + field.Description
		}

		out = append(out, "{{ ."+strings.Join(field.Path, ".")+" }}"+desc)
	}

	return out
}

// walkFields recursively dereferences in, and selects the correct analyser
// function for the resulting deferenced type.
func walkFields(in reflect.Value, parent *field) (children []field) {
	in = deref(in)

	switch in.Kind() {
	case reflect.Struct:
		return walkStruct(in, parent)

	case reflect.Slice:
		return walkSlice(in, parent)

	default:
		// Any other types are not recursed into, they are usable template
		// fields
		if parent == nil {
			return []field{}
		}
		return []field{*parent}
	}
}

// walkStruct describes the fields of in, using parent as the field path to in.
//
// walkStruct recursively calls walkFields to build a chain of fields for a
// given input - this is safe as any types that would cause recursion beyond the
// stack limit is a type that needs to a good refactor.
func walkStruct(in reflect.Value, parent *field) []field {
	out := []field{}

	in = deref(in)
	rt := in.Type()

	for i := 0; i < in.NumField(); i++ {
		field := in.Field(i)
		fieldType := rt.Field(i)

		// Skip unexported fields
		if fieldType.PkgPath != "" {
			continue
		}

		// Add this field to the path and recurse
		thisField := parent.Copy()
		thisField.Path = append(thisField.Path, fieldType.Name)
		thisField.Description = fieldType.Tag.Get("docs")

		out = append(out, walkFields(field, thisField)...)
	}

	return out
}

// walkSlice determines the element type the slice is comprised of, and calls
// walkFields for an instantiated element of in.
//
// See recursion notes for walkStruct.
func walkSlice(in reflect.Value, parent *field) []field {
	// Get the type of the slice elements
	elemType := in.Type().Elem()

	// Instansiate a new element
	v := reflect.New(elemType)
	v = deref(v)

	// If a slice is passed directly to walkFields, parent will be nil
	if parent == nil {
		return walkFields(v, &field{
			Path: []string{"[]"},
		})
	}

	// Prepend "[]" to the parent field name§
	i := len(parent.Path) - 1
	tail := "[]" + parent.Path[i]
	parent.Path[i] = tail

	return walkFields(v, parent)
}

// deref recursively dereferences in.
func deref(in reflect.Value) reflect.Value {
	// Recursively dereference any pointers
	t := in.Type()
	for in.Kind() == reflect.Ptr {
		// If the pointer is the zero value (nil) instantiate
		if !in.IsValid() {
			in = reflect.New(t)
		}

		t = t.Elem()
		in = in.Elem()
	}

	// If the resulting value is the zero-value, instantiate an instance so
	// there is something to work with
	if !in.IsValid() {
		in = reflect.New(t)
	}

	if in.Kind() == reflect.Ptr {
		return deref(in)
	}

	return in
}
