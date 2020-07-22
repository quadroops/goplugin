# pkg/executor

Package: `github.com/quadroops/goplugin/pkg/executor`

**Overview**

We need to create new instance of `executor` and register all available hosts.

## Usages

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