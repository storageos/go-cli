package host

import (
	"os"
)

// Get - get current node hostname, used by `mount` package to determine client that's mounting
// a volume
func Get() (hostname string, err error) {
	return os.Hostname()
}
