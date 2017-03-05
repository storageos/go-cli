package version

import (
	"runtime"

	"github.com/storageos/go-api/types"
)

const (
	ProductName string = "storageos"
	APIVersion         = "1"
)

// Revision that was compiled. This will be filled in by the compiler.
var Revision string

// BuildDate is when the binary was compiled.  This will be filled in by the
// compiler.
var BuildDate string

// Version number that is being run at the moment.  Version should use semver.
var Version string

// Experimental is intended to be used to enable alpha features.
var Experimental string

// GetStorageOSVersion returns version info.
func GetStorageOSVersion() types.VersionInfo {
	v := types.VersionInfo{
		Name:       ProductName,
		Revision:   Revision,
		BuildDate:  BuildDate,
		Version:    Version,
		APIVersion: APIVersion,
		GoVersion:  runtime.Version(),
		OS:         runtime.GOOS,
		Arch:       runtime.GOARCH,
	}

	// kernelVersion := "<unknown>"
	// if kv, err := kernel.GetKernelVersion(); err != nil {
	// 	log.Warnf("Could not get kernel version: %v", err)
	// } else {
	// 	kernelVersion = kv.String()
	// }
	// v.KernelVersion = kernelVersion

	return v
}
