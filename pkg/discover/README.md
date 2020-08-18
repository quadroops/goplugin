# pkg/discover

Package: `github.com/quadroops/goplugin/pkg/discover`

We need a ways to discover our registered and available plugins.

1. We are using single file configuration.  And using toml as main configuration's format 
2. This configuration's file should be located at user's home dir, such as: `/home/my/.goplugin/config.toml`
3. User able to customize this configuration filepath using environment variable: `GOPLUGIN_DIR`

Package's functionalities:

**Explore**

- Used to discover "quadroops goplugins's path". 
- Check if envvar `GOPLUGIN_DIR` is exist, if exist then read the config's file from the path, if
config file not exist then throw a panic error.  
- If envvar not exist, then by default try to seek user's home dir like `/home/my/.goplugin` and read the config path
from there, if still not exist then throw a panic error.
- If all process success, then should be return a config filepath (string)

**Parser**

- Used to parse configuration from config file and return a `PluginConfig`

---

## Toml Configuration Values

```

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

---

## Usages

```go
import (
	"github.com/quadroops/goplugin/pkg/discover"
	"github.com/quadroops/goplugin/pkg/discover/driver"
)

d := discover.NewConfigChecker(driver.NewOSChecker(), driver.NewDefaultChecker())
conf, err := d.Explore()
if err != nil {
    // error handling 
}

parser := discover.NewConfigParser(driver.NewTomlParser(), driver.NewFileReader())

// will return *discover.PluginConfig if success no error
config, err := parser.Load(conf)
if err != nil {
    // error handling
}
```