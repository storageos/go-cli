package get

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/cmd/runwrappers"
)

type diagnosticsCommand struct {
	config  ConfigProvider
	client  Client
	display Displayer

	outputPath string

	writer io.Writer
}

func (c *diagnosticsCommand) runWithCtx(ctx context.Context, cmd *cobra.Command, _ []string) error {
	// If no output path specified by the user, set the output file path to:
	//
	// 		<working_dir>/diagnostics-<current_time>.gz
	//
	if c.outputPath == "" {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		filename := fmt.Sprintf("diagnostics-%v.gz", time.Now().Format(time.RFC3339Nano))

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

	bundle, err := c.client.GetDiagnostics(ctx)
	if err != nil {
		// If we failed to get a bundle then remove the file we created.
		os.Remove(outputfile.Name())
		return err
	}
	defer bundle.Close()

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
		Example: `
$ storageos get diagnostics

$ storageos get diagnostics --output-file ~/my-diagnostics
`,
		PreRun: func(_ *cobra.Command, _ []string) {
			c.display = SelectDisplayer(c.config)
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			run := runwrappers.RunWithTimeout(c.config)(c.runWithCtx)
			return run(context.Background(), cmd, args)
		},

		SilenceUsage: true,
	}

	cobraCommand.Flags().StringVar(&c.outputPath, "output-file", "", "writes the generated diagnostic bundle to a specified file path")

	return cobraCommand
}
