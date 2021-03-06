package goplugin

import (
	"time"

	"github.com/quadroops/goplugin/pkg/caller/driver"
	"github.com/quadroops/goplugin/pkg/discover"
	"github.com/quadroops/goplugin/pkg/executor"
	"github.com/quadroops/goplugin/pkg/host"
	"github.com/quadroops/goplugin/pkg/process"
	"github.com/quadroops/goplugin/pkg/supervisor"
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
	hostPlugins     PluginMapper
	configChecker   *discover.ConfigChecker
	configParser    *discover.ConfigParser
	processInstance *process.Instance
	identityChecker host.IdentityChecker
}

// Option used to customize default objects
type Option func(*GoPlugin)

// InstallationOptions used to store any options on install
type InstallationOptions struct {
	RetryTimeoutCaller int
}

// PluginMapper used as registry to store plugin's config
type PluginMapper map[string]*PluginConf

// PluginConf used to store plugin's configurations
type PluginConf struct {
	Protocol *ProtocolOption
}

// Registry used as wrapper of executor object
type Registry struct {
	hostPlugins []*GoPlugin
	hosts       []*host.Builder
	exec        *executor.Exec
}

// HostPlugins used to store all plugins from some host
type HostPlugins struct {
	Host    string
	Plugins host.Plugins
}

// PluginSupervisor is main struct used to store supervisor's configurations
type PluginSupervisor struct {
	pluggable   *Registry
	interval    int
	ticker      *time.Ticker
	tickerDone  chan bool
	hostPlugins []*HostPlugins
	runner      *supervisor.Runner
	driver      supervisor.Driver
	handlers    []supervisor.OnErrorHandler
}

// PluginSupervisorOption used to customize supervisor values
type PluginSupervisorOption func(*PluginSupervisor)
