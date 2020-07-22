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