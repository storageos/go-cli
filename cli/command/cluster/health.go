package cluster

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/dnephin/cobra"

	apiTypes "github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/command/formatter"
	"github.com/storageos/go-cli/pkg/constants"
	cliTypes "github.com/storageos/go-cli/types"
)

type healthOpt struct {
	errReturn bool
	quiet     bool
	format    string
	timeout   int
}

func newHealthCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	opt := &healthOpt{}

	cmd := &cobra.Command{
		Use:   "health",
		Short: `Displays the cluster's health.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runHealth(storageosCli, opt)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opt.quiet, "quiet", "q", false, "Display minimal cluster health info.  Can be used with format.")
	flags.BoolVarP(&opt.errReturn, "err", "e", false, "Return non-zero exit code if, and only if, cluster is fully operational")
	flags.IntVarP(&opt.timeout, "timeout", "t", constants.DefaultAPITimeout, "Timeout in seconds.")
	flags.StringVar(&opt.format, "format", "", "Pretty-print health with formats: table (default), detailed, cp, dp or raw.")

	return cmd
}

func runHealth(storageosCli *command.StorageOSCli, opt *healthOpt) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(opt.timeout))
	defer cancel()

	status, err := storageosCli.Client().ClusterHealth(ctx)
	if err != nil {
		return err
	}

	sort.Slice(status, func(i, j int) bool {
		return cliTypes.HumanisedStringLess(status[i].NodeName, status[j].NodeName)
	})

	format := opt.format
	if len(format) == 0 {
		if len(storageosCli.ConfigFile().ClusterHealthFormat) > 0 && !opt.quiet {
			format = storageosCli.ConfigFile().ClusterHealthFormat
		} else {
			format = formatter.TableFormatKey
		}
	}

	clusterHealthCtx := formatter.Context{
		Output: storageosCli.Out(),
		Format: formatter.NewClusterHealthFormat(format, opt.quiet),
	}
	if err := formatter.ClusterHealthWrite(clusterHealthCtx, status); err != nil {
		return err
	}
	if opt.errReturn {
		for _, s := range status {
			// Check for all the submodules that need to be alive for us to be
			// functioning
			if err := checkSubmodules(s); err != nil {
				return err
			}
		}
	}
	return nil
}

func checkSubmodules(nodeSubmodules *apiTypes.ClusterHealthNode) error {
	submoduleErrors := []string{}

	submoduleStates := map[string]apiTypes.SubModuleStatus{
		"DirectFSInitiator": nodeSubmodules.Submodules.DirectFSInitiator,
		"Director":          nodeSubmodules.Submodules.Director,
		"KV":                nodeSubmodules.Submodules.KV,
		"KVWrite":           nodeSubmodules.Submodules.KVWrite,
		"NATS":              nodeSubmodules.Submodules.NATS,
		"Presentation":      nodeSubmodules.Submodules.Presentation,
		"RDB":               nodeSubmodules.Submodules.RDB,
	}

	for name, state := range submoduleStates {
		if state.Status != "alive" {
			submoduleErrors = append(submoduleErrors, fmt.Sprintf("name: %q, message: '%q'", name, state.Message))
		}
	}

	if len(submoduleErrors) > 0 {
		return fmt.Errorf("\nnode: %q has unhealthy submodules:\n%q", nodeSubmodules.NodeName, strings.Join(submoduleErrors, "\n"))
	}
	return nil
}
