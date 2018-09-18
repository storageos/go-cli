package configfile

import (
	"fmt"
	"runtime"
	"strings"
)

// A CredStore maps host addresses to credentials which can be used to
// authenticate users with them. A credentials store is usually backed
// by a config file, although passwords are stored in supported keychain
// applications where available.
type CredStore map[string]credentials

// HasCredentials returns the existence of a host address in the credentials
// store.
func (c CredStore) HasCredentials(host string) bool {
	_, ok := c[host]
	return ok
}

// GetCredentials returns the credentials stored for the given host address,
// if present. If keychain usage is enabled, the password will be retrieved
// from there.
func (c CredStore) GetCredentials(host string) (username string, password string, err error) {
	creds, ok := c[host]
	if !ok {
		return "", "", ErrUnknownHost
	}

	username = creds.Username
	if creds.UseKeychain {
		password, err = creds.passwordFromKeychain(host)
	} else {
		password = string(*creds.Password)
	}

	return
}

// SetCredentials writes the given credentials to the credentials stores entry
// for the specified host, replacing any prior credentials. The password is
// encoded before being stored, whilst a keychain application is used to store
// it where a supported application is present.
func (c CredStore) SetCredentials(host string, username string, password string) error {
	pass := encodedPassword(password)

	creds := credentials{
		Username:    username,
		Password:    &pass,
		UseKeychain: runtime.GOOS == "darwin",
	}

	var errs []error
	for _, h := range strings.Split(host, ",") {
		c[h] = creds
		if creds.UseKeychain {
			if err := creds.saveToKeychain(h); err != nil {
				errs = append(errs, err)
			}
		}
	}

	if errs != nil {
		return fmt.Errorf("error: %+v", errs)
	}
	return nil
}

// DeleteCredentials will remove the entry corresponding to the provided host
// from the credentials store.
func (c CredStore) DeleteCredentials(host string) {
	for _, h := range strings.Split(host, ",") {
		delete(c, h)
	}
}
