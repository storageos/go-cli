package get

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/c2-cli/cmd/runwrappers"
	"code.storageos.net/storageos/c2-cli/pkg/cmdcontext"
)

type diagnosticsCommand struct {
	config  ConfigProvider
	client  Client
	display Displayer

	outputPath string

	writer io.Writer
}

func (c *diagnosticsCommand) runWithCtx(ctx context.Context, cmd *cobra.Command, _ []string) error {

	// Slightly inefficient in case of fs failure, but the client now needs to
	// care about the server specified name.
	bundle, err := c.client.GetDiagnostics(ctx)
	if err != nil {
		switch v := err.(type) {

		case apiclient.IncompleteDiagnosticsError:
			// If the error is for an incomplete diagnostic bundle, extract
			// the data and warn the user.
			bundle = v.BundleReadCloser()
			fmt.Fprintf(os.Stderr, "\nWarning: %v\n\n", v)

		default:
			// If we failed to get a bundle at all then remove the file we created.
			return err
		}
	}
	defer bundle.Close()

	// If no output path specified by the user, try to use the name provided by
	// the server. Failing that default to:
	//
	// 		<working_dir>/diagnostics-<current_time>.gz
	//
	if c.outputPath == "" {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		var filename string
		if provided, ok := bundle.Named(); ok {
			filename = provided
		} else {
			filename = fmt.Sprintf("diagnostics-%v.gz", time.Now().Format(time.RFC3339Nano))

		}

		c.outputPath = filepath.Join(wd, filename)
	}

	// Create a new file at the target output path, but only if one doesn't
	// already exist.
	//
	// Using the current time in the default filename should provide enough
	// entropy to not be annoying to a user.
	outputfile, err := os.OpenFile(c.outputPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	defer outputfile.Close()

	// Move the bundle to the output path.
	_, err = io.Copy(outputfile, bundle)
	if err != nil {
		return err
	}

	return c.display.GetDiagnostics(ctx, c.writer, c.outputPath)
}

func newDiagnostics(w io.Writer, client Client, config ConfigProvider) *cobra.Command {
	c := &diagnosticsCommand{
		config: config,
		client: client,
		writer: w,
	}
	cobraCommand := &cobra.Command{
		Use:   "diagnostics",
		Short: "Fetch a cluster diagnostic bundle",
		Long:  "Fetch a cluster diagnostic bundle from the target node. Due to the work involved this command will run with a minimum command timeout duration of 1h, although accepts longer durations",
		Example: `
$ storageos get diagnostics

$ storageos get diagnostics --output-file ~/my-diagnostics.gz
`,
		PreRun: func(_ *cobra.Command, _ []string) {
			c.display = SelectDisplayer(c.config)
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			run := runwrappers.Chain(
				runwrappers.RunWithTimeout(
					// Diagnostic retrieval may take a while - force a minimum
					// command timeout of an hour (warning if changed from the
					// user's value)
					cmdcontext.NewMinimumTimeoutProvider(
						c.config,
						time.Hour,
						os.Stderr,
					),
				),
				runwrappers.AuthenticateClient(c.config, c.client),
			)(c.runWithCtx)
			return run(context.Background(), cmd, args)
		},

		SilenceUsage: true,
	}

	cobraCommand.Flags().StringVar(&c.outputPath, "output-file", "", "writes the generated diagnostic bundle to a specified file path")

	return cobraCommand
}
