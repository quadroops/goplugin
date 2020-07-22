# pkg/process

Package: `github.com/quadroops/goplugin/pkg/process`

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

## Usages

```go
import (
	"github.com/quadroops/goplugin/pkg/process"
	"github.com/quadroops/goplugin/pkg/process/driver"
)

p := process.New(
    driver.NewSubProcess(), 
    driver.NewProcesses(driver.NewRegistry()),
)

// run new subprocess
ch, err := p.Run(1, "test", "test", 1001)
if err != nil {
    // error handling
}

// get plugin info
plugin := <-ch

// kill plugin
err = p.Kill() 

// kill all plugins
errs := p.KillAll()
```