package configfile

import (
	"runtime"
)

type credStore map[string]credentials

func (c credStore) HasCredetials(host string) bool {
	_, ok := c[host]
	return ok
}

func (c credStore) GetCredetials(host string) (username string, password string, err error) {
	creds, ok := c[host]
	if !ok {
		return "", "", ErrUnknownHost
	}

	username = creds.Username
	if creds.UseKeychain {
		password, err = creds.passwordFromKeychain(host)
	} else {
		password = string(creds.Password)
	}

	return
}

func (c credStore) SetCredetials(host string, username string, password string) error {
	creds := credentials{
		Username:    username,
		Password:    encodedPassword(password),
		UseKeychain: runtime.GOOS == "darwin",
	}

	c[host] = creds
	if creds.UseKeychain {
		return creds.saveToKeychain(host)
	}
	return nil
}

func (c credStore) DeleteCredetials(host string) {
	delete(c, host)
}
