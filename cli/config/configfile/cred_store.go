package configfile

import (
	"runtime"
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

	c[host] = creds
	if creds.UseKeychain {
		return creds.saveToKeychain(host)
	}
	return nil
}

func (c CredStore) DeleteCredentials(host string) {
	delete(c, host)
}
