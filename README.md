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
