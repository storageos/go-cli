package system

import (
	"fmt"
	"runtime"
	"time"

	"context"

	"github.com/dnephin/cobra"
	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/pkg/templates"
	"github.com/storageos/go-cli/version"
)

var versionTemplate = `Client:
 Version:      {{.Client.Version}}
 API version:  {{.Client.APIVersion}}
 Go version:   {{.Client.GoVersion}}
 Git commit:   {{.Client.Revision}}
 Built:        {{.Client.BuildDate}}
 OS/Arch:      {{.Client.OS}}/{{.Client.Arch}}{{if .ServerOK}}

Server:
 Version:      {{.Server.Version}}
 API version:  {{.Server.APIVersion}}
 Go version:   {{.Server.GoVersion}}
 Git commit:   {{.Server.Revision}}
 Built:        {{.Server.BuildDate}}
 OS/Arch:      {{.Server.OS}}/{{.Server.Arch}}
 Experimental: {{.Server.Experimental}}{{end}}`

type versionOptions struct {
	format string
}

// NewVersionCommand creates a new cobra.Command for `docker version`
func NewVersionCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	var opts versionOptions

	cmd := &cobra.Command{
		Use:   "version [OPTIONS]",
		Short: "Show the StorageOS version information",
		Args:  cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runVersion(storageosCli, &opts)
		},
	}

	flags := cmd.Flags()

	flags.StringVarP(&opts.format, "format", "f", "", "Format the output using the given Go template")

	return cmd
}

func runVersion(storageosCli *command.StorageOSCli, opts *versionOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	templateFormat := versionTemplate
	if opts.format != "" {
		templateFormat = opts.format
	}

	tmpl, err := templates.Parse(templateFormat)
	if err != nil {
		return cli.StatusError{StatusCode: 64,
			Status: "Template parsing error: " + err.Error()}
	}

	APIVersion := storageosCli.Client().ClientVersion()
	if defaultAPIVersion := storageosCli.DefaultVersion(); APIVersion != defaultAPIVersion {
		APIVersion = fmt.Sprintf("%s (downgraded from %s)", APIVersion, defaultAPIVersion)
	}

	vd := types.VersionResponse{
		Client: &types.VersionInfo{
			Name:       version.ProductName,
			Revision:   version.Revision,
			BuildDate:  version.BuildDate,
			Version:    version.Version,
			APIVersion: version.APIVersion,
			GoVersion:  runtime.Version(),
			OS:         runtime.GOOS,
			Arch:       runtime.GOARCH,
		},
	}

	serverVersion, err := storageosCli.Client().ServerVersion(ctx)
	if err == nil {
		vd.Server = serverVersion
	}

	// first we need to make BuildDate more human friendly
	t, errTime := time.Parse(time.RFC3339Nano, vd.Client.BuildDate)
	if errTime == nil {
		vd.Client.BuildDate = t.Format(time.ANSIC)
	}

	if vd.ServerOK() {
		t, errTime = time.Parse(time.RFC3339Nano, vd.Server.BuildDate)
		if errTime == nil {
			vd.Server.BuildDate = t.Format(time.ANSIC)
		}
	}

	if err2 := tmpl.Execute(storageosCli.Out(), vd); err2 != nil && err == nil {
		err = err2
	}
	storageosCli.Out().Write([]byte{'\n'})
	return err
}
