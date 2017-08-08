// +build !darwin

package secret

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"os/user"
)

const currentConfigVersion = 1

type storedSecrets struct {
	Version int               `json:"version"`
	Secrets map[string]string `json:"secrets"`
}

func secretGet(service string, username string) (string, error) {
	user, err := user.Lookup(username)
	if err != nil {
		return "", err
	}

	if user.HomeDir == "" {
		return "", errors.New("No home directory")
	}

	filepath := user.HomeDir + "/.storageoscli"
	if err := os.MkdirAll(filepath, 0700); err != nil {
		return "", err
	}

	// If the file doesn't exist, create it for next time and return
	filepath += "/config.json"
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		f, err := os.Create(filepath)
		if err != nil {
			return "", err
		}
		defer f.Close()

		return "", ErrNotFound
	}

	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	// There may be no json in the file yet
	if len(data) == 0 {
		return "", ErrNotFound
	}

	ss := &storedSecrets{}
	if err := json.Unmarshal(data, ss); err != nil {
		return "", err
	}

	s, ok := ss.Secrets[service]
	if !ok {
		return "", ErrNotFound
	}

	return s, nil

}

func secretSet(service string, username string, secret string) error {
	user, err := user.Lookup(username)
	if err != nil {
		return err
	}

	if user.HomeDir == "" {
		return errors.New("No home directory")
	}

	filepath := user.HomeDir + "/.storageoscli"
	if err := os.MkdirAll(filepath, 0700); err != nil {
		return err
	}

	// If the file doesn't exist, create it and write a new json blob to it
	filepath += "/config.json"
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		data, err := json.MarshalIndent(&storedSecrets{
			Version: currentConfigVersion,
			Secrets: map[string]string{
				service: secret,
			},
		}, "", "    ")

		if err != nil {
			return err
		}

		return ioutil.WriteFile(filepath, data, 0600)
	}

	// If the file does exist read it in first
	ss := &storedSecrets{}
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	if len(data) > 0 {
		if err := json.Unmarshal(data, ss); err != nil {
			return err
		}
	} else {
		ss.Version = currentConfigVersion
		ss.Secrets = make(map[string]string)
	}

	// add a new entry and then write it back
	ss.Secrets[service] = secret
	data, err = json.MarshalIndent(ss, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath, data, 0600)
}
