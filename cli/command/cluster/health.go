package cluster

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"time"

	"github.com/dnephin/cobra"

	apiTypes "github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/command/formatter"
	cliNode "github.com/storageos/go-cli/cli/command/node"
	"github.com/storageos/go-cli/discovery"
	"github.com/storageos/go-cli/pkg/constants"
	cliTypes "github.com/storageos/go-cli/types"
)

type healthOpt struct {
	cluster string
	quiet   bool
	format  string
	timeout int
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
	flags.IntVarP(&opt.timeout, "timeout", "t", constants.DefaultAPITimeout, "Timeout in seconds.")
	flags.StringVar(&opt.format, "format", "", "Pretty-print health with formats: table (default), detailed, cp, dp or raw.")

	return cmd
}

func runHealth(storageosCli *command.StorageOSCli, opt *healthOpt) error {

	nodes, err := getNodes(storageosCli, opt)
	if err != nil {
		return err
	}

	format := opt.format
	if len(format) == 0 {
		if len(storageosCli.ConfigFile().ClusterHealthFormat) > 0 && !opt.quiet {
			format = storageosCli.ConfigFile().ClusterHealthFormat
		} else {
			format = formatter.TableFormatKey
		}
	}

	sort.Sort(cliTypes.NodeByName(nodes))
	for _, node := range nodes {
		if err := runNodeHealth(node, opt.timeout); err != nil {
			return err
		}
	}

	clusterHealthCtx := formatter.Context{
		Output: storageosCli.Out(),
		Format: formatter.NewClusterHealthFormat(format, opt.quiet),
	}
	return formatter.ClusterHealthWrite(clusterHealthCtx, nodes)
}

func runNodeHealth(node *cliTypes.Node, timeout int) error {
	addr := node.AdvertiseAddress

	u, err := url.Parse(node.AdvertiseAddress)
	if err == nil && u.Host != "" {
		addr = u.Host
	}

	cliNode.UpdateNodeHealth(node, addr, timeout)

	return nil
}

func getNodes(storageosCli *command.StorageOSCli, opt *healthOpt) ([]*cliTypes.Node, error) {

	if opt.cluster != "" {
		return getDiscoveryNodes(storageosCli.GetDiscovery(), opt.cluster)
	}
	return getAPINodes(storageosCli, opt.timeout)
}

func getDiscoveryNodes(discoveryHost, clusterID string) ([]*cliTypes.Node, error) {

	client, err := discovery.NewClient(discoveryHost, "", "")
	if err != nil {
		return nil, err
	}

	cluster, err := client.ClusterStatus(clusterID)
	if err != nil {
		return nil, err
	}

	return cluster.Nodes, nil

}

func getAPINodes(storageosCli *command.StorageOSCli, timeout int) ([]*cliTypes.Node, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()

	listOptions := apiTypes.ListOptions{
		Context: ctx,
	}
	apiNodes, err := storageosCli.Client().NodeList(listOptions)
	if err != nil {
		return nil, fmt.Errorf("API not responding to list nodes: %v", err)
	}

	var nodes []*cliTypes.Node
	for _, n := range apiNodes {
		node := &cliTypes.Node{
			ID:               n.ID,
			Name:             n.Name,
			AdvertiseAddress: n.Address,
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}
