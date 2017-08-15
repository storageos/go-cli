package cluster

import (
	"net/url"

	"github.com/dnephin/cobra"

	apiTypes "github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/command/formatter"
	"github.com/storageos/go-cli/discovery"
	cliTypes "github.com/storageos/go-cli/types"
)

type healthOpt struct {
	cluster  string
	quiet    bool
	format   string
	selector string
}

func newHealthCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	opt := &healthOpt{}

	cmd := &cobra.Command{
		Use:   "health [CLUSTER_ID]",
		Short: `Displays the cluster's health.  When a cluster id is provided, uses the discovery service to discover nodes.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				opt.cluster = args[0]
			}
			return runHealth(storageosCli, opt)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opt.quiet, "quiet", "q", false, "Display minimal cluster health info.  Can be used with format.")
	flags.StringVar(&opt.format, "format", "", "Pretty-print health with formats: table (default), cp, dp or raw.")
	// flags.StringVarP(&opt.selector, "selector", "s", "", "Provide selector (e.g. to list all nodes with label region=eu ' --selector=region=eu')")

	return cmd
}

func runHealth(storageosCli *command.StorageOSCli, opt *healthOpt) error {

	nodes, err := getNodes(storageosCli, opt)
	if err != nil {
		return err
	}

	format := opt.format
	if len(format) == 0 {
		if len(storageosCli.ConfigFile().PoolsFormat) > 0 && !opt.quiet {
			format = storageosCli.ConfigFile().ClusterHealthFormat
		} else {
			format = formatter.TableFormatKey
		}
	}

	for _, node := range nodes {
		u, err := url.Parse(node.AdvertiseAddress)
		if err != nil {
			return err
		}

		cpHealth, err := storageosCli.Client().CPHealth(u.Hostname())
		if err != nil {
			// Don't exit if we can't collect health for a node
			continue
		}
		node.Health.CP = cpHealth

		dpHealth, err := storageosCli.Client().DPHealth(u.Hostname())
		if err != nil {
			// Don't exit if we can't collect health for a node
			continue
		}
		node.Health.DP = dpHealth
	}

	clusterHealthCtx := formatter.Context{
		Output: storageosCli.Out(),
		Format: formatter.NewClusterHealthFormat(format, opt.quiet),
	}
	return formatter.ClusterHealthWrite(clusterHealthCtx, nodes)
}

func getNodes(storageosCli *command.StorageOSCli, opt *healthOpt) ([]*cliTypes.Node, error) {

	if opt.cluster != "" {
		return getDiscoveryNodes(opt.cluster)
	}
	return getAPINodes(storageosCli)
}

func getDiscoveryNodes(clusterID string) ([]*cliTypes.Node, error) {

	client, err := discovery.NewClient("", "", "")
	if err != nil {
		return nil, err
	}

	cluster, err := client.ClusterStatus(clusterID)
	if err != nil {
		return nil, err
	}

	return cluster.Nodes, nil

}

func getAPINodes(storageosCli *command.StorageOSCli) ([]*cliTypes.Node, error) {

	apiNodes, err := storageosCli.Client().ControllerList(apiTypes.ListOptions{})
	if err != nil {
		return nil, err
	}

	var nodes []*cliTypes.Node
	for _, n := range apiNodes {
		node := &cliTypes.Node{
			ID:               n.ID,
			Name:             n.Name,
			AdvertiseAddress: "http://" + n.Address,
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}
