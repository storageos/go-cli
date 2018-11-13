// +build !linux

package volume

import (
	"github.com/spf13/cobra"
	"github.com/storageos/go-cli/cli/command"
)

// Mount commands only supported on linux
func newMountCommand(storageosCli *command.StorageOSCli) *cobra.Command {
	return &cobra.Command{}
}
