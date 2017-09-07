package configfile

import (
	"encoding/base64"
	"encoding/json"
	"os/exec"
	"regexp"
	"runtime"
	"syscall"
)

type credentials struct {
	Username    string           `json:"username"`
	Password    *encodedPassword `json:"password,omitempty"`
	UseKeychain bool             `json:"useKeychain,omitempty"`
}

func (c credentials) saveToKeychain(host string) error {
	if runtime.GOOS != "darwin" {
		return ErrNotDarwin
	}

	return exec.Command(
		"/usr/bin/security", "add-generic-password",
		"-s", "storageos_cli",
		"-a", host,
		"-w", string(*c.Password),
		"-U",
	).Run()
}

func (c credentials) passwordFromKeychain(host string) (string, error) {
	if runtime.GOOS != "darwin" {
		return "", ErrNotDarwin
	}

	com := exec.Command(
		"/usr/bin/security", "find-generic-password",
		"-s", "storageos_cli",
		"-a", host,
		"-g",
	)

	out, err := com.CombinedOutput()
	exitCode := com.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()

	switch {
	case err != nil && exitCode == 44:
		return "", ErrNotFound

	case err != nil:
		return "", err

	default:
		matches := regexp.MustCompile("password: \"(.+)\"").FindStringSubmatch(string(out))
		if len(matches) != 2 {
			return "", ErrNotFound
		}

		return matches[1], nil
	}

}

type encodedPassword string

func (e *encodedPassword) MarshalJSON() ([]byte, error) {
	return json.Marshal(base64.StdEncoding.EncodeToString([]byte(*e)))
}

func (e *encodedPassword) UnmarshalJSON(data []byte) error {
	var encoded string

	err := json.Unmarshal(data, &encoded)
	if err != nil {
		return err
	}

	bytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return err
	}

	*e = encodedPassword(bytes)
	return nil
}
