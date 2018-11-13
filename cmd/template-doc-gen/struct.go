package main

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/storageos/go-api/types"
	clitypes "github.com/storageos/go-cli/types"
)

var allFieldStructs = []interface{}{
	clitypes.Cluster{},
	types.ConnectivityResult{},
	types.Licence{},
	types.Namespace{},
	types.Node{},
	types.Policy{},
	types.Pool{},
	types.Rule{},
	types.User{},
	types.Volume{},
	types.VersionResponse{},
}

// writeStructFields writes template doc recursively, based on struct fields
func writeStructFields(w io.Writer, t reflect.Type, level int) {
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
		fmt.Fprintln(w)

		switch ft.Kind() {
		case reflect.Struct:
			writeStructFields(w, ft, level+1)
		case reflect.Slice:
			writeStructFields(w, ft.Elem(), level+1)
		}
	}
}
