package get

import (
	"context"
	"errors"
	"io"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/output/jsonformat"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/volume"
)

var errRequiresNamespace = errors.New("namespace not specified")

type volumeCommand struct {
	config  ConfigProvider
	client  GetClient
	display GetDisplayer

	namespaceID string

	writer io.Writer
}

func (c *volumeCommand) run(cmd *cobra.Command, args []string) error {
	timeout, err := c.config.DialTimeout()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	switch len(args) {
	case 1:
		return c.getVolume(ctx, args)
	default:
		return c.listVolumes(ctx, args)
	}
}

func (c *volumeCommand) getVolume(ctx context.Context, args []string) error {
	uid := id.Volume(args[0])

	if c.namespaceID == "" {
		return errRequiresNamespace
	}

	volume, err := c.client.GetVolume(
		ctx,
		id.Namespace(c.namespaceID),
		uid,
	)
	if err != nil {
		return err
	}

	return c.display.GetVolume(ctx, c.writer, volume)
}

func (c *volumeCommand) listVolumes(ctx context.Context, args []string) error {
	var volumes []*volume.Resource
	var err error

	uids := make([]id.Volume, len(args))
	for i, a := range args {
		uids[i] = id.Volume(a)
	}

	if c.namespaceID != "" {
		volumes, err = c.client.GetNamespaceVolumes(
			ctx,
			id.Namespace(c.namespaceID),
			uids...,
		)
	} else {
		if len(uids) > 0 {
			return errRequiresNamespace
		}
		volumes, err = c.client.GetAllVolumes(ctx)
	}

	if err != nil {
		return err
	}

	return c.display.GetVolumeList(ctx, c.writer, volumes)
}

func newVolume(w io.Writer, client GetClient, config ConfigProvider) *cobra.Command {
	c := &volumeCommand{
		config: config,
		client: client,
		display: jsonformat.NewDisplayer(
			jsonformat.DefaultEncodingIndent,
		),

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

		// If a legitimate error occurs as part of the VERB volume command
		// we don't need to barf the usage template.
		SilenceUsage: true,
	}

	cobraCommand.Flags().StringVarP(&c.namespaceID, "namespace", "n", "", "the id of the namespace to retrieve the volume resources from. if not set all namespaces are included")

	return cobraCommand
}
