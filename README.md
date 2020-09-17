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
package main

import (
    "log"

    "github.com/quadroops/goplugin"
)

func main() {
    hostA := goplugin.New("hostA")
    hostB := goplugin.New("hostB")
    
    // stop all available plugins, when main system is shutting down
    defer func() {
        hostA.GetPluginInstance().KillAll()
        hostB.GetPluginInstance().KillAll()
    }()

    pluggable, err := goplugin.Register(
        hostA,
        hostB,
    )
    if err != nil {
        panic(err)
    }

    err = pluggable.Setup()
    if err != nil {
        panic(err)
    }

    aPlugins, err := pluggable.FromHost("hostA")
    if err != nil {
        panic(err)
    }

    bPlugins, err := pluggable.FromHost("hostB")
    if err != nil {
        panic(err)
    }

    err = aPlugins.Run("plugin1", 8081)
    if err != nil {
        panic(err)
    }

    err = bPlugins.Run("plugin2", 8082)
    if err != nil {
        panic(err)
    }

    // get plugin's instance, you dont need to decide which protocol will be used 
    // just use our helper, or even you can create your own helper as long as, your helper
    // function follow caller.Build type
    //
    // after get your plugin's instance, now you can communicate with your plugin via REST/GRPC
    plugin1, err := aPlugins.Get(
        "plugin1", 
        8081, 
        goplugin.BuildProtocol(&goplugin.ProtocolOption{
            RestAddr: "http://localhost",
            RestTimeout: 10, // this configuration is optional
        })
    )

    if err != nil {
        panic(err)
    }

    // call plugin's ping method
    respPing, err := plugin1.Ping()
    if err != nil {
        panic(err)
    }

    log.Printf("Resp: %s", respPing)

    // call plugin's exec method
    respExec, err := plugin1.Exec("test.command", []byte("hello"))
    if err != nil {
        panic(err)
    }

    log.Printf("Resp: %s", respExec)

    // you can use a same ways for bPlugins
}
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

hostA := goplugin.New("hostA", 
    goplugin.WithCustomConfigChecker(Check), 
    goplugin.WithCustomConfigParser(MongoSourceReader, MongoParse),
)

// stop all available plugins, when main system is shutting down
defer func() {
    hostA.GetPluginInstance().KillAll()
}()

pluggable := goplugin.Register(
    hostA,
)

err = pluggable.Setup()
if err != nil {
    panic(err)
}
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