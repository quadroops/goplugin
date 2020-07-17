# quadroops/goplugin 

This package is a library to implement modular and pluggable architecture for application using
Golang's language.

This library inspired from [hashicorp's plugin architecture](https://github.com/hashicorp/go-plugin), but just to implement
it in different ways, such as:

1. Using single file configuration to explore all available plugins
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

## Discover

Package: `github.com/quadroops/goplugin/discover`

We need a ways to discover our registered and available plugins.

1. We are using single file configuration.  And using toml as main configuration's format 
2. This configuration's file should be located at user's home dir, such as: `/home/my/.goplugin/config.toml`
3. User able to customize this configuration filepath using environment variable: `GOPLUGIN_DIR`

An example of valid configuration values:

```toml

[meta]
version = "1.0.0"
author = "hiraq|hiraq.dev@gmail.com"
contributors = [
    "a|a@test.com",
    "b|b@test.com
]

# global configurations
[settings]
debug = true 

# Used as main plugin registries
#
# a plugin should provide basic five informations about their self
#
# - Author
# - Md5Sum.  To make sure that we should not be able to exec a plugin which will harm us
# - Exec path
# - Exec start time. Used to wait a plugin to start processes until they are ready to consume by caller
# - Communication type (grpc/rest)
#
# Each of registered plugin MUST have a unique's name
[plugins]

    [plugins.name_1]
    author = "author_1|author_1@gmail.com"
    md5 = "d194b7bad208c2ddfa0ef597fd4abcc5"
    exec = "/path/to/exec"
    exec_args = ["--port", "8080"]
    exec_file = "/path/to/exec"
    exec_time = 5
    comm_type = "grpc"
    comm_port = "8080"
    
    [plugins.name_2]
    author = "author_2|author_2@gmail.com"
    md5 = "d194b7bad208c2ddfa0ef597fd4abcc5"
    exec = "/path/to/exec"
    exec_args = ["--port", "8080"]
    exec_file = "/path/to/exec"
    exec_time = 10
    comm_type = "rest"
    comm_port = "8081"
    
    [plugins.name_3]
    author = "author_3|author_3@gmail.com"
    md5 = "d194b7bad208c2ddfa0ef597fd4abcc5"
    exec = "/path/to/exec"
    exec_args = ["--port", "8080"]
    exec_file = "/path/to/exec"
    exec_time = 20
    comm_type = "nano"
    comm_port = "8082"

# Used as service registries
# A service is an application that consume / using plugins
# They are not allowed to access a plugin that not registered to service's plugin registries
[hosts]

    [hosts.host_1]
    plugins = ["name_1", "name_2"]
    
    [hosts.host_2]
    plugins = ["name_3"]
```

Package's functionalities:


**Initialize**

- Used to discover "quadroops goplugins's path". 
- Check if envvar `GOPLUGIN_DIR` is exist, if exist then read the config's file from the path, if
config file not exist then throw a panic error.  
- If envvar not exist, then by default try to seek user's home dir like `/home/my/.goplugin` and read the config path
from there, if still not exist then throw a panic error.
- If all process success, then should be return a config filepath (string)

**Parser**

- Used to parse configuration from config file and return a `PluginConfig`

---

## Exec

There are two domain need to manage:

- Host
- Process

### Host

Package: `github.com/quadroops/goplugin/host`

**Overview**

- Setup hostname. Each of application/service using `goplugin` must have a unique hostname
- Validate available plugins.  This process will check all available plugins, and also check for their `MD5Sum`, 
should return a map of plugin name and their exec path

**Behaviors**

- Setup hostname
- Install.  This process will automatically explore all available plugins based on given `hostname` and `PluginConfig`

    - Get a list of plugins 
    - Check if exec path is exist
    - Get `MD5Sum` and compare it with `PluginConfig`
    - Create a map based on plugin's name and their exec path

---

### Process

Package: `github.com/quadroops/goplugin/process`

**Overview**

This package provide functionalities to manage subprocesses, including running and killing process.  

Flows:

- Run individual plugin, based on plugin's `exec` path.  A host can run a plugin based on their plugin's name
- When running a plugin, a host should be able to wait based plugin's `exec_time` to make sure plugin has been running successfully
- Each of executed plugins, will be isolated and have their own process, mapped using plugin's name and their `ID`
- A host can kill individual or all executed plugins based on their `ID`
- When a host killed, should be able to kill all executed plugins, to make sure there are no zombie process running on OS

**Behaviors**

- `Run` individual plugin based on plugin's name
- `Kill` stop individual plugin's process
- `KillAll` kill all running plugins from the `Registry`

---

### Executor 

Package: `github.com/quadroops/goplugin/executor`

**Overview**

We need to create new instance of `executor` and register all available hosts.

**Usages**

```go
h1 := host.New("host1", config1, md5Checker)
h2 := host.New("host2", config2, md5Checker)

// each of hosts, will have their own isolated processes
process1 := process.New(runner, processes)
process2 := process.New(runner, processes)

run := executor.New(
    executor.Register(h1, process1), 
    executor.Register(h2, process2),
)
```

And then, we can run a plugin from a host

```go
// each of hosts will have their own self isolated container 
// each of container is a new instance 
container1, err := run.FromHost("host1")
container2, err := run.FromHost("host2")

// start plugins
err = container1.Run("plugin_1")
err = container1.Run("plugin_2")
err = container2.Run("plugin_3")
```

---

## Caller 

How should we handle a communication between application's host and our plugins ?

Each of plugins must implement a same communication channels.

There are two available actions:

1. Ping
2. Exec

**Ping**

Used to check plugin's process health.  It's simple just a request with payloand message of "ping" 
and should've got a response of "pong"

**Exec**

Payload:

- Command: A string of command action's name
- Content: An encoded byte data 

A host and plugin should be agreed to define command's names and also their content payload (byte encoded) you can use 
protocol buffer or bson or messagepack to encode the data into binary, and you have to convert the binary data into string
using `hex` mechanism, example `hex` in golang:

```go

// payload is []byte
resp, err := client.Exec(context.Background(), &pbPlugin.ExecRequest{
    Command: cmdName,
    Payload: hex.EncodeToString(payload),
})

```

A plugin may have multiple actions, but only has a single entrypoint, an `exec` path.  When a plugin have multiple actions,
they can manage it by routing the action based on command's name.

For an example, a plugin of analyze log's data, it may have processes like:

- Store management, command name: `plugin.log.store`
- Filtering management, command name: `plugin.log.filter`
- Retention management, command name: `plugin.log.retention`

A command name should be used like routing key in RabbitMQ.

---

### Protocol: REST

There are two endpoints need to implement :

Endpoint: `/ping`

```
Method: GET
Request headers:
- X-Host: <hostname>

Response status header: 200 (OK)
Response body:
{
    "status": "success",
    "data": {
        "response": "pong"
    }
}
```

Endpoint: `/exec`

```
Method: POST
Request header:
- X-Host: <hostname>

Request body:
{
    "command": "<defined.command.name>",
    "payload": "<hexadecimal_encoded_bytes_in_string>" 
}

Response status header: 202 (accepted)
Response body:
{
    "status": "success",
    "data": {
        "response": "<hexadecimal_encoded_bytes_in_string>"
    }
}
```

An example of hex encoded byte (in golang):

```go
import "encoding/hex"

payload := "test"
b, _ := msgpack.Marshal(payload)
str := hex.EncodeToString(b) 
```

Unmarshalling

```go
// will return a bytes, you doesn't need to decode the hex
b, err := rest.Exec("rest.testing", []byte("test"))
assert.NoError(t, err)

var s string
err = msgpack.Unmarshal(b, &s)
assert.NoError(t, err)
```

Why hex ? We cannot take any assumptions here what kind of data type you will used, the most safety way to manage the request and response data is by sending
the data in the byte's form and encode the data (byte) using hex.  By using this method, you can manage you request and response payload based on your needs.
Please remember, even your request or response data is on "string" type, you have to convert it into byte and encode it using hex.

Data flows:

```
<your_data_type> -> convert into byte -> encode with hex
```

---

### Protocol: GRPC

For grpc's channel, a plugin must use and implement available proto's file

```protobuf
syntax = "proto3";
package plugin;

import "google/protobuf/empty.proto";

service Plugin {
    rpc Ping (google.protobuf.Empty) returns (PingResponse);
    rpc Exec (ExecRequest) returns (ExecResponse);
}

message Data {
    string Response = 1;
}

message PingResponse {
    string Status = 1;
    Data Data = 2;
}

message ExecRequest {
    string Command = 1;
    string Payload = 2;
}

message ExecResponse {
    string Status = 1;
    Data Data = 2;
}
```

---

We should using a same interface for our communication protocols :

```go
type Caller interface {
    Ping() (string, error)
    Exec(cmdName string, payload []byte) ([]byte, error)
}
```

Each of communication protocol, must implement `Caller` interface, usages:

```go
// a communication caller's helper used to construct a communicaton instance
// between rest or grpc
buildCaller := func(commType, commPort string) caller.Caller {
    switch commType {
        case "rest":
            return driver.NewREST(commPort)
        case "grpc":
            return driver.NewGRPC(func() (pbPlugin.PluginClient, error) {
                conn, err := grpc.Dial(commPort)
                if err != nil {
                    return nil, err
                }

                return pbPlugin.NewPluginClient(conn), nil
            })
    }

    return &caller.Plugin{}
}

// get plugin's caller from container by passing plugin's name and the buildCaller helper
// so we doesnt need to kno what kind of communication's types used
plugin, err := container.Get("plugin_1", buildCaller)

// ping request
resp, err := plugin.Ping()

// exec request
resp, err := plugin.Exec("plugin.command.name", []byte("payload example"))
```

## Plugin

To create a plugin, an author must provide a `goplugin.toml` inside their plugin's directory, with the format like below:

```toml
[plugin]
name = "pluginname"
author = "author|author@email.com"
md5 = "d194b7bad208c2ddfa0ef597fd4abcc5"
exec = "/path/to/exec"
exec_args = ["--myparam", "mycustom"]
exec_file = "/path/to/exec"
exec_time = 5
comm_type = "grpc"
hosts = ["host_1", "host_2"]
```

Your plugin must be can to start from command line and provide a parameter to set a port using `--port` to using custom port:

```
python ./main.py --port 8080

node ./main.js --port 8080
```

And your plugin must have to implement a same communication interface methods for REST or GRPC:

- Ping
- Exec