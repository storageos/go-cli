package cmd

import (
	"bytes"
	"fmt"

	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"

	"code.storageos.net/storageos/c2-cli/config/environment"
)

func newEnvConfigHelpTopic() *cobra.Command {
	w := &bytes.Buffer{}

	fmt.Fprintf(w, "The StorageOS CLI allows the user to provide their own defaults for some configuration settings through environment variables.\n\nAvailable Settings:\n")

	table := uitable.New()
	table.MaxColWidth = 80
	table.Separator = "  "
	table.Wrap = true

	for _, doc := range environment.EnvConfigHelp {
		table.AddRow(fmt.Sprintf("  %s", doc.Name), doc.Help)
	}
	fmt.Fprintln(w, table)

	return &cobra.Command{
		Use:   "env",
		Short: "View documentation for configuration settings which can be set in the environment",
		Long:  w.String(),
	}
}
