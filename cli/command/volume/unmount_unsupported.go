// +build !linux

package volume

import (
	"fmt"
	"github.com/storageos/go-cli/cli/command"
	"runtime"
)

func runUnmount(storageosCli *command.StorageOSCli, opt unmountOptions) error {
	return fmt.Errorf("volume unmount not natively supported on %s", runtime.GOOS)
}
