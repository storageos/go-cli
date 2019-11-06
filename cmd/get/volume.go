package get

import (
	"errors"
	"io"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/volume"
)

type volumeCommand struct {
	client  GetClient
	display GetDisplayer

	writer io.Writer
}

func (c *volumeCommand) run(cmd *cobra.Command, args []string) error {
	switch len(args) {
	case 0:
		return errors.New("not implemented")
	case 1:
		return c.getVolume(cmd, args)
	default:
		return c.listVolumes(cmd, args)
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

func (c *volumeCommand) listVolumes(_ *cobra.Command, args []string) error {
	requestedMap := map[id.Namespace]map[id.Volume]bool{}
	for _, a := range args {
		ns, uid, err := id.ParseFQVN(a)
		if err != nil {
			return err
		}

		if _, ok := requestedMap[ns]; !ok {
			requestedMap[ns] = map[id.Volume]bool{}
		}

		requestedMap[ns][uid] = true
	}

	requestedVolumes := []*volume.Resource{}
	for ns := range requestedMap {
		nsVolumes, err := c.client.GetNamespaceVolumes(ns)
		if err != nil {
			return err
		}

		for _, v := range nsVolumes {
			if _, ok := requestedMap[ns][v.ID]; ok {
				requestedVolumes = append(requestedVolumes, v)
			}
		}
	}

	return c.display.WriteGetVolumeList(c.writer, requestedVolumes)
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
