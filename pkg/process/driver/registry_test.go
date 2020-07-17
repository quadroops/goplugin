package driver_test

import (
	"testing"

	"github.com/quadroops/goplugin/pkg/process"
	"github.com/quadroops/goplugin/pkg/process/driver"
	"github.com/stretchr/testify/assert"
)

func createMockPlugin(name string) process.Plugin {
	return process.Plugin{
		ID: process.ID(1),
		Name: name,
	}
}

func TestRegistryRegisterSuccess(t *testing.T) {
	registry := driver.NewRegistry()
	registry.Register("test", createMockPlugin("test"))
	assert.True(t, registry.IsExist("test"))
}

func TestRegistryDeleteSuccess(t *testing.T) {
	registry := driver.NewRegistry()
	registry.Register("test", createMockPlugin("test"))
	registry.Delete("test")

	assert.False(t, registry.IsExist("test"))
}

func TestRegistryGetSuccess(t *testing.T) {
	registry := driver.NewRegistry()
	registry.Register("test", createMockPlugin("test"))

	plugin, err := registry.Get("test")
	assert.NoError(t, err)
	assert.Equal(t, plugin.Name, "test")
}

func TestRegistryGetErrorNotFound(t *testing.T) {
	registry := driver.NewRegistry()
	plugin, err := registry.Get("test")
	assert.Error(t, err)
	assert.Empty(t, plugin.Name)
}