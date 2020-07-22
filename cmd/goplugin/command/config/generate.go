package config

import (
	"fmt"
	"log"
	"os"

	"github.com/pelletier/go-toml"

	"github.com/quadroops/goplugin/pkg/discover"
	"github.com/quadroops/goplugin/pkg/discover/driver"

	"github.com/spf13/cobra"
)

const (
	// Version is current goplugin's version
	Version = "1.0.0"

	// Author current goplugin's author
	Author = "hiraq|hiraq.dev@gmail.com"

	// Debug is current goplugin's settings for debugging
	Debug = false
)

// ConfigGenerateCmd used to generate main skeleton config
var ConfigGenerateCmd = &cobra.Command{
	Use:     "generate",
	Short:   "Generate skeleton config",
	Example: "goplugin config generate --config-path=/home/my/.custom",
	Run: func(cmd *cobra.Command, args []string) {
		var confpath string

		confPathFromArg, err := cmd.Flags().GetString("config-path")
		if err != nil {
			log.Fatalf("Error parse argument config path: %v", err)
		}

		if confPathFromArg != "" {
			confpath = confPathFromArg
		} else {
			d := discover.NewConfigChecker(driver.NewOsChecker(), driver.NewDefaultChecker())
			dfile, err := d.Explore()
			if err != nil {
				log.Fatalf("Error exploring plugins: %v", err)
			}

			confpath = dfile
		}

		meta := discover.PluginMeta{
			Version: Version,
			Author:  Author,
		}

		settings := discover.PluginSettings{
			Debug: Debug,
		}

		config := discover.PluginConfig{
			Meta:     meta,
			Settings: settings,
		}

		err = generate(confpath, config)
		if err != nil {
			log.Fatalf("Error generate config: %v", err)
		}

		fmt.Println("Config has been generated")
	},
}

func init() {
	ConfigGenerateCmd.Flags().StringP("config-path", "", "", "Set custom config path")
}

func generate(confpath string, config discover.PluginConfig) error {
	b, err := toml.Marshal(config)
	if err != nil {
		return err
	}

	f, err := os.Create(confpath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(b)
	return err
}
