package cluster

import (
	"fmt"
	"net/url"

	"github.com/dnephin/cobra"

	apiTypes "github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/command/formatter"
	"github.com/storageos/go-cli/discovery"
	cliTypes "github.com/storageos/go-cli/types"
)

type healthOpt struct {
	cpHealth bool
	dpHealth bool
	cluster  string
}

func (h healthOpt) cp() bool {
	return h.cpHealth || !(h.cpHealth || h.dpHealth)
}

func (h healthOpt) dp() bool {
	return h.dpHealth || !(h.dpHealth || h.cpHealth)
}

func newHealthCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	opt := &healthOpt{}

	cmd := &cobra.Command{
		Use:   "health [--cp | --dp] [CLUSTER_ID]",
		Short: `Displays the cluster's health.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				opt.cluster = args[0]
			}
			return runHealth(storageosCli, opt)
		},
	}

	flag := cmd.Flags()
	flag.BoolVar(&opt.cpHealth, "cp", false, "Display the output from the control plane only")
	flag.BoolVar(&opt.dpHealth, "dp", false, "Display the output from the data plane only")

	return cmd
}

func runCPHealth(storageosCli *command.StorageOSCli, nodes []*cliTypes.Node) error {
	clusterHealth := &apiTypes.ClusterHealthCP{}

	for _, node := range nodes {
		fmt.Printf("node: %#v\n", node)
		u, err := url.Parse(node.AdvertiseAddress)
		if err != nil {
			return err
		}

		status, err := storageosCli.Client().CPHealth(u.Hostname())
		if err != nil {
			return err
		}

		clusterHealth.Add(node.ID, status)
	}

	return formatter.ClusterHealthCPWrite(formatter.Context{
		Output: storageosCli.Out(),
		Format: formatter.NewHealthCPFormat(formatter.TableFormatKey),
	}, clusterHealth)
}

func runDPHealth(storageosCli *command.StorageOSCli, nodes []*cliTypes.Node) error {
	clusterHealth := &apiTypes.ClusterHealthDP{}

	for _, node := range nodes {
		u, err := url.Parse(node.AdvertiseAddress)
		if err != nil {
			return err
		}

		status, err := storageosCli.Client().DPHealth(u.Hostname())
		if err != nil {
			return err
		}

		clusterHealth.Add(node.ID, status)
	}

	return formatter.ClusterHealthDPWrite(formatter.Context{
		Output: storageosCli.Out(),
		Format: formatter.NewHealthDPFormat(formatter.TableFormatKey),
	}, clusterHealth)
}

func runHealth(storageosCli *command.StorageOSCli, opt *healthOpt) error {

	nodes, err := getNodes(storageosCli, opt)
	if err != nil {
		return err
	}

	switch {
	case opt.cp() && opt.dp():
		fmt.Fprintln(storageosCli.Out(), "Controlplane:")
		if err := runCPHealth(storageosCli, nodes); err != nil {
			return err
		}

		fmt.Fprintln(storageosCli.Out(), "\nDataplane:")
		return runDPHealth(storageosCli, nodes)

	case opt.cp():
		return runCPHealth(storageosCli, nodes)

	default:
		return runDPHealth(storageosCli, nodes)
	}
	return nil
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
