package get

import (
	"io"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/volume"
)

type volumeCommand struct {
	client  GetClient
	display GetDisplayer

	namespaceID string

	writer io.Writer
}

func (c *volumeCommand) run(cmd *cobra.Command, args []string) error {
	switch len(args) {
	case 1:
		return c.getVolume(cmd, args)
	default:
		return c.listVolumes(cmd, args)
	}
}

func (c *volumeCommand) getVolume(_ *cobra.Command, args []string) error {
	uid := id.Volume(args[0])

	volume, err := c.client.GetVolume(
		id.Namespace(c.namespaceID),
		uid,
	)
	if err != nil {
		return err
	}

	return c.display.WriteGetVolume(c.writer, volume)
}

func (c *volumeCommand) listVolumes(_ *cobra.Command, args []string) error {
	var volumes []*volume.Resource
	var err error

	uids := make([]id.Volume, len(args))
	for i, a := range args {
		uids[i] = id.Volume(a)
	}

	if c.namespaceID != "" {
		volumes, err = c.client.GetNamespaceVolumes(
			id.Namespace(c.namespaceID),
			uids...,
		)
	} else {
		volumes, err = c.client.GetAllVolumes()
	}

	if err != nil {
		return err
	}

	return c.display.WriteGetVolumeList(c.writer, volumes)
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
		Short:   "volume retrieves basic information about StorageOS volumes",
		Example: `
$ storageos get volume banana
`,

		RunE: c.run,
	}

	cobraCommand.Flags().StringVarP(&c.namespaceID, "namespace", "n", "", "the id of the namespace to retrieve the volume resources from")

	return cobraCommand
}
