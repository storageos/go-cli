// +build !linux

package volume

import (
	"github.com/dnephin/cobra"
	"github.com/storageos/go-cli/cli/command"
)

// Mount commands only supported on linux
func newMountCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	return &cobra.Command{}
}
