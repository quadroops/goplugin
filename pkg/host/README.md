# pkg/host

Package: `github.com/quadroops/goplugin/pkg/host`

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

## Usages

```go
import (
    "github.com/quadroops/goplugin/pkg/host"
    "github.com/quadroops/goplugin/pkg/host/driver"
)

h := host.New("hostname", config, driver.NewMD5Check())
installed, err := h.Install(h.Setup())

if err != nil {
    // error handling
}
```