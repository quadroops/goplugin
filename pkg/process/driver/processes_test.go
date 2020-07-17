package driver_test

import (
	"errors"
	"testing"

	"github.com/quadroops/goplugin/pkg/errs"
	"github.com/quadroops/goplugin/pkg/process/driver"
	"github.com/stretchr/testify/assert"
)

func TestAddExistSuccess(t *testing.T) {
	plugin := createMockPlugin("test")

	registry := driver.NewRegistry()
	processes := driver.NewProcesses(registry)
	
	err := processes.Add(plugin)
	assert.NoError(t, err)
	assert.True(t, processes.IsExist(plugin.Name))
}

func TestAddGetSuccess(t *testing.T) {
	plugin := createMockPlugin("test")

	registry := driver.NewRegistry()
	processes := driver.NewProcesses(registry)
	
	err := processes.Add(plugin)
	assert.NoError(t, err)
	assert.True(t, processes.IsExist(plugin.Name))

	p, err := processes.Get("test")
	assert.NoError(t, err)
	assert.Equal(t, plugin.Name, p.Name)
}

func TestGetError(t *testing.T) {
	plugin := createMockPlugin("test")
	registry := driver.NewRegistry()
	processes := driver.NewProcesses(registry)

	err := processes.Add(plugin)
	assert.NoError(t, err)

	_, err = processes.Get("test2")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrPluginNotFound))
}

func TestGetErrorEmptyProcesses(t *testing.T) {
	registry := driver.NewRegistry()
	processes := driver.NewProcesses(registry)
	
	_, err := processes.Get("test2")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrEmptyProcesses))
}

func TestIsExistFromRegistry(t *testing.T) {
	plugin := createMockPlugin("test")
	registry := driver.NewRegistry()
	processes := driver.NewProcesses(registry)

	err := processes.Add(plugin)
	assert.NoError(t, err)
	
	p, err := processes.Get("test")
	assert.NoError(t, err)
	assert.Equal(t, plugin.Name, p.Name)

	assert.True(t, processes.IsExist("test"))
}

func TestAddErrorExist(t *testing.T) {
	plugin := createMockPlugin("test")
	
	registry := driver.NewRegistry()
	processes := driver.NewProcesses(registry)
	
	err := processes.Add(plugin)
	assert.NoError(t, err)

	err = processes.Add(plugin)
	assert.Error(t, err)
}

func TestRemoveUnexistedProcess(t *testing.T) {
	plugin := createMockPlugin("test")
	
	registry := driver.NewRegistry()
	processes := driver.NewProcesses(registry)
	
	err := processes.Add(plugin)
	assert.NoError(t, err)
	assert.True(t, processes.IsExist(plugin.Name))

	err = processes.Remove("test2")
	assert.Error(t, err)
	assert.True(t, errors.Is(errs.ErrPluginNotFound, err))
}

func TestRemoveSuccess(t *testing.T) {
	plugin := createMockPlugin("test")
	plugin2 := createMockPlugin("test2")
	plugin3 := createMockPlugin("test3")

	registry := driver.NewRegistry()
	processes := driver.NewProcesses(registry)
	
	err := processes.Add(plugin)
	assert.NoError(t, err)
	assert.True(t, processes.IsExist(plugin.Name))

	processes.Add(plugin2)
	processes.Add(plugin3)
	assert.True(t, processes.IsExist(plugin2.Name))
	assert.True(t, processes.IsExist(plugin3.Name))

	err = processes.Remove("test")
	assert.NoError(t, err)

	_, err = processes.Listen()
	assert.NoError(t, err)
	assert.False(t, processes.IsExist(plugin.Name))
	assert.True(t, processes.IsExist(plugin2.Name))
	assert.True(t, processes.IsExist(plugin3.Name))
}

func TestRemoveEmptyProcesses(t *testing.T) {
	registry := driver.NewRegistry()
	processes := driver.NewProcesses(registry)
	err := processes.Remove("test")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrEmptyProcesses))
}

func TestResetSuccess(t *testing.T) {
	plugin := createMockPlugin("test")
	registry := driver.NewRegistry()
	processes := driver.NewProcesses(registry)
	
	err := processes.Add(plugin)
	assert.NoError(t, err)
	assert.True(t, processes.IsExist(plugin.Name))

	processes.Reset()
	assert.False(t, processes.IsExist(plugin.Name))
}