package licence

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
)

type applyOptions struct {
	filename string
	stdin    bool
}

func newApplyCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	opt := applyOptions{}

	cmd := &cobra.Command{
		Use: "apply [OPTIONS]",
		Short: `Apply a new licence, Either provide the filename of the licence file or write to stdin.
		E.g. "storageos licence apply --filename=licence"
		E.g. "cat licence | storageos licence apply --stdin"`,
		Args: cli.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runApply(cmd, storageosCli, opt)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opt.filename, "filename", "", "Provide the filename of the licence file in json format.")
	flags.BoolVar(&opt.stdin, "stdin", false, "Read licence input from stdin")
	return cmd
}

func runApply(cmd *cobra.Command, storageosCli *command.StorageOSCli, opt applyOptions) error {
	switch {
	case opt.stdin:
		if opt.filename != "" {
			return fmt.Errorf("Please provide stdin or use other methods. (Not both)")
		}
		return runApplyFromStdin(storageosCli, opt)

	case opt.filename != "":
		return runApplyFromFile(storageosCli, opt)

	default:
		return fmt.Errorf(
			"Please provide input files or set the stdin flag\nSee '%s --help'.\n\nUsage:  %s\n\n%s",
			cmd.CommandPath(),
			cmd.UseLine(),
			cmd.Short,
		)
	}
}

func runApplyFromFile(storageosCli *command.StorageOSCli, opt applyOptions) error {
	data, err := ioutil.ReadFile(opt.filename)
	if err != nil {
		return err
	}
	return sendKey(storageosCli, data)
}

func runApplyFromStdin(storageosCli *command.StorageOSCli, opt applyOptions) error {
	buf, err := ioutil.ReadAll(storageosCli.In())
	if err != nil {
		return fmt.Errorf("failed to read stdin: %s", err)
	}
	return sendKey(storageosCli, buf)
}

func sendKey(storageosCli *command.StorageOSCli, data []byte) error {
	return storageosCli.Client().LicenceApply(string(bytes.TrimSpace(data)))
}
