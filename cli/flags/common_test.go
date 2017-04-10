package flags

import (
	"path/filepath"
	"testing"

	"github.com/docker/docker/pkg/testutil/assert"
	"github.com/spf13/pflag"
	cliconfig "github.com/storageos/go-cli/cli/config"
)

func _TestCommonOptionsInstallFlags(t *testing.T) {
	flags := pflag.NewFlagSet("testing", pflag.ContinueOnError)
	opts := NewCommonOptions()
	opts.InstallFlags(flags)

	err := flags.Parse([]string{
		"--tlscacert=\"/foo/cafile\"",
		"--tlscert=\"/foo/cert\"",
		"--tlskey=\"/foo/key\"",
	})
	assert.NilError(t, err)
	assert.Equal(t, opts.TLSOptions.CAFile, "/foo/cafile")
	assert.Equal(t, opts.TLSOptions.CertFile, "/foo/cert")
	assert.Equal(t, opts.TLSOptions.KeyFile, "/foo/key")
}

func defaultPath(filename string) string {
	return filepath.Join(cliconfig.Dir(), filename)
}

func _TestCommonOptionsInstallFlagsWithDefaults(t *testing.T) {
	flags := pflag.NewFlagSet("testing", pflag.ContinueOnError)
	opts := NewCommonOptions()
	opts.InstallFlags(flags)

	err := flags.Parse([]string{})
	assert.NilError(t, err)
	assert.Equal(t, opts.TLSOptions.CAFile, defaultPath("ca.pem"))
	assert.Equal(t, opts.TLSOptions.CertFile, defaultPath("cert.pem"))
	assert.Equal(t, opts.TLSOptions.KeyFile, defaultPath("key.pem"))
}
