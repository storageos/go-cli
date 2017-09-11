package configfile_test

import (
	"github.com/storageos/go-cli/cli/config"
	"io/ioutil"
	"runtime"
	"testing"
)

func initTestFile() (dir string, err error) {
	file := []byte(`{
	"knownHosts": {
		"localhost": {
			"username": "storageos",
			"password": "c3RvcmFnZW9z"
		},
		"otherhost": {
			"username": "notstorageos",
			"password": "bm90c3RvcmFnZW9z"
		}
	}
}`)

	d, err := ioutil.TempDir("", "")
	if err != nil {
		return "", err
	}

	err = ioutil.WriteFile(d+"/"+config.ConfigFileName, file, 0644)
	if err != nil {
		return "", err
	}

	return d, nil
}

func TestConfigReading(t *testing.T) {
	d, err := initTestFile()
	if err != nil {
		t.Fatal(err)
	}

	configFile, err := config.Load(d)
	if err != nil {
		t.Fatal(err)
	}

	store := configFile.CredentialsStore

	if !store.HasCredentials("localhost") {
		t.Fatal("Config file not correctly reading in knownHosts")
	}
	if !store.HasCredentials("otherhost") {
		t.Fatal("Config file not correctly reading in knownHosts")
	}
}

func TestConfigCredGet(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("Not testing cred get tests on darwin")
	}

	d, err := initTestFile()
	if err != nil {
		t.Fatal(err)
	}

	configFile, err := config.Load(d)
	if err != nil {
		t.Fatal(err)
	}

	user, pass, err := configFile.CredentialsStore.GetCredentials("localhost")
	if err != nil {
		t.Fatalf("Failed to get user and pass: %v", err)
	}

	if user != "storageos" {
		t.Fatalf("Got username (%v), expecting (storageos)", user)
	}
	if pass != "storageos" {
		t.Fatalf("Got password (%v), expecting (storageos)", pass)
	}

	user, pass, err = configFile.CredentialsStore.GetCredentials("otherhost")
	if err != nil {
		t.Fatalf("Failed to get user and pass: %v", err)
	}

	if user != "notstorageos" {
		t.Fatalf("Got username (%v), expecting (notstorageos)", user)
	}
	if pass != "notstorageos" {
		t.Fatalf("Got password (%v), expecting (notstorageos)", pass)
	}
}

func TestConfigCredStore(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("Not testing cred set tests on darwin")
	}

	d, err := initTestFile()
	if err != nil {
		t.Fatal(err)
	}

	configFile, err := config.Load(d)
	if err != nil {
		t.Fatal(err)
	}

	err = configFile.CredentialsStore.SetCredentials("foo", "bar", "baz")
	if err != nil {
		t.Fatal(err)
	}

	err = configFile.Save()
	if err != nil {
		t.Fatal(err)
	}

	reloaded, err := config.Load(d)
	if err != nil {
		t.Fatal(err)
	}

	user, pass, err := reloaded.CredentialsStore.GetCredentials("foo")
	if err != nil {
		t.Fatalf("Failed to get user and pass: %v", err)
	}

	if user != "bar" {
		t.Fatalf("Got username (%v), expecting (bar)", user)
	}
	if pass != "baz" {
		t.Fatalf("Got password (%v), expecting (baz)", pass)
	}
}
