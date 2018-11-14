package main

import (
	"fmt"

	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli/command/formatter"
	"github.com/storageos/go-cli/pkg/templates"
	clitypes "github.com/storageos/go-cli/types"
)

func main() {
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

	for _, v := range formatter.AllObjects {
		fmt.Println(templates.MethodUsage(v))
	}
	for _, v := range allFieldStructs {
		fmt.Println(templates.FieldUsage(v))
	}
}
