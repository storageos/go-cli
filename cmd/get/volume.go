package get

import (
	"errors"
	"io"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/pkg/id"
)

type volumeCommand struct {
	client  GetClient
	display GetDisplayer

	writer io.Writer
}

func (c *volumeCommand) run(cmd *cobra.Command, args []string) error {
	switch len(args) {
	case 1:
		return c.getVolume(cmd, args)
	default:
		return errors.New("not implemented yet")
	}
}

func (c *volumeCommand) getVolume(_ *cobra.Command, args []string) error {
	ns, uid, err := id.ParseFQVN(args[0])
	if err != nil {
		return err
	}

	volume, err := c.client.GetVolume(ns, uid)
	if err != nil {
		return err
	}

	return c.display.WriteGetVolume(c.writer, volume)
}

func newVolume(w io.Writer, client GetClient, display GetDisplayer) *cobra.Command {
	c := &volumeCommand{
		client:  client,
		display: display,

		writer: w,
	}

	cobraCommand := &cobra.Command{
		Aliases: []string{"volumes"},
		Use:     "volume [volume ids...]",
		Short:   "volume retrieves basic information about StorageOS nodes",
		Example: `
$ storageos get volume banana
`,

		RunE: c.run,
	}

	return cobraCommand
}
