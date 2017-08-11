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
	cpHealth  bool
	dpHealth  bool
	clusterID string
	nodeID    string
}

func (h healthOptions) cp() bool {
	return h.cpHealth || !(h.cpHealth || h.dpHealth)
}

func (h healthOptions) dp() bool {
	return h.dpHealth || !(h.dpHealth || h.cpHealth)
}

func newHealthCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	opt := &healthOptions{}

	cmd := &cobra.Command{
		Use:   "health cp|dp NODE_ID ",
		Short: "Display detailed information on a given node",
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opt.nodeID = args[0]

			if opt.clusterID != "" {
				return runHealthFromClusterID(storageosCli, opt)
			}

			return runHealthFromENV(storageosCli, opt)
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&opt.cpHealth, "cp", false, "Display the output from the control plane only")
	flags.BoolVar(&opt.dpHealth, "dp", false, "Display the output from the data plane only")
	flags.StringVar(&opt.clusterID, "cluster", "", "Find the node's IP address from a cluster token")

	return cmd
}

func runHealthFromAddr(storageosCli *command.StorageOSCli, addr string, opt *healthOptions) error {
	u, err := url.Parse(addr)
	if err != nil {
		return err
	}

	switch {
	case opt.cp() && opt.dp():
		cphealth, err := storageosCli.Client().CPHealth(u.Hostname())
		if err != nil {
			return err
		}

		fmt.Fprintln(storageosCli.Out(), "Controlplane:")
		if err := formatter.NodeHealthWrite(formatter.Context{
			Output: storageosCli.Out(),
			Format: formatter.NewNodeHealthFormat(formatter.TableFormatKey),
		}, cphealth.ToNamedSubmodules()); err != nil {
			return err
		}

		fmt.Fprintln(storageosCli.Out(), "\nDataplane:")
		dphealth, err := storageosCli.Client().DPHealth(u.Hostname())
		if err != nil {
			return err
		}

		return formatter.NodeHealthWrite(formatter.Context{
			Output: storageosCli.Out(),
			Format: formatter.NewNodeHealthFormat(formatter.TableFormatKey),
		}, dphealth.ToNamedSubmodules())

	case opt.cp():
		health, err := storageosCli.Client().CPHealth(u.Hostname())
		if err != nil {
			return err
		}

		return formatter.NodeHealthWrite(formatter.Context{
			Output: storageosCli.Out(),
			Format: formatter.NewNodeHealthFormat(formatter.TableFormatKey),
		}, health.ToNamedSubmodules())

	default:
		health, err := storageosCli.Client().DPHealth(u.Hostname())
		if err != nil {
			return err
		}

		return formatter.NodeHealthWrite(formatter.Context{
			Output: storageosCli.Out(),
			Format: formatter.NewNodeHealthFormat(formatter.TableFormatKey),
		}, health.ToNamedSubmodules())
	}
}

func runHealthFromClusterID(storageosCli *command.StorageOSCli, opt *healthOptions) error {
	client, err := discovery.NewClient("", "", "")
	if err != nil {
		return err
	}

	cluster, err := client.ClusterStatus(opt.clusterID)
	if err != nil {
		return err
	}

	for _, node := range cluster.Nodes {
		if node.ID == opt.nodeID {
			return runHealthFromAddr(storageosCli, node.AdvertiseAddress, opt)
		}
	}

	return fmt.Errorf("Failed to find node (%s) in cluster (%s)", opt.nodeID, opt.clusterID)
}

func runHealthFromENV(storageosCli *command.StorageOSCli, opt *healthOptions) error {
	node, err := storageosCli.Client().Controller(opt.nodeID)
	if err != nil {
		return err
	}

	// runHealthFromAddr runs url.Parse on the given url, this means that the url
	// must have a scheme to be valid. The only field used from the url is
	// hostname and the scheme is ignored. For this reason the scheme "scheme"
	// has been used here, to make it more obvious that this scheme doesn't
	// change behaviour at all.
	return runHealthFromAddr(storageosCli, "scheme://"+node.Address, opt)
}
