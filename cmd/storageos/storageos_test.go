package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/testutil/assert"
	"github.com/storageos/go-cli/cli/command"
	"github.com/storageos/go-cli/cli/debug"
)

func TestClientDebugEnabled(t *testing.T) {
	defer debug.Disable()

	cmd := newStorageOSCommand(&command.StorageOSCli{})
	cmd.Flags().Set("debug", "true")

	err := cmd.PersistentPreRunE(cmd, []string{})
	assert.NilError(t, err)
	assert.Equal(t, os.Getenv("DEBUG"), "1")
	assert.Equal(t, logrus.GetLevel(), logrus.DebugLevel)
}

func TestExitStatusForInvalidSubcommandWithHelpFlag(t *testing.T) {
	discard := ioutil.Discard
	cmd := newStorageOSCommand(command.NewStorageOSCli(os.Stdin, discard, discard))
	cmd.SetArgs([]string{"help", "invalid"})
	err := cmd.Execute()
	assert.Error(t, err, "unknown help topic: invalid")
}
