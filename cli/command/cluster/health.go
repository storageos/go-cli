package cluster

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/dnephin/cobra"
	log "github.com/sirupsen/logrus"

	apiTypes "github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/command/formatter"
	"github.com/storageos/go-cli/discovery"
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
	flags.IntVarP(&opt.timeout, "timeout", "t", 1, "Timeout in seconds.")
	flags.StringVar(&opt.format, "format", "", "Pretty-print health with formats: table (default), cp, dp or raw.")

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

	for _, node := range nodes {
		if err := runNodeHealth(storageosCli, node, opt.timeout); err != nil {
			return err
		}
	}

	clusterHealthCtx := formatter.Context{
		Output: storageosCli.Out(),
		Format: formatter.NewClusterHealthFormat(format, opt.quiet),
	}
	return formatter.ClusterHealthWrite(clusterHealthCtx, nodes)
}

func runNodeHealth(storageosCli *command.StorageOSCli, node *cliTypes.Node, timeout int) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()

	addr := node.AdvertiseAddress

	u, err := url.Parse(node.AdvertiseAddress)
	if err == nil && u.Host != "" {
		addr = u.Host
	}

	// Ignore errors and carry on
	cpHealth, err := storageosCli.Client().CPHealth(ctx, addr)
	if err != nil {
		log.Debugf("error updating cp health: %v", err)
	}
	node.Health.CP = cpHealth

	dpHealth, err := storageosCli.Client().DPHealth(ctx, addr)
	if err != nil {
		log.Debugf("error updating dp health: %v", err)
	}
	node.Health.DP = dpHealth

	return nil
}

func getNodes(storageosCli *command.StorageOSCli, opt *healthOpt) ([]*cliTypes.Node, error) {

	if opt.cluster != "" {
		return getDiscoveryNodes(opt.cluster)
	}
	return getAPINodes(storageosCli, opt.timeout)
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

func getAPINodes(storageosCli *command.StorageOSCli, timeout int) ([]*cliTypes.Node, error) {

	ctx, _ := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	listOptions := apiTypes.ListOptions{
		Context: ctx,
	}
	apiNodes, err := storageosCli.Client().ControllerList(listOptions)
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
