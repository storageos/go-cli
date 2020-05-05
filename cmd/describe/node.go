package describe

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/cmd/flagutil"
	"code.storageos.net/storageos/c2-cli/cmd/runwrappers"
	"code.storageos.net/storageos/c2-cli/namespace"
	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/output"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/selectors"
	"code.storageos.net/storageos/c2-cli/volume"
)

type nodeCommand struct {
	config  ConfigProvider
	client  Client
	display Displayer

	selectors []string

	writer io.Writer
}

func (c *nodeCommand) runWithCtx(ctx context.Context, cmd *cobra.Command, args []string) error {
	switch len(args) {
	case 1:
		n, hostedVols, err := c.describeNode(ctx, args[0])
		if err != nil {
			return err
		}

		namespaceForID, err := c.getNamespaceForID(ctx)
		if err != nil {
			return err
		}

		return c.display.DescribeNode(
			ctx,
			c.writer,
			output.NewNodeDescription(n, hostedVols, namespaceForID),
		)
	default:
		set, err := selectors.NewSetFromStrings(c.selectors...)
		if err != nil {
			return err
		}

		nodes, hostedVolsMap, err := c.describeListNodes(ctx, args)
		if err != nil {
			return err
		}

		namespaceForID, err := c.getNamespaceForID(ctx)
		if err != nil {
			return err
		}

		nodeDescriptions := []*output.NodeDescription{}

		for _, n := range set.FilterNodes(nodes) {
			nodeDescriptions = append(
				nodeDescriptions,
				output.NewNodeDescription(n, hostedVolsMap[n.ID], namespaceForID),
			)
		}

		return c.display.DescribeListNodes(ctx, c.writer, nodeDescriptions)
	}
}

// describeNode retrieves a single node resource and a list of volume resources
// which it hosts a deployment for using the API client.
//
// The node is retrieved by name or ID based on the current command
// configuration.
func (c *nodeCommand) describeNode(ctx context.Context, ref string) (*node.Resource, []*volume.Resource, error) {
	useIDs, err := c.config.UseIDs()
	if err != nil {
		return nil, nil, err
	}

	var n *node.Resource

	if !useIDs {
		n, err = c.client.GetNodeByName(ctx, ref)
	} else {
		uid := id.Node(ref)
		n, err = c.client.GetNode(ctx, uid)
	}
	if err != nil {
		return nil, nil, err
	}

	allVols, err := c.client.GetAllVolumes(ctx)
	if err != nil {
		return nil, nil, err
	}

	var hosted []*volume.Resource

	for _, v := range allVols {
		if v.Master.Node == n.ID {
			hosted = append(hosted, v)
			continue
		}

		for _, r := range v.Replicas {
			if r.Node == n.ID {
				hosted = append(hosted, v)
				break
			}
		}
	}

	return n, hosted, nil
}

// describeListNodes retrieves a list of node resources and a mapping from
// nodes to volumes which they host a deployment for using the API client.
//
// The nodes are retrieved by name or ID based on the current command
// configuration.
func (c *nodeCommand) describeListNodes(ctx context.Context, refs []string) ([]*node.Resource, map[id.Node][]*volume.Resource, error) {
	useIDs, err := c.config.UseIDs()
	if err != nil {
		return nil, nil, err
	}

	var nodeList []*node.Resource

	if !useIDs {
		nodeList, err = c.client.GetListNodesByName(ctx, refs...)
	} else {
		uids := make([]id.Node, len(refs))
		for i, ref := range refs {
			uids[i] = id.Node(ref)
		}
		nodeList, err = c.client.GetListNodesByUID(ctx, uids...)
	}
	if err != nil {
		return nil, nil, err
	}

	allVols, err := c.client.GetAllVolumes(ctx)
	if err != nil {
		return nil, nil, err
	}

	nodeToHosted := make(map[id.Node][]*volume.Resource)
	for _, v := range allVols {
		nodeToHosted[v.Master.Node] = append(nodeToHosted[v.Master.Node], v)

		for _, r := range v.Replicas {
			nodeToHosted[r.Node] = append(nodeToHosted[r.Node], v)
		}
	}

	return nodeList, nodeToHosted, nil
}

func (c *nodeCommand) getNamespaceForID(ctx context.Context) (map[id.Namespace]*namespace.Resource, error) {
	namespaces, err := c.client.ListNamespaces(ctx)
	if err != nil {
		return nil, err
	}

	namespaceForID := make(map[id.Namespace]*namespace.Resource)
	for _, ns := range namespaces {
		namespaceForID[ns.ID] = ns
	}

	return namespaceForID, nil
}

func newNode(w io.Writer, client Client, config ConfigProvider) *cobra.Command {
	c := &nodeCommand{
		config: config,
		client: client,
		writer: w,
	}
	cobraCommand := &cobra.Command{
		Aliases: []string{"nodes"},
		Use:     "node [node names...]",
		Short:   "Retrieve detailed information for nodes in the cluster",
		Example: `
$ storageos describe node my-node-name
`,

		PreRunE: func(_ *cobra.Command, _ []string) error {
			c.display = SelectDisplayer(c.config)

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			run := runwrappers.Chain(
				runwrappers.RunWithTimeout(c.config),
				runwrappers.EnsureTargetOrSelectors(&c.selectors),
				runwrappers.AuthenticateClient(c.config, c.client),
			)(c.runWithCtx)
			return run(context.Background(), cmd, args)
		},

		SilenceUsage: true,
	}

	flagutil.SupportSelectors(cobraCommand.Flags(), &c.selectors)

	return cobraCommand
}
