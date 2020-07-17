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
	Version      string
	Author       string
	Contributors []string
}

// PluginSettings used to save all global [settings] informations
type PluginSettings struct {
	Debug bool
}

// PluginInfo used to save all plugin's basic informations
type PluginInfo struct {
	Author   string
	MD5      string
	Exec     string
	ExecArgs []string `toml:"exec_args"`
	ExecFile string   `toml:"exec_file"`
	ExecTime int      `toml:"exec_time"`
	RPCType  string   `toml:"comm_type"`
	RPCAddr  string   `toml:"comm_port"`
}

// PluginHost used to save all registered service's plugins
type PluginHost struct {
	Plugins []string
}

// PluginConfig used to store all plugin configuration values
type PluginConfig struct {
	Meta     PluginMeta
	Settings PluginSettings
	Plugins  map[string]PluginInfo
	Hosts    map[string]PluginHost
}

// Host used as main host's configurations
type Host struct {
	Name   string
	Config PluginConfig
}
