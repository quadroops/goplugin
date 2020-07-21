package discover

// Checker is basic interface and should be implemented
// by any objects that want to check config file path
type Checker interface {
	// Check should be return a string of config file path
	// and just throw a panic if current file not exist
	Check() string
}

// Parser is an interface to abstract main function to parse toml's file
type Parser interface {
	Parse(content []byte) (*PluginConfig, error)
}

// FileReader is an interface to abstract filesystem activity
type FileReader interface {
	ReadFile(filepath string) ([]byte, error)
}

// PluginInit used to save initialize states
type PluginInit struct {
	SvcName        string
	ConfigFilePath string
}

// PluginMeta used to save all [meta] informations
type PluginMeta struct {
	Version      string   `toml:"version"`
	Author       string   `toml:"author"`
	Contributors []string `toml:"contributors"`
}

// PluginSettings used to save all global [settings] informations
type PluginSettings struct {
	Debug bool `toml:"debug"`
}

// PluginInfo used to save all plugin's basic informations
type PluginInfo struct {
	Author       string   `toml:"author"`
	MD5          string   `toml:"md5"`
	Exec         string   `toml:"exec"`
	ExecArgs     []string `toml:"exec_args"`
	ExecFile     string   `toml:"exec_file"`
	ExecTime     int      `toml:"exec_time"`
	ProtocolType string   `toml:"protocol_type"`
}

// PluginHost used to save all registered service's plugins
type PluginHost struct {
	Plugins []string `toml:"plugins"`
}

// PluginConfig used to store all plugin configuration values
type PluginConfig struct {
	Meta     PluginMeta            `toml:"meta"`
	Settings PluginSettings        `toml:"settings"`
	Plugins  map[string]PluginInfo `toml:"plugins"`
	Hosts    map[string]PluginHost `toml:"hosts"`
}

// Host used as main host's configurations
type Host struct {
	Name   string
	Config PluginConfig
}
