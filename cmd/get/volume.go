package get

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/apiclient"
	"code.storageos.net/storageos/c2-cli/output/jsonformat"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/volume"
)

type volumeCommand struct {
	config  ConfigProvider
	client  GetClient
	display GetDisplayer

	usingIDs bool

	writer io.Writer
}

func (c *volumeCommand) run(cmd *cobra.Command, args []string) error {
	timeout, err := c.config.CommandTimeout()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	switch len(args) {
	case 1:
		v, err := c.getVolume(ctx, args)
		if err != nil {
			return err
		}

		return c.display.GetVolume(ctx, c.writer, v)
	case 0:
		volumes, err := c.client.GetAllVolumes(ctx)
		if err != nil {
			return err
		}

		return c.display.GetVolumeList(ctx, c.writer, volumes)
	default:
		if c.usingIDs {
			volumes, err := c.listVolumesUsingID(ctx, args)
			if err != nil {
				return err
			}

			return c.display.GetVolumeList(ctx, c.writer, volumes)
		}

		volumes, err := c.listVolumes(ctx, args)
		if err != nil {
			return err
		}

		return c.display.GetVolumeList(ctx, c.writer, volumes)
	}
}

func (c *volumeCommand) getVolume(ctx context.Context, args []string) (*volume.Resource, error) {
	volumeReference := args[0]

	if !c.usingIDs {
		nsName, volName, err := volume.ParseReferenceName(args[0])
		if err != nil {
			return nil, err
		}

		ns, err := c.client.GetNamespaceByName(ctx, nsName)
		if err != nil {
			return nil, err
		}

		return c.client.GetVolumeByName(ctx, ns.ID, volName)
	}

	nsID, volID, err := volume.ParseReferenceID(volumeReference)
	switch err {
	case nil:
	case volume.ErrNoNamespace:
		// if no namespace is supplied then resolve the id of the default one
		defaultNs, err := c.client.GetNamespaceByName(ctx, "default")
		if err != nil {
			return nil, err
		}
		nsID = defaultNs.ID
	default:
		return nil, err
	}

	return c.client.GetVolume(ctx, nsID, volID)
}

// listVolumes retrieves a list of volume resources from the provided set of
// name-based reference strings using the API client.
func (c *volumeCommand) listVolumes(ctx context.Context, nameRefs []string) ([]*volume.Resource, error) {
	nsVols := map[string][]string{}
	nsIDForName := map[string]id.Namespace{}

	for _, ref := range nameRefs {
		nsName, volName, err := volume.ParseReferenceName(ref)

		if err != nil {
			return nil, err
		}

		nsVols[nsName] = append(nsVols[nsName], volName)
	}

	// Get the namespace ID → name mapping to identify requested
	// volumes.
	namespaces, err := c.client.GetAllNamespaces(ctx)
	if err != nil {
		return nil, err
	}

	for _, ns := range namespaces {
		if _, ok := nsVols[ns.Name]; ok {
			nsIDForName[ns.Name] = ns.ID
		}
	}

	resources := []*volume.Resource{}

	for nsName, volNames := range nsVols {
		nsID, ok := nsIDForName[nsName]
		if !ok {
			return nil, errors.New("namespace with name %v not found")
		}

		nsResources, err := c.client.GetNamespaceVolumesByName(ctx, nsID, volNames...)
		if err != nil {
			return nil, err
		}

		resources = append(resources, nsResources...)
	}

	return resources, nil
}

// listVolumesUsingID retrieves a list of volume resources from the provided
// set of ID-based reference strings using the API client.
func (c *volumeCommand) listVolumesUsingID(ctx context.Context, idRefs []string) ([]*volume.Resource, error) {
	nsVols := map[id.Namespace][]id.Volume{}
	defaultNsVols := []id.Volume{}

	for _, ref := range idRefs {
		nsID, volID, err := volume.ParseReferenceID(ref)

		switch err {
		case nil:
			nsVols[nsID] = append(nsVols[nsID], volID)
		case volume.ErrNoNamespace:
			defaultNsVols = append(defaultNsVols, volID)
		default:
			return nil, err
		}
	}

	if len(defaultNsVols) > 0 {
		// Get the default ns id and put in the map
		defaultNs, err := c.client.GetNamespaceByName(ctx, "default")
		if err != nil {
			return nil, err
		}

		nsVols[defaultNs.ID] = append(nsVols[defaultNs.ID], defaultNsVols...)
	}

	resources := []*volume.Resource{}

	for nsID, volIDs := range nsVols {
		nsResources, err := c.client.GetNamespaceVolumes(ctx, nsID, volIDs...)
		if err != nil {
			return nil, err
		}

		resources = append(resources, nsResources...)
	}

	return resources, nil
}

func newVolume(w io.Writer, initClient func() (*apiclient.Client, error), config ConfigProvider) *cobra.Command {
	c := &volumeCommand{
		config: config,
		display: jsonformat.NewDisplayer(
			jsonformat.DefaultEncodingIndent,
		),

		writer: w,
	}

	cobraCommand := &cobra.Command{
		Aliases: []string{"volumes"},
		Use:     "volume [volume names...]",
		Short:   "volume retrieves basic information about StorageOS volumes",
		Example: `
$ storageos get volume fruits/banana
`,

		PreRunE: func(_ *cobra.Command, _ []string) error {
			client, err := initClient()
			if err != nil {
				return fmt.Errorf("error initialising api client: %w", err)
			}
			c.client = client
			return nil
		},
		RunE: c.run,

		// If a legitimate error occurs as part of the VERB volume command
		// we don't need to barf the usage template.
		SilenceUsage: true,
	}

	cobraCommand.Flags().BoolVar(&c.usingIDs, "use-id", false, "request StorageOS volumes by their namespace ID and volume ID instead of by their namespace name and volume name")

	return cobraCommand
}
