package plugin

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/pelletier/go-toml"

	"github.com/quadroops/goplugin/pkg/discover/driver"

	"github.com/quadroops/goplugin/pkg/discover"

	"github.com/spf13/cobra"
)

const (
	gopluginToml           = "goplugin.toml"
	gopluginTomlTemp       = "goplugin.temp.toml"
	gopluginTomlBackupTemp = "goplugin.backup.toml"
)

// Info used to store plugin's main data
type Info struct {
	Name          string
	Hosts         []string
	Author        string
	MD5           string
	Exec          string
	ExecArgs      []string `toml:"exec_args"`
	ExecFile      string   `toml:"exec_file"`
	ExecTime      int      `toml:"exec_time"`
	PrototcolType string   `toml:"protocol_type"`
}

// Plugin used as main data entity used to store all plugin's informations
// from goplugin.toml
type Plugin struct {
	Plugin Info
}

// PluginRegisterCmd will execute plugin registration process including for validation
// registration flows:
// - get current working directory
// - check if any goplugin.toml exist
// - parse the data from goplugin.toml
// - validate plugin:
//		- has valid author
//		- check if file/binary is exist
//		- compare MD5 sum
// - merge the data into global goplugin.toml in config directory or from given custom filepath
var PluginRegisterCmd = &cobra.Command{
	Use:     "register",
	Short:   "Register new plugin",
	Long:    "You can register a plugin from current working directory or by giving custom plugin path",
	Example: "goplugin plugin register --pluginpath=/home/my/myplugin",
	Run: func(cmd *cobra.Command, args []string) {
		var pluginPath, configpath string

		pluginPathFromArg, err := cmd.Flags().GetString("pluginpath")
		if err != nil {
			log.Fatalf("Error parse argument pluginpath: %v", err)
		}

		// first, we need to check if user giving a custom plugin file path
		if pluginPathFromArg != "" {
			pluginPath = pluginPathFromArg
		} else {
			// we move to second strategy, to trying to check current working directory
			// to make sure if any goplugin.toml
			currentDir, err := os.Getwd()
			if err != nil {
				log.Fatalf("Error get current working directory: %v", err)
			}

			pluginFromCurrentDir := fmt.Sprintf("%s/%s", currentDir, gopluginToml)
			_, err = os.Stat(pluginFromCurrentDir)
			if os.IsNotExist(err) {
				log.Fatalln("Plugin configuration file not found")
			}

			pluginPath = pluginFromCurrentDir
		}

		// we need to check if user giving a custom config path or not
		configPathFromArg, err := cmd.Flags().GetString("configpath")
		if err != nil {
			log.Fatalf("Error parsing configpath from arg: %v", err)
		}

		if configPathFromArg != "" {
			configpath = configPathFromArg
		} else {
			// if user not giving a custom config path, we need to move
			// to second strategy, to get the config path using our discover's package
			d := discover.NewConfigChecker(driver.NewOsChecker(), driver.NewDefaultChecker())
			dfile, err := d.Explore()
			if err != nil {
				log.Fatalf("Error exploring plugins: %v", err)
			}

			configpath = dfile
		}

		// parse plugin toml file
		plugin, err := parsePluginToml(pluginPath)
		if err != nil {
			log.Fatalf("Error parsing plugin toml: %v", err)
		}

		if plugin == nil {
			log.Fatalln("Plugin configuration is empty")
		}

		// parse config toml
		conf, err := parseConfToml(configpath)
		if err != nil {
			log.Fatalf("Error parsing config toml: %v", err)
		}

		confMerged := mergeConfig(conf, plugin)
		confTempPath, err := createConfigTemp(configpath, confMerged)
		if err != nil {
			log.Fatalf("Error create temporary config: %v", err)
		}

		err = os.Rename(configpath, rebuildConfPath(configpath, gopluginTomlBackupTemp))
		if err != nil {
			log.Fatalf("Error renaming original file: %v", err)
		}

		err = os.Rename(confTempPath, rebuildConfPath(confTempPath, gopluginToml))
		if err != nil {
			log.Fatalf("Error renaming temporary file: %v", err)

			// restore original config
			os.Rename(rebuildConfPath(configpath, gopluginTomlBackupTemp), configpath)
		}

		err = os.Remove(rebuildConfPath(configpath, gopluginTomlBackupTemp))
		if err != nil {
			log.Fatalf("Error deleting backup file: %v", err)

			// delete generated config
			os.Remove(configpath)

			// restore original config
			os.Rename(rebuildConfPath(configpath, gopluginTomlBackupTemp), configpath)
		}

		// finish , send the message to user
		fmt.Printf("Plugin: %s has been registered", plugin.Plugin.Name)
	},
}

func init() {
	PluginRegisterCmd.Flags().StringP("configpath", "c", "", "Set custom config filepath")
	PluginRegisterCmd.Flags().StringP("pluginpath", "p", "", "Set custom plugin filepath")
}

func parsePluginToml(filepath string) (*Plugin, error) {
	var plugin Plugin

	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	tomlBytes, err := toml.LoadBytes(content)
	if err != nil {
		return nil, err
	}

	err = tomlBytes.Unmarshal(&plugin)
	if err != nil {
		return nil, err
	}

	return &plugin, nil
}

func parseConfToml(confpath string) (*discover.PluginConfig, error) {
	var conf discover.PluginConfig

	content, err := ioutil.ReadFile(confpath)
	if err != nil {
		return nil, err
	}

	tomlBytes, err := toml.LoadBytes(content)
	if err != nil {
		return nil, err
	}

	err = tomlBytes.Unmarshal(&conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

func mergeConfig(confToml *discover.PluginConfig, plugin *Plugin) *discover.PluginConfig {
	pluginInfo := discover.PluginInfo{
		Author:       plugin.Plugin.Author,
		Exec:         plugin.Plugin.Exec,
		ExecArgs:     plugin.Plugin.ExecArgs,
		ExecFile:     plugin.Plugin.ExecFile,
		ExecTime:     plugin.Plugin.ExecTime,
		MD5:          plugin.Plugin.MD5,
		ProtocolType: plugin.Plugin.PrototcolType,
	}

	if len(plugin.Plugin.Hosts) >= 1 {
		for _, host := range plugin.Plugin.Hosts {
			if hostPlugins, exist := confToml.Hosts[host]; exist {
				hostPlugins.Plugins = append(hostPlugins.Plugins, plugin.Plugin.Name)
				confToml.Hosts[host] = hostPlugins
			} else {
				confToml.Hosts[host] = discover.PluginHost{
					Plugins: []string{plugin.Plugin.Name},
				}
			}
		}
	}

	confToml.Plugins[plugin.Plugin.Name] = pluginInfo
	return confToml
}

func createConfigTemp(confpath string, conf *discover.PluginConfig) (string, error) {
	confTempPath := rebuildConfPath(confpath, gopluginTomlTemp)

	b, err := toml.Marshal(conf)
	if err != nil {
		return "", err
	}

	f, err := os.Create(confTempPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = f.Write(b)
	return confTempPath, err
}

func rebuildConfPath(confpath, newPath string) string {
	splittedConf := strings.Split(confpath, "/")
	confDir := splittedConf[0 : len(splittedConf)-1]
	confRebuild := fmt.Sprintf("%s/%s", strings.Join(confDir, "/"), newPath)
	return confRebuild
}
