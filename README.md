# quadroops/goplugin 

This package is a library to implement modular and pluggable architecture for application using
Golang's language.

This library inspired from [hashicorp's plugin architecture](https://github.com/hashicorp/go-plugin), but just to implement
it in different ways, such as:

1. Using single file configuration to explore all available plugins (customizable)
2. Adding REST as an alternatives of communication protocols 

We have a same direction with hashicorp's plugin:

1. Using subprocess
2. Supported multiple languages (as long as support with GRPC/REST)

Why we doesn't support [golang's native plugin](https://golang.org/pkg/plugin/) ? It's because using `Symbol` for exported variable 
and methods, it's easy for simple plugin but it will harder to maintain for complex application

---

## Flows

```
// need a way to discover available plugins
[host] -----<discover>---> [plugin]


// start subprocess, should be sharing a same PID
[host] ------<exec>-----> [plugin]


// communicate via REST or GRPC 
[host] ---<communicate>------> [plugin]
```

We have three main flows:

- Discover
- Exec
- Communication via rpc

Current supported communication protocols:

- [REST](https://en.wikipedia.org/wiki/Representational_state_transfer)
- [GRPC](https://grpc.io/)

---

## Usages

Basic usages:

```go
// need to save plugin's configuration ports
// assume this plugin using rest
hostAPlugin1Conf := goplugin.PluginConf{
	Protocol: &goplugin.ProtocolOption{
		RESTOpts: &driver.RESTOptions{
			Addr: "localhost",
			Port: 8080,
			Timeout: 10, // in seconds
		}
	}
}

// assume this plugin using grpc
hostAPlugin2Conf := goplugin.PluginConf{
	Protocol: &goplugin.ProtocolOption{
		GRPCOpts: &driver.GRPCOptions{
			Addr: "host",
			Port: 8081,
		}
	}
}

// each hosts will have to manage their own plugin's configurations
// we need to map available plugins to their specific configurations
hostA := goplugin.New("hostA").SetupPlugin(
	goplugin.Map("plugin1", hostAPlugin1PConf), 
	goplugin.Map("plugin2", hostAPlugin2Conf),
)

// you can register multiple hosts here
// and whel Install method called, it will install all plugins from
// all registered hosts
pluggable, err := goplugin.Register(hostA).Install()
if err != nil {
	panic(err)
}

defer func() {
	// for each successfull installed hosts need to shutdown
	// their plugins if something bad happened
	pluggable.KillPlugins()
}

// when load a plugin, you need to specify from which host
// each host will have their own plugin registries
// if plugin not started yet, it should also run plugin via os subprocess
plugin1, err := pluggable.Get("hostA", "plugin1")
if err != nil {
	panic(err)
}

// start working with available plugins
respPing, err := plugin1.Ping()
if err != nil {
	panic(err)
}

log.Printf("Response ping: %s", respPing)
```

Assume you need to customize config `Checker` and `Parser`, imagine you want to using a database
like `Mongodb` as configuration source.

```go
// implement discover.Checker
func Check() string {
    return 'mongodb://localhost:27017/dbname'
}

// implement discover.SourceReader
func MongoSourceReader(source string) ([]byte, error) {
    // parse configurations from mongodb
}

// implement discover.Parser
func MongoParse(content []byte) (*discover.PluginConfig, error) {
    // map configurations to discover.PluginConfig
}

// need to save plugin's configuration ports
// assume this plugin using rest
hostAPlugin1Conf := goplugin.PluginConf{
	Protocol: &goplugin.ProtocolOption{
		RESTOpts: &driver.RESTOptions{
			Addr: "localhost",
			Port: 8080,
			Timeout: 10, // in seconds
		}
	}
}

hostA := goplugin.New("hostA", 
    goplugin.WithCustomConfigChecker(Check), 
    goplugin.WithCustomConfigParser(MongoSourceReader, MongoParse),
).SetupPlugin(goplugin.Map("plugin1", hostAPlugin1PConf))

// you can register multiple hosts here
// and whel Install method called, it will install all plugins from
// all registered hosts
pluggable, err := goplugin.Register(hostA).Install()
if err != nil {
	panic(err)
}

defer func() {
	// for each successfull installed hosts need to shutdown
	// their plugins if something bad happened
	pluggable.KillPlugins()
}

// when load a plugin, you need to specify from which host
// each host will have their own plugin registries
// if plugin not started yet, it should also run plugin via os subprocess
plugin1, err := pluggable.Get("hostA", "plugin1")
if err != nil {
	panic(err)
}

// start working with available plugins
respPing, err := plugin1.Ping()
if err != nil {
	panic(err)
}

log.Printf("Response ping: %s", respPing)
```

Customizated components available for:

- ConfigChecker (default: get toml file path) 
- ConfigParser (default: toml parser)
- ProcessInstance (default: subprocess)
- HostIdentityChecker (default: MD5Checker)

Available options for customization:

- `goplugin.WithCustomConfigChecker`
- `goplugin.WithCustomConfigParser`
- `goplugin.WithCustomIdentityChecker`
- `goplugin.WithCustomProcess`

---

## Examples

Please refer to our [examples](https://github.com/quadroops/goplugin)