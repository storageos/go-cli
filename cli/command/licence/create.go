package licence

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/dnephin/cobra"
	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
)

type createOptions struct {
	filename string
	stdin    bool
}

func newCreateCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	opt := createOptions{}

	cmd := &cobra.Command{
		Use: "create [OPTIONS]",
		Short: `Create a new licence, Either provide the filename of the licence file or write to stdin.
		E.g. "storageos licence create --filename=licence"
		E.g. "cat licence | storageos licence create --stdin"`,
		Args: cli.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreate(cmd, storageosCli, opt)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opt.filename, "filename", "", "Provide the filename of the licence file in json format.")
	flags.BoolVar(&opt.stdin, "stdin", false, "Read licence input from stdin")
	return cmd
}

func runCreate(cmd *cobra.Command, storageosCli *command.StorageOSCli, opt createOptions) error {
	switch {
	case opt.stdin:
		if opt.filename != "" {
			return fmt.Errorf("Please provide stdin or use other methods. (Not both)")
		}
		return runCreateFromStdin(storageosCli, opt)

	case opt.filename != "":
		return runCreateFromFile(storageosCli, opt)

	default:
		return fmt.Errorf(
			"Please provide input files or set the stdin flag\nSee '%s --help'.\n\nUsage:  %s\n\n%s",
			cmd.CommandPath(),
			cmd.UseLine(),
			cmd.Short,
		)
	}
}

func runCreateFromFile(storageosCli *command.StorageOSCli, opt createOptions) error {
	data, err := ioutil.ReadFile(opt.filename)
	if err != nil {
		return err
	}
	return sendKey(storageosCli, data)
}

func runCreateFromStdin(storageosCli *command.StorageOSCli, opt createOptions) error {
	buf, err := ioutil.ReadAll(storageosCli.In())
	if err != nil {
		return fmt.Errorf("failed to read stdin: %s", err)
	}
	return sendKey(storageosCli, buf)
}

func sendKey(storageosCli *command.StorageOSCli, data []byte) error {
	return storageosCli.Client().LicenceCreate(&types.LicenceKeyContainer{Key: string(bytes.TrimSpace(data))})
}
