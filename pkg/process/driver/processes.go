package driver

import (
	"context"

	"github.com/quadroops/goplugin/pkg/errs"
	"github.com/quadroops/goplugin/pkg/process"
	"github.com/reactivex/rxgo/v2"
)

type proccesses struct {
	suppliers []rxgo.Supplier
	registry  process.RegistryBuilder
}

// NewProcesses used to create new instance of processes
func NewProcesses(registry process.RegistryBuilder) process.ProcessesBuilder {
	suppliers := []rxgo.Supplier{}
	return &proccesses{suppliers: suppliers, registry: registry}
}

func (p *proccesses) Get(name string) (process.Plugin, error) {
	if len(p.suppliers) < 1 {
		return process.Plugin{}, errs.ErrEmptyProcesses
	}

	if p.registry.IsExist(name) {
		return p.registry.Get(name)
	}

	observer, err := p.Listen()
	if err != nil {
		return process.Plugin{}, err
	}

	chPlugin := make(chan rxgo.Item)
	observer.Filter(func(val interface{}) bool {
		plugin, ok := val.(process.Plugin)
		if !ok {
			return false
		}

		return plugin.Name == name
	}, rxgo.WithCPUPool()).Map(func(_ context.Context, val interface{}) (interface{}, error) {
		plugin, ok := val.(process.Plugin)
		if !ok {
			return nil, errs.ErrCastInterface
		}

		return plugin, nil
	}).Send(chPlugin)

	for item := range chPlugin {
		plugin, ok := item.V.(process.Plugin)
		if !ok {
			return process.Plugin{}, errs.ErrCastInterface
		}

		p.registry.Register(name, plugin)
		return plugin, nil
	}

	return process.Plugin{}, errs.ErrPluginNotFound
}

func (p *proccesses) Reset() {
	p.suppliers = []rxgo.Supplier{}
	p.registry.Reset()
}

func (p *proccesses) Remove(name string) error {
	if len(p.suppliers) < 1 {
		return errs.ErrEmptyProcesses
	}

	exist := p.IsExist(name)
	if !exist {
		return errs.ErrPluginNotFound
	}

	observer, err := p.Listen()
	if err != nil {
		return err
	}

	values, err := observer.Filter(func(val interface{}) bool {
		plugin, ok := val.(process.Plugin)
		if !ok {
			return false
		}

		return plugin.Name != name
	}, rxgo.WithCPUPool()).ToSlice(len(p.suppliers))

	if err != nil {
		return err
	}

	var suppliers []rxgo.Supplier
	for _, value := range values {
		plugin, ok := value.(process.Plugin)
		if ok {
			suppliers = append(suppliers, func(_ context.Context) rxgo.Item {
				return rxgo.Of(plugin)
			})
		}
	}

	p.suppliers = suppliers
	p.registry.Delete(name)
	return nil
}

func (p *proccesses) Add(plugin process.Plugin) error {
	if p.IsExist(plugin.Name) {
		return errs.ErrPluginStarted
	}

	p.suppliers = append(p.suppliers, func(_ context.Context) rxgo.Item {
		return rxgo.Of(plugin)
	})

	return nil
}

func (p *proccesses) IsExist(name string) bool {
	observer, err := p.Listen()
	if err != nil {
		return false
	}

	if p.registry.IsExist(name) {
		return true
	}

	out, err := observer.Contains(func(val interface{}) bool {
		v, ok := val.(process.Plugin)
		if !ok {
			return false
		}

		return v.Name == name
	}, rxgo.WithCPUPool()).Get()

	if err != nil {
		return false
	}

	return out.V.(bool)
}

func (p *proccesses) Listen() (rxgo.Observable, error) {
	if len(p.suppliers) < 1 {
		return nil, errs.ErrEmptyProcesses
	}

	return rxgo.Start(p.suppliers), nil
}
