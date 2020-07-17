package executor

import (
	"context"

	"github.com/quadroops/goplugin/pkg/errs"
	"github.com/quadroops/goplugin/pkg/caller"
	"github.com/quadroops/goplugin/pkg/host"
	"github.com/quadroops/goplugin/pkg/process"
	"github.com/reactivex/rxgo/v2"
)

// Register used to register a host and their isolated processes
func Register(h *host.Builder, proc *process.Instance) *Registry {
	return &Registry{
		Host: h,
		Process: proc,
	}
} 

// New used to create new instance, with parameter a list of registries
func New(processes ...*Registry) *Exec {
	var sources []rxgo.Supplier
	if len(processes) >= 1 {
		for _, p := range processes {
			rebuild := p 
			source := func(_ context.Context) rxgo.Item {
				return rxgo.Of(rebuild)
			}

			sources = append(sources, source)
		}
	}

	return &Exec{
		processes: sources,
	}
}

// ProcessLength used to count registered processes
func (e *Exec) ProcessLength() int {
	return len(e.processes)
}

// FromHost used to start initialize host process
func (e *Exec) FromHost(host string) (*Container, error) {
	var container Container

	observable := rxgo.Start(e.processes)
	<-observable.Filter(func(i interface{}) bool {
		reg, ok := i.(*Registry)
		if !ok { return false }

		return reg.Host.Hostname == host
	}, rxgo.WithCPUPool()).DoOnNext(func(i interface{}) {
		reg, ok := i.(*Registry)
		if ok {
			container = Container{
				Registry: reg,
			}
		}
	})

	// setup host
	err := container.Setup()
	return &container, err
}

// Setup used to check any available host and install it
func (c *Container) Setup() error {
	// only run this step if host has not been installed
	if !c.installed {
		plugins, err := c.Registry.Host.Install(c.Registry.Host.Setup())
		if err != nil { return err }

		c.plugins = plugins[host.Name(c.Registry.Host.Hostname)]
		c.installed = true
	}

	return nil
}

// IsInstalled used to check if current host has been installed or not
func (c *Container) IsInstalled() bool {
	return c.installed
}

// PluginLength used to count how many available plugins from given host
func (c *Container) PluginLength() int {
	return len(c.plugins) 
}

// Run used to start plugin's process
func (c *Container) Run(name string) error {
	pluginMeta, exist := c.plugins[host.PluginName(name)]
	if !exist { return errs.ErrPluginNotFound }

	chanPlugin, err := c.Registry.Process.Run(
		pluginMeta.ExecTime, 
		name, 
		pluginMeta.ExecPath, 
		pluginMeta.ExecArgs...)

	c.Registry.Process.Watch(chanPlugin, err)
	return err
}

// Get used to create plugin's instance
func (c *Container) Get(name string, builder caller.Builder) (*caller.Plugin, error) {
	pluginMeta, exist := c.plugins[host.PluginName(name)]
	if !exist { return nil, errs.ErrPluginNotFound }

	// need to make sure if current plugin's rpc type supported
	var protocolAllowed bool
	for _, protocol := range caller.AllowedProtocols {
		if pluginMeta.RPCType == protocol {
			protocolAllowed = true
			break
		}
	}

	if !protocolAllowed { return nil, errs.ErrProtocolUnknown }
	return caller.New(pluginMeta, builder(pluginMeta.RPCType, pluginMeta.RPCPort)), nil
}

// GetPluginMeta used to get plugin's metadata
func (c *Container) GetPluginMeta(name string) (*host.Registry, error) {
	pluginMeta, exist := c.plugins[host.PluginName(name)]
	if !exist { return nil, errs.ErrPluginNotFound }

	return pluginMeta, nil
}