package host

// MD5Checker used to parse a md5 value from given filepath
type MD5Checker interface {
	Parse(file string) (string, error)
}

// Name used for host name
type Name string

// PluginName used for plugin's name
type PluginName string

// Registry used for storing validated plugins
type Registry struct {
	ExecPath string
	ExecArgs []string
	ExecFile string
	ExecTime int
	MD5Sum   string
	RPCType  string
	RPCPort  string
}

// Plugins is a mapper a plugin and their metadata
type Plugins map[PluginName]*Registry

// Host is a mapper between a host and their available plugins
type Host map[Name]Plugins
