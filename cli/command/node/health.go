package node

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dnephin/cobra"

	api "github.com/storageos/go-api"
	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/command/formatter"
	"github.com/storageos/go-cli/pkg/constants"
	cliTypes "github.com/storageos/go-cli/types"
)

type healthOptions struct {
	name    string
	quiet   bool
	format  string
	timeout int
}

func newHealthCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	opt := &healthOptions{}

	cmd := &cobra.Command{
		Use:   "health [OPTIONS] NODE",
		Short: "Display detailed information on a given node",
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opt.name = args[0]
			return runHealth(storageosCli, opt)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opt.quiet, "quiet", "q", false, "Display minimal node health info.  Can be used with format.")
	flags.StringVar(&opt.format, "format", "", "Pretty-print health with formats: table (default), cp, dp or raw.")
	flags.IntVarP(&opt.timeout, "timeout", "t", constants.DefaultAPITimeout, "Timeout in seconds.")

	return cmd
}

func runHealth(storageosCli *command.StorageOSCli, opt *healthOptions) error {

	c, err := storageosCli.Client().Node(opt.name)
	if err != nil {
		return err
	}

	node := &cliTypes.Node{
		ID:               c.ID,
		Name:             c.Name,
		AdvertiseAddress: c.Address,
	}

	UpdateNodeHealth(node, node.AdvertiseAddress, opt.timeout)

	format := opt.format
	if len(format) == 0 {
		if len(storageosCli.ConfigFile().NodeHealthFormat) > 0 && !opt.quiet {
			format = storageosCli.ConfigFile().NodeHealthFormat
		} else {
			format = formatter.TableFormatKey
		}
	}

	nodeHealthCtx := formatter.Context{
		Output: storageosCli.Out(),
		Format: formatter.NewNodeHealthFormat(format, opt.quiet),
	}
	return formatter.NodeHealthWrite(nodeHealthCtx, node)
}

// UpdateNodeHealth updates the health status of a given node by querying the
// node endpoints.
func UpdateNodeHealth(node *cliTypes.Node, address string, timeout int) error {
	healthEndpointFormat := "http://%s:%s/v1/" + api.HealthAPIPrefix

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()

	client := &http.Client{}

	var healthStatus types.HealthStatus
	cpURL := fmt.Sprintf(healthEndpointFormat, address, api.DefaultPort)
	cpReq, err := http.NewRequest("GET", cpURL, nil)
	if err != nil {
		return err
	}

	cpResp, err := client.Do(cpReq.WithContext(ctx))
	if err != nil {
		return err
	}

	if err := json.NewDecoder(cpResp.Body).Decode(&healthStatus); err != nil {
		return err
	}
	node.Health.CP = healthStatus.ToCPHealthStatus()
	node.Health.DP = healthStatus.ToDPHealthStatus()

	return nil
}
