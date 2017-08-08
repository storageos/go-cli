package secret

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"os/user"
	"sync"
	"syscall"
)

const DefaultUserSecretName = "default_username"
const DefaultPassSecretName = "default_password"

var globalMu = &sync.Mutex{}

var ErrNotFound = errors.New("Unknown secret")

func GetSecret(secretName string) (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}

	globalMu.Lock()
	defer globalMu.Unlock()

	return secretGet(secretName, user.Username)
}

func SetSecret(secretName string, secret string) error {
	user, err := user.Current()
	if err != nil {
		return err
	}

	globalMu.Lock()
	defer globalMu.Unlock()

	return secretSet(secretName, user.Username, secret)
}

func PromptOrGetSecret(prompt string, secretName string, hideInput bool) (string, error) {
	secret, err := GetSecret(secretName)
	if err != nil && err != ErrNotFound {
		fmt.Println("Secret-file error: " + err.Error())
	}

	// If we got a secret return it
	if err == nil {
		return secret, nil
	}

	// Secret not found, get it
	secret = ""

	restorePoint, err := terminal.MakeRaw(syscall.Stdin)
	if err != nil {
		return "", err
	}
	defer terminal.Restore(syscall.Stdin, restorePoint)

	term := terminal.NewTerminal(os.Stdin, prompt+"> ")

	if hideInput {
		secret, err = term.ReadPassword(prompt + "> ")
		if err != nil {
			return "", err
		}
	} else {
		secret, err = term.ReadLine()
		if err != nil {
			return "", err
		}
	}

	if err := SetSecret(secretName, secret); err != nil {
		fmt.Println("Secret-file error: " + err.Error())
	}

	return secret, nil
}
