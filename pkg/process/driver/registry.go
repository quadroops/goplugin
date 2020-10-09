package driver

import (
	"fmt"

	"github.com/quadroops/goplugin/pkg/errs"
	"github.com/quadroops/goplugin/pkg/process"
)

type mapper struct {
	data map[string]process.Plugin
}

// NewRegistry used to create new mapper instance that implement registry builder
func NewRegistry() process.RegistryBuilder {
	data := make(map[string]process.Plugin)
	return &mapper{
		data: data,
	}
}

func (m *mapper) Reset() {
	m.data = make(map[string]process.Plugin)
}

func (m *mapper) Register(name string, plugin process.Plugin) {
	m.data[name] = plugin
}

func (m *mapper) IsExist(name string) bool {
	_, exist := m.data[name]
	return exist
}

func (m *mapper) Delete(name string) {
	delete(m.data, name)
}

func (m *mapper) Get(name string) (process.Plugin, error) {
	plugin, exist := m.data[name]
	if exist {
		return plugin, nil
	}

	return process.Plugin{}, fmt.Errorf("%q: %w", name, errs.ErrPluginNotFound)
}
