package goplugin

import (
	"github.com/quadroops/goplugin/pkg/caller/driver"
	"github.com/quadroops/goplugin/pkg/discover"
	"github.com/quadroops/goplugin/pkg/executor"
	"github.com/quadroops/goplugin/pkg/host"
	"github.com/quadroops/goplugin/pkg/process"
)

// ProtocolOption used to configure rest or grpc options
// You have to choose between rest or grpc, if your plugin
// using rest, than ignore grpc option, and vice versa, but you
// can't ignore them both
type ProtocolOption struct {
	RESTOpts *driver.RESTOptions
	GRPCOpts *driver.GrpcOptions
}

// GoPlugin used as main struct to store goplugin's state
type GoPlugin struct {
	hostName        string
	configChecker   *discover.ConfigChecker
	configParser    *discover.ConfigParser
	processInstance *process.Instance
	identityChecker host.IdentityChecker
}

// Option used to customize default objects
type Option func(*GoPlugin)

// PluginMapper used as registry to store plugin's config
type PluginMapper map[string]*PluginConf

// PluginConf used to store plugin's configurations
type PluginConf struct {
	Port     int
	Protocol *ProtocolOption
}

// Registry used as wrapper of executor object
type Registry struct {
	hosts []*host.Builder
	exec  *executor.Exec
}
