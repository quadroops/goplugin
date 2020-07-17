package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd is our main cmd provider
var rootCmd = &cobra.Command{
	Use:   "goplugin",
	Short: "Goplugin is golang plugin library to manage plugin architecture",
}

func init() {
	rootCmd.AddCommand(discoverCmd)
}

// Execute will provide main cmd application
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(fmt.Sprintf("Error: %v", err))
		os.Exit(1)
	}
}