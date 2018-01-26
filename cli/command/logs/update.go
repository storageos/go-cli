package logs

import (
	"context"
	"fmt"
	"time"

	"github.com/storageos/go-api/types"
	"github.com/storageos/go-cli/cli/command"
)

func runUpdate(storageosCli *command.StorageOSCli, opt logOptions) error {

	client := storageosCli.Client()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(opt.timeout))
	defer cancel()

	config := types.LoggerUpdateOptions{
		Context: ctx,
	}

	if opt.level != "" {
		config.Level = opt.level
		config.Fields = append(config.Fields, "level")
	}

	// Filter can be set to empty
	if opt.filter != "" || opt.clearFilter {
		config.Filter = opt.filter
		config.Fields = append(config.Fields, "filter")
	}

	// Set nodes to update, if any specified (empty list means all)
	config.Nodes = opt.nodes

	if _, err := client.LoggerUpdate(config); err != nil {
		return fmt.Errorf("Failed to update logging: %v", err)
	}

	fmt.Fprintln(storageosCli.Out(), "OK")
	return nil
}
