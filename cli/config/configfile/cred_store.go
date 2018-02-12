package configfile

import (
	"fmt"
	"runtime"
	"strings"
)

type CredStore map[string]credentials

func (c CredStore) HasCredentials(host string) bool {
	_, ok := c[host]
	return ok
}

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

func (c CredStore) DeleteCredentials(host string) {
	for _, h := range strings.Split(host, ",") {
		delete(c, h)
	}
}
