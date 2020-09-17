package goplugin

import (
	"fmt"

	"github.com/hashicorp/go-multierror"

	"github.com/quadroops/goplugin/pkg/errs"
	"github.com/quadroops/goplugin/pkg/executor"
	"github.com/quadroops/goplugin/pkg/host"
)

// Registry used as wrapper of executor object
type Registry struct {
	hosts []*host.Builder
	exec  *executor.Exec
}

// Register used to collect available hosts and run
// executor behind it
func Register(hostPlugins ...*GoPlugin) (*Registry, error) {
	var registries []*executor.Registry
	var hosts []*host.Builder

	if len(hostPlugins) >= 1 {
		for _, h := range hostPlugins {
			host, err := h.Build()
			if err != nil {
				return nil, err
			}

			hosts = append(hosts, host)
			reg := executor.Register(host, h.GetProcessInstance())
			registries = append(registries, reg)
		}
	}

	return &Registry{
		exec:  executor.New(registries...),
		hosts: hosts,
	}, nil
}

// Setup .
func (r *Registry) Setup() error {
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

// FromHost here is just a proxy to call FromHot from executor
func (r *Registry) FromHost(hostname string) (*executor.Container, error) {
	return r.exec.FromHost(hostname)
}
