package command

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/quadroops/goplugin/pkg/discover"
	"github.com/quadroops/goplugin/pkg/discover/driver"
	"github.com/spf13/cobra"
)

// discoverCmd used to fetch all information about current plugin's information
// including meta, settings, available plugins and hosts
var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Used to discover available plugins",
	Run: func(cmd *cobra.Command, args []string) {
		defer func() {
			if r := recover(); r != nil {
				log.Fatalf("Panic catched: %v", r)
			}
		}()

		var filepath string

		fileFromArg, err := cmd.Flags().GetString("filepath")
		if err != nil {
			log.Fatalf("Error exploring plugins: %v", err)
		}

		if fileFromArg == "" {
			d := discover.NewConfigChecker(driver.NewOsChecker(), driver.NewDefaultChecker())
			dfile, err := d.Explore()
			if err != nil {
				log.Fatalf("Error exploring plugins: %v", err)
			}

			filepath = dfile
		} else {
			filepath = fileFromArg
		}

		fileReader := driver.NewFileReader()
		tomlParser := driver.NewTomlParser()

		parser := discover.NewConfigParser(tomlParser, fileReader)
		config, err := parser.Load(filepath)
		if err != nil {
			log.Fatalf("Error parse toml: %v", err)
		}

		onlyPlugins, err := cmd.Flags().GetBool("only-plugins")
		if err != nil {
			log.Fatalf("Error exploring plugins: %v", err)
		}

		onlyHosts, err := cmd.Flags().GetBool("only-hosts")
		if err != nil {
			log.Fatalf("Error exploring plugins: %v", err)
		}

		if onlyPlugins {
			renderPlugins(config)
			return
		}

		if onlyHosts {
			renderHosts(config)
			return
		}

		renderTableMeta(config)
		fmt.Println("")
		renderPlugins(config)
		fmt.Println("")
		renderHosts(config)
	},
}

func init() {
	discoverCmd.Flags().StringP("filepath", "f", "", "Set custom goplugin filepath")
	discoverCmd.Flags().BoolP("only-plugins", "", false, "Used to show only available plugins")
	discoverCmd.Flags().BoolP("only-hosts", "", false, "Used to show only available hosts")
}

func renderTableMeta(conf *discover.PluginConfig) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Version", "Author", "Contributors"})

	meta := []string{conf.Meta.Version, conf.Meta.Author, strings.Join(conf.Meta.Contributors, ",")}
	data := [][]string{meta}

	for _, v := range data {
		table.Append(v)
	}

	table.SetCaption(true, "Meta information")
	table.Render()
}

func renderPlugins(conf *discover.PluginConfig) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "MD5", "Author", "Exec Command", "Exec File", "Exec Time", "Protocol Type", "Protocol Port"})

	var data [][]string
	for name, plugin := range conf.Plugins {
		d := []string{name, plugin.MD5, plugin.Author, plugin.Exec, plugin.ExecFile, strconv.Itoa(plugin.ExecTime), plugin.ProtocolType}
		data = append(data, d)
	}

	table.AppendBulk(data)
	table.SetCaption(true, "Plugins information")
	table.Render()
}

func renderHosts(conf *discover.PluginConfig) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Plugins"})

	var data [][]string
	for name, host := range conf.Hosts {
		d := []string{name, strings.Join(host.Plugins, ",")}
		data = append(data, d)
	}

	table.AppendBulk(data)
	table.SetCaption(true, "Hosts information")
	table.Render()
}
