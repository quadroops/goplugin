package goplugin

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-multierror"

	"github.com/quadroops/goplugin/pkg/caller"
	"github.com/quadroops/goplugin/pkg/errs"
	"github.com/quadroops/goplugin/pkg/executor"
	"github.com/quadroops/goplugin/pkg/host"
)

// Register used to collect available hosts and run
// executor behind it
func Register(hostPlugins ...*GoPlugin) *Registry {
	return &Registry{
		hostPlugins: hostPlugins,
	}
}

// Install used to install available hostPlugins by calling Build
func (r *Registry) Install() (*Registry, error) {
	var hosts []*host.Builder
	var registries []*executor.Registry

	if len(r.hostPlugins) >= 1 {
		for _, h := range r.hostPlugins {
			host, err := h.Build()
			if err != nil {
				return nil, err
			}

			hosts = append(hosts, host)
			reg := executor.Register(host, h.GetProcessInstance())
			registries = append(registries, reg)
		}
	}

	r.hosts = hosts
	r.exec = executor.New(registries...)

	// setup all hosts
	err := r.setup()
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Registry) setup() error {
	if len(r.hosts) >= 1 {
		var errGroups error
		for _, h := range r.hosts {
			container, err := r.exec.FromHost(h.Hostname)
			if err != nil {
				return err
			}

			err = container.Setup()
			if err != nil {
				errGroups = multierror.Append(errGroups, err)
			}
		}

		if errGroups != nil {
			return fmt.Errorf("%w: %q", errs.ErrNoHosts, errGroups)
		}
	}

	return nil
}

// Get plugin's instance.  If plugin not started yet, this step will also
// run the plugin via os subprocess, before using this method you have to make
// sure that all hosts has been installed
func (r *Registry) Get(host, plugin string) (*caller.Plugin, error) {
	if len(r.hosts) < 1 {
		return nil, errs.ErrPluginNotFound
	}

	var hostPlugin *GoPlugin
	for _, hp := range r.hostPlugins {
		if hp.hostName == host {
			hostPlugin = hp
			break
		}
	}

	if hostPlugin == nil {
		return nil, errs.ErrNoHosts
	}

	for _, h := range r.hosts {
		if h.Hostname == host {
			container, err := r.exec.FromHost(host)
			if err != nil {
				return nil, err
			}

			pluginConf, err := hostPlugin.GetPluginConf(plugin)
			if err != nil {
				return nil, err
			}

			meta, err := container.GetPluginMeta(plugin)
			if err != nil {
				return nil, err
			}

			var port int
			if meta.ProtocolType == "rest" {
				port = pluginConf.Protocol.RESTOpts.Port
			} else {
				port = pluginConf.Protocol.GRPCOpts.Port
			}

			// if plugin not ready yet, we need to run it
			if !container.IsPluginReady(plugin) {
				err = container.Run(plugin, port)
				if err != nil {
					return nil, err
				}
			}

			p, err := container.Get(plugin, port, BuildProtocol(pluginConf.Protocol))
			if err != nil {
				return nil, err
			}

			return p, nil
		}
	}

	return nil, errs.ErrPluginNotFound
}

// KillPlugins used to kill all plugins from all installed hosts
func (r *Registry) KillPlugins() {
	if len(r.hostPlugins) >= 1 {
		for _, host := range r.hostPlugins {
			log.Printf("Killing all plugins from host: %s ...", host.hostName)
			h := host.GetProcessInstance()
			h.KillAll()
		}
	}
}
