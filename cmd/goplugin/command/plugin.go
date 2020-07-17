package command

import (
	"github.com/quadroops/goplugin/cmd/goplugin/command/plugin"
	"github.com/spf13/cobra"
)

var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Manage plugins",
}

func init() {
	pluginCmd.AddCommand(plugin.PluginRegisterCmd)
}
