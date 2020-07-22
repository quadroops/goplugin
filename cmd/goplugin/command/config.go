package command

import (
	"github.com/quadroops/goplugin/cmd/goplugin/command/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
}

func init() {
	configCmd.AddCommand(config.ConfigGenerateCmd)
}
