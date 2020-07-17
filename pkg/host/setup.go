package host

import (
	"context"
	"fmt"

	"github.com/quadroops/goplugin/pkg/errs"
	"github.com/quadroops/goplugin/pkg/discover"
	"github.com/quadroops/goplugin/pkg/host/flow"
	"github.com/reactivex/rxgo/v2"
)

// Builder used to create host metadata
type Builder struct {
	Hostname   string
	Config     *discover.PluginConfig
	md5Checker MD5Checker
}

// New used to create new instance of host
func New(hostname string, config *discover.PluginConfig, md5Checker MD5Checker) *Builder {
	return &Builder{
		Hostname:   hostname,
		Config:     config,
		md5Checker: md5Checker,
	}
}

// Setup used to get available plugins based on hostname
func (b *Builder) Setup() Plugins {
	if b.Config == nil {
		return nil
	}

	source := func(ctx context.Context) rxgo.Item {
		return rxgo.Of(b.Config)
	}

	hostPlugins := make(Plugins)
	f := flow.NewSetup(b.Hostname, b.Config)
	observable := rxgo.Start([]rxgo.Supplier{source}).
		Filter(f.FilterByHostname, rxgo.WithCPUPool()).
		Map(f.MapToPlugins, rxgo.WithCPUPool()).
		FlatMap(f.FlatToNewObservable)

	<-observable.
		Filter(f.FilterByPlugin, rxgo.WithCPUPool()).
		DoOnNext(func(i interface{}) {
			plugin, ok := i.(string)
			if ok {
				pluginInfo, exist := b.Config.Plugins[plugin]
				if exist {
					hostPlugins[PluginName(plugin)] = &Registry{
						ExecFile: pluginInfo.ExecFile,
						ExecArgs: pluginInfo.ExecArgs,
						ExecPath: pluginInfo.Exec,
						ExecTime: pluginInfo.ExecTime,
						MD5Sum:   pluginInfo.MD5,
						RPCPort:  pluginInfo.RPCAddr,
						RPCType:  pluginInfo.RPCType,
					}
				}
			}
		})

	return hostPlugins
}

// Install used to validate given available plugins, and return a Host
// a mapper between host and their plugins
func (b *Builder) Install(plugins Plugins) (Host, error) {
	if plugins == nil {
		return nil, fmt.Errorf("%w", errs.ErrEmptyPlugins)
	}

	host := make(Host)
	rebuildlugins := make(Plugins)

	source := func(_ context.Context, next chan<- rxgo.Item) {
		for name, p := range plugins {
			flowInstallRegistry := flow.RegistryProxy{
				ExecFile: p.ExecFile,
				ExecArgs: p.ExecArgs,
				ExecPath: p.ExecPath,
				ExecTime: p.ExecTime,
				MD5Sum: p.MD5Sum,
				RPCPort: p.RPCPort,
				RPCType: p.RPCType,
			}

			flowPlugin := flow.Plugin{
				Name: string(name),
				Registry: flowInstallRegistry,
			}

			next <- rxgo.Of(flowPlugin) 
		}
	}

	f := flow.NewInstall(b.md5Checker)
	<-rxgo.Defer([]rxgo.Producer{source}).
		Filter(f.FilterByExecFile, rxgo.WithCPUPool()).
		Filter(f.FilterByMD5, rxgo.WithCPUPool()).
		DoOnNext(func(v interface{}) {
			plugin, ok := v.(flow.Plugin)
			if ok {
				rebuildlugins[PluginName(plugin.Name)] = &Registry{
					ExecFile: plugin.Registry.ExecFile,
					ExecArgs: plugin.Registry.ExecArgs,
					ExecPath: plugin.Registry.ExecPath,
					ExecTime: plugin.Registry.ExecTime,
					MD5Sum: plugin.Registry.MD5Sum,
					RPCPort: plugin.Registry.RPCPort,
					RPCType: plugin.Registry.RPCType,
				}
			}
		})

	if len(rebuildlugins) < 1 {
		return nil, errs.ErrNoPlugins
	}

	host[Name(b.Hostname)] = rebuildlugins 
	return host, nil
}
