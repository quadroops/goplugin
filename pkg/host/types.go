package host

// IdentityChecker used to check identity for security purpose
// By default, we will use md5 file checker to check plugin's md5 value
type IdentityChecker interface {
	Parse(source string) (string, error)
}

// Name used for host name
type Name string

// PluginName used for plugin's name
type PluginName string

// Registry used for storing validated plugins
type Registry struct {
	ExecPath     string
	ExecArgs     []string
	ExecFile     string
	ExecTime     int
	MD5Sum       string
	ProtocolType string
}

// Plugins is a mapper a plugin and their metadata
type Plugins map[PluginName]*Registry

// Host is a mapper between a host and their available plugins
type Host map[Name]Plugins
