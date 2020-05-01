package cmd

import (
	"bytes"
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"code.storageos.net/storageos/c2-cli/config/file"
)

type configFilePathProvider interface {
	ConfigFilePath() (string, error)
}

func newConfigFileHelpTopic(configPathProvider configFilePathProvider) *cobra.Command {
	w := &bytes.Buffer{}

	var locationString string
	if path, err := configPathProvider.ConfigFilePath(); err == nil {
		locationString = fmt.Sprintf("The config file is currently sourced from %q.\n", path)
	}

	fmt.Fprintf(w, "The StorageOS CLI allows the user to provide their own defaults for some configuration settings through a YAML configuration file.\n%v\nBelow is an example generated config file:\n\n", locationString)

	// This is covered by a unit test and will not pass CI if the example
	// config file fails to encode.
	_ = yaml.NewEncoder(w).Encode(file.ExampleConfigFile)

	fmt.Fprintln(w)

	return &cobra.Command{
		Use:   "config-file",
		Short: "View help information for using a configuration file",
		Long:  w.String(),
	}
}
