package cmd

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/kr/pretty"
	"gopkg.in/yaml.v3"

	"code.storageos.net/storageos/c2-cli/config/file"
)

func TestEncodeExampleConfig(t *testing.T) {
	t.Parallel()

	wantEncoded := fmt.Sprintf(`noAuthCache: "false"
endpoints:
  - http://localhost:5705
cacheDir: %v
timeout: 15s
username: storageos
useIds: "false"
namespace: default
output: text
`, *file.ExampleConfigFile.RawCacheDir)
	w := &bytes.Buffer{}

	if err := yaml.NewEncoder(w).Encode(file.ExampleConfigFile); err != nil {
		t.Fatalf("example config file cannot be encoded to yaml")
	}

	if w.String() != wantEncoded {
		pretty.Ldiff(t, w.String(), wantEncoded)
		t.Fatalf("encoded example config file does not match desired")
	}
}
