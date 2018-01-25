package netutil

import (
	"errors"
	"fmt"
	sosErr "github.com/storageos/go-cli/pkg/errors"
	"strings"
)

func errAllFailed(addrs []string) error {
	msg := fmt.Sprintf("failed to dial all known cluster members, (%s)", strings.Join(addrs, ","))
	help := "ensure that the value of $STORAGEOS_HOST (or the -H flag) is correct, and that there are healthy StorageOS nodes in this cluster"

	return sosErr.NewTypedStorageOSError(sosErr.APIUncontactable, nil, msg, help)
}

func newInvalidNodeError(err error) error {
	msg := fmt.Sprintf("invalid node format: %s", err)
	help := "please check the format of $STORAGEOS_HOST (or the -H flag) complies with the StorageOS JOIN format"

	return sosErr.NewTypedStorageOSError(sosErr.InvalidHostConfig, err, msg, help)
}

var errNoAddresses = errors.New("the MultiDialer instance has not been initialised with client addresses")
var errUnsupportedScheme = errors.New("unsupported URL scheme")
var errInvalidPortNumber = errors.New("invalid port number")
