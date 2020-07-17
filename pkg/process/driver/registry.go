package driver

import (
	"fmt"
	"sync"

	"github.com/quadroops/goplugin/pkg/errs"
	"github.com/quadroops/goplugin/pkg/process"
)

type mapper struct {
	data   map[string]process.Plugin

	// need to implement primitiv locking to make sure
	// consistency
	mutex sync.RWMutex
}

// NewRegistry used to create new mapper instance that implement registry builder
func NewRegistry() process.RegistryBuilder {
	data := make(map[string]process.Plugin)
	return &mapper{
		data: data,
	}
}

func (m *mapper) Reset() {
	m.mutex.Lock()
	m.data = make(map[string]process.Plugin)
	m.mutex.Unlock()
}

func (m *mapper) Register(name string, plugin process.Plugin) {
	m.mutex.Lock()
	m.data[name] = plugin
	m.mutex.Unlock()
}

func (m *mapper) IsExist(name string) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	_, exist := m.data[name]
	return exist
}

func (m *mapper) Delete(name string) {
	m.mutex.Lock()
	delete(m.data, name)
	m.mutex.Unlock()
}

func (m *mapper) Get(name string) (process.Plugin, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	plugin, exist := m.data[name]
	if exist {
		return plugin, nil
	}

	return process.Plugin{}, fmt.Errorf("%q: %w", name, errs.ErrPluginNotFound) 
}