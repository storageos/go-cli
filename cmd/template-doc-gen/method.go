package main

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

// writeStructMethods writes template doc recursively, based on struct methods
func writeStructMethods(w io.Writer, t reflect.Type, level int) {
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
