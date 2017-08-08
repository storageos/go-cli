// +build darwin

package secret

import (
	"github.com/tmc/keyring"
)

const secretPrefix = "storageoscli_"

func secretGet(service string, username string) (string, error) {
	service = secretPrefix + service

	rtn, err := keyring.Get(service, username)
	if err == keyring.ErrNotFound {
		err = ErrNotFound
	}

	return rtn, err
}

func secretSet(service string, username string, secret string) error {
	service = secretPrefix + service

	return keyring.Set(service, username, secret)
}
