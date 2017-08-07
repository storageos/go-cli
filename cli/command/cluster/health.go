package cluster

import (
	"fmt"
	"net/url"

	"github.com/dnephin/cobra"

	apiTypes "github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/command/formatter"
	"github.com/storageos/go-cli/discovery"
	cliTypes "github.com/storageos/go-cli/types"
)

func newHealthCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "health cp|dp CLUSTER_TOKEN",
		Short: `Displays the cluster's health information from a cluster token (as given by cluster create)`,
		Args:  cli.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runHealth(storageosCli, args[1], args[0])
		},
	}

	return cmd
}

func runCPHealth(storageosCli *command.StorageOSCli, cluster *cliTypes.Cluster) error {
	clusterHealth := &apiTypes.ClusterHealthCP{}

	for _, node := range cluster.Nodes {
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

func runDPHealth(storageosCli *command.StorageOSCli, cluster *cliTypes.Cluster) error {
	clusterHealth := &apiTypes.ClusterHealthDP{}

	for _, node := range cluster.Nodes {
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

func runHealth(storageosCli *command.StorageOSCli, clusterID string, cpdp string) error {
	client, err := discovery.NewClient("", "", "")
	if err != nil {
		return err
	}

	cluster, err := client.ClusterStatus(clusterID)
	if err != nil {
		return err
	}

	if cpdp == "cp" {
		return runCPHealth(storageosCli, cluster)
	}
	if cpdp == "dp" {
		return runDPHealth(storageosCli, cluster)
	}

	return fmt.Errorf("Unknown instance type selector: %s (expecting [cp|dp])", cpdp)
}
