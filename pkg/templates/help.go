package templates

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"
)

// FieldUsage provides a usage help message for template variable v based on
// its fields
func FieldUsage(v interface{}) string {
	w := &bytes.Buffer{}
	t := reflect.TypeOf(v)
	writeStructTemplate(w, t, "  ", writeStructFields)
	return w.String()
}

// MethodUsage provides a usage help message for template variable v based on
// its methods
func MethodUsage(v interface{}) string {
	w := &bytes.Buffer{}
	t := reflect.TypeOf(v)
	writeStructTemplate(w, t, "  ", writeStructMethods)
	return w.String()
}

// writeStructTemplate writes the template doc for a struct
func writeStructTemplate(w io.Writer, t reflect.Type, indent string, writeFieldsFunc func(io.Writer, reflect.Type, string, int)) {
	name := typeName(t)
	fmt.Fprintf(w, "Available fields in %s format:\n", name)
	writeFieldsFunc(w, t, indent, 0)
}

func typeName(t reflect.Type) string {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	name := t.Name()
	name = strings.TrimSuffix(name, "Context")
	name = strings.TrimSuffix(name, "Response")
	name = strings.TrimSuffix(name, "Result")
	name = strings.ToLower(name)
	return name
}

// writeStructFields writes template doc recursively, based on struct fields
func writeStructFields(w io.Writer, t reflect.Type, indent string, level int) {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Anonymous || field.PkgPath != "" {
			continue
		}

		ft := field.Type
		for ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}

		fmt.Fprint(w, strings.Repeat(indent, level+1))
		switch ft.Kind() {
		case reflect.Slice:
			fmt.Fprintf(w, "{{index .%s <index>}}", field.Name)
		case reflect.Map:
			fmt.Fprintf(w, "{{index .%s <key>}}", field.Name)
		default:
			fmt.Fprintf(w, "{{.%s}}", field.Name)
		}
		if desc := field.Tag.Get("tdesc"); desc != "" {
			fmt.Fprintf(w, ": %s", desc)
		}
		fmt.Fprintln(w)

		switch ft.Kind() {
		case reflect.Struct:
			writeStructFields(w, ft, indent, level+1)
		case reflect.Slice:
			writeStructFields(w, ft.Elem(), indent, level+1)
		}
	}
}

// writeStructMethods writes template doc recursively, based on struct methods
func writeStructMethods(w io.Writer, t reflect.Type, indent string, level int) {
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		if method.PkgPath != "" {
			continue
		}
		switch method.Name {
		case "MarshalJSON", "FullHeader", "AddHeader":
			continue
		}
		mt := method.Type
		fmt.Fprint(w, strings.Repeat(indent, level+1))
		fmt.Fprintf(w, "{{.%s", method.Name)
		for i := 1; i < mt.NumIn(); i++ {
			fmt.Fprint(w, " ")
			pt := mt.In(i)
			fmt.Fprintf(w, "<%s>", pt.Name())
		}
		fmt.Fprintln(w, "}}")
	}
}
