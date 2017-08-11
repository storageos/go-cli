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
		Use:   "health [--cp | --dp] CLUSTER_TOKEN",
		Short: `Displays the cluster's health information from a cluster token (as given by cluster create)`,
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opt.cluster = args[0]
			return runHealth(storageosCli, opt)
		},
	}

	flag := cmd.Flags()
	flag.BoolVar(&opt.cpHealth, "cp", false, "Display the output from the control plane only")
	flag.BoolVar(&opt.dpHealth, "dp", false, "Display the output from the data plane only")

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

func runHealth(storageosCli *command.StorageOSCli, opt *healthOpt) error {
	client, err := discovery.NewClient("", "", "")
	if err != nil {
		return err
	}

	cluster, err := client.ClusterStatus(opt.cluster)
	if err != nil {
		return err
	}

	switch {
	case opt.cp() && opt.dp():
		fmt.Fprintln(storageosCli.Out(), "Controlplane:")
		if err := runCPHealth(storageosCli, cluster); err != nil {
			return err
		}

		fmt.Fprintln(storageosCli.Out(), "\nDataplane:")
		return runDPHealth(storageosCli, cluster)

	case opt.cp():
		return runCPHealth(storageosCli, cluster)

	default:
		return runDPHealth(storageosCli, cluster)
	}
}
