// +build !linux

package volume

import (
	"fmt"
	"github.com/storageos/go-cli/cli/command"
	"runtime"
)

func runMount(storageosCli *command.StorageOSCli, opt mountOptions) error {
	return fmt.Errorf("volume mount not natively supported on %s", runtime.GOOS)
}
