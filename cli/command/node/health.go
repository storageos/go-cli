package node

import (
	"fmt"
	"github.com/dnephin/cobra"
	"net/url"
	// storageos "github.com/storageos/go-api"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/command/formatter"
	"github.com/storageos/go-cli/discovery"
)

type healthOptions struct {
	clusterID string
}

func newHealthCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	opt := &healthOptions{}

	cmd := &cobra.Command{
		Use:   "health cp|dp NODE_ID ",
		Short: "Display detailed information on a given node",
		Args:  cli.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opt.clusterID != "" {
				return runHealthFromClusterID(storageosCli, args[0], args[1], opt.clusterID)
			}

			return runHealthFromENV(storageosCli, args[0], args[1])
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opt.clusterID, "cluster", "", "Find the node's IP address from a cluster token")

	return cmd
}

func runHealthFromAddr(storageosCli *command.StorageOSCli, addr string, cpdp string) error {
	u, err := url.Parse(addr)
	if err != nil {
		return err
	}

	if cpdp == "cp" {
		health, err := storageosCli.Client().CPHealth(u.Hostname())
		if err != nil {
			return err
		}

		return formatter.NodeHealthWrite(formatter.Context{
			Output: storageosCli.Out(),
			Format: formatter.NewNodeHealthFormat(formatter.TableFormatKey),
		}, health.ToNamedSubmodules())
	}

	if cpdp == "dp" {
		health, err := storageosCli.Client().DPHealth(u.Hostname())
		if err != nil {
			return err
		}

		return formatter.NodeHealthWrite(formatter.Context{
			Output: storageosCli.Out(),
			Format: formatter.NewNodeHealthFormat(formatter.TableFormatKey),
		}, health.ToNamedSubmodules())
	}

	return fmt.Errorf("Unknown instance type selector: %s (expecting [cp|dp])", cpdp)
}

func runHealthFromClusterID(storageosCli *command.StorageOSCli, cpdp string, nodeID string, clusterID string) error {
	client, err := discovery.NewClient("", "", "")
	if err != nil {
		return err
	}

	cluster, err := client.ClusterStatus(clusterID)
	if err != nil {
		return err
	}

	for _, node := range cluster.Nodes {
		if node.ID == nodeID {
			return runHealthFromAddr(storageosCli, node.AdvertiseAddress, cpdp)
		}
	}

	return fmt.Errorf("Failed to find node (%s) in cluster (%s)", nodeID, clusterID)
}

func runHealthFromENV(storageosCli *command.StorageOSCli, cpdp string, nodeID string) error {
	node, err := storageosCli.Client().Controller(nodeID)
	if err != nil {
		return err
	}

	return runHealthFromAddr(storageosCli, node.Address, cpdp)
}
