package executor_test

import (
	"context"
	"errors"
	"testing"

	"github.com/quadroops/goplugin/pkg/caller"
	discoverDriver "github.com/quadroops/goplugin/pkg/discover/driver"
	"github.com/quadroops/goplugin/pkg/errs"
	"github.com/quadroops/goplugin/pkg/executor"
	"github.com/quadroops/goplugin/pkg/host"
	"github.com/quadroops/goplugin/pkg/process"

	callerMock "github.com/quadroops/goplugin/pkg/caller/mocks"
	hostMock "github.com/quadroops/goplugin/pkg/host/mocks"
	processMock "github.com/quadroops/goplugin/pkg/process/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	tomlContent = `
	[meta]
	version = "1.0.0"
	author = "hiraq|hiraq@ruangguru.com"
	contributors = [
		"a|a@ruangguru.com",
		"b|b@ruangguru.com"
	]

	# global configurations
	[settings]
	debug = true 

	# Used as main plugin registries
	#
	# a plugin should provide basic five informations about their self
	#
	# - Author
	# - Md5Sum.  To make sure that we should not be able to exec a plugin which will harm us
	# - Exec path
	# - Exec start time. Used to wait a plugin to start processes until they are ready to consume by caller
	# - Rpc type (grpc/rest/nano)
	#
	# Each of registered plugin MUST have a unique's name
	[plugins]

		[plugins.name_1]
		author = "author_1|author_1@gmail.com"
		md5 = "d41d8cd98f00b204e9800998ecf8427e"
		exec = "./tmp/test"
    	exec_file = "./tmp/test"
		exec_time = 5
		comm_type = "grpc"
		
		[plugins.name_2]
		author = "author_2|author_2@gmail.com"
		md5 = "d41d8cd98f00b204e9800998ecf8427e"
		exec = "./tmp/test"
    	exec_file = "./tmp/test"
		exec_time = 10
		comm_type = "grpc"
		
		[plugins.name_3]
		author = "author_3|author_3@gmail.com"
		md5 = "d41d8cd98f00b204e9800998ecf8427e"
		exec = "./tmp/test"
    	exec_file = "./tmp/test"
		exec_time = 20
		comm_type = "unknown"
		
	# Used as service registries
	# A service is an application that consume / using plugins
	# They are not allowed to access a plugin that not registered to service's plugin registries
	[hosts]

		[hosts.host_1]
		plugins = ["name_1", "name_2", "name_3"]		
		
		[hosts.host_2]
		plugins = []		

		[hosts.host_3]
		plugins = ["name_2", "name_3"]
	`
)

func createMockPlugin(name string) process.Plugin {
	_, cancel := context.WithCancel(context.Background())
	return process.Plugin{
		Kill: cancel,
		Name: name,
	}
}

func createMockChanPlugin(plugin process.Plugin) <-chan process.Plugin {
	c := make(chan process.Plugin)
	go func() {
		c <- plugin
		close(c)
	}()

	return c
}

func TestRegisterSuccess(t *testing.T) {
	toml, err := discoverDriver.NewTomlParser().Parse([]byte(tomlContent))
	assert.NoError(t, err)
	assert.Len(t, toml.Plugins, 3)
	assert.Len(t, toml.Hosts, 3)

	md5 := new(hostMock.MD5Checker)
	h := host.New("host_1", toml, md5)
	assert.NotNil(t, h)

	runner := new(processMock.Runner)
	processes := new(processMock.ProcessesBuilder)
	p := process.New(runner, processes)

	registry := executor.Register(h, p)
	assert.NotNil(t, registry)
	assert.Equal(t, "host_1", registry.Host.Hostname)
}

func TestExecNewSuccess(t *testing.T) {
	toml, err := discoverDriver.NewTomlParser().Parse([]byte(tomlContent))
	assert.NoError(t, err)

	md5 := new(hostMock.MD5Checker)
	h := host.New("host_1", toml, md5)
	h2 := host.New("host_2", toml, md5)
	h3 := host.New("host_3", toml, md5)

	runner := new(processMock.Runner)
	processes := new(processMock.ProcessesBuilder)
	p := process.New(runner, processes)
	p2 := process.New(runner, processes)
	p3 := process.New(runner, processes)

	exec := executor.New(
		executor.Register(h, p),
		executor.Register(h2, p2),
		executor.Register(h3, p3),
	)

	assert.Equal(t, exec.ProcessLength(), 3)
}

func TestExecEmpty(t *testing.T) {
	exec := executor.New()
	assert.Equal(t, exec.ProcessLength(), 0)
}

func TestExecFromHostSuccess(t *testing.T) {
	toml, err := discoverDriver.NewTomlParser().Parse([]byte(tomlContent))
	assert.NoError(t, err)

	md5 := new(hostMock.MD5Checker)
	md5.On("Parse", mock.Anything).Return("d41d8cd98f00b204e9800998ecf8427e", nil)

	h := host.New("host_1", toml, md5)
	h2 := host.New("host_2", toml, md5)
	h3 := host.New("host_3", toml, md5)

	runner := new(processMock.Runner)
	processes := new(processMock.ProcessesBuilder)
	p := process.New(runner, processes)
	p2 := process.New(runner, processes)
	p3 := process.New(runner, processes)

	exec := executor.New(
		executor.Register(h, p),
		executor.Register(h2, p2),
		executor.Register(h3, p3),
	)

	container1, err := exec.FromHost("host_1")
	assert.NoError(t, err)
	assert.True(t, container1.IsInstalled())
	assert.Equal(t, 3, container1.PluginLength())

	_, err = exec.FromHost("host_2")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrNoPlugins))

	container3, err := exec.FromHost("host_3")
	assert.NoError(t, err)
	assert.True(t, container3.IsInstalled())
	assert.Equal(t, container3.PluginLength(), 2)
}

func TestRunSuccess(t *testing.T) {
	toml, err := discoverDriver.NewTomlParser().Parse([]byte(tomlContent))
	assert.NoError(t, err)

	md5 := new(hostMock.MD5Checker)
	md5.On("Parse", mock.Anything).Return("d41d8cd98f00b204e9800998ecf8427e", nil)

	h := host.New("host_1", toml, md5)

	mockPlugin := createMockPlugin("test")
	runner := new(processMock.Runner)
	runner.On("Run", 5, "name_1", "./tmp/test", 1001).Once().Return(createMockChanPlugin(mockPlugin), nil)

	processes := new(processMock.ProcessesBuilder)
	processes.On("IsExist", "name_1").Once().Return(false)
	processes.On("Add", mock.Anything).Once().Return(nil)

	p := process.New(runner, processes)

	exec := executor.New(
		executor.Register(h, p),
	)

	container1, err := exec.FromHost("host_1")
	assert.NoError(t, err)

	err = container1.Run("name_1", 1001)
	assert.NoError(t, err)
}

func TestRunNoPlugin(t *testing.T) {
	toml, err := discoverDriver.NewTomlParser().Parse([]byte(tomlContent))
	assert.NoError(t, err)

	md5 := new(hostMock.MD5Checker)
	md5.On("Parse", mock.Anything).Return("d41d8cd98f00b204e9800998ecf8427e", nil)

	h := host.New("host_1", toml, md5)

	runner := new(processMock.Runner)
	processes := new(processMock.ProcessesBuilder)

	p := process.New(runner, processes)

	exec := executor.New(
		executor.Register(h, p),
	)

	container1, err := exec.FromHost("host_1")
	assert.NoError(t, err)

	err = container1.Run("unknown", 1001)
	assert.Error(t, err)
	assert.True(t, errors.Is(errs.ErrPluginNotFound, err))
}

func TestGetPluginSuccess(t *testing.T) {
	mockCaller := new(callerMock.Caller)
	toml, err := discoverDriver.NewTomlParser().Parse([]byte(tomlContent))
	assert.NoError(t, err)

	md5 := new(hostMock.MD5Checker)
	md5.On("Parse", mock.Anything).Return("d41d8cd98f00b204e9800998ecf8427e", nil)

	h := host.New("host_1", toml, md5)
	runner := new(processMock.Runner)
	processes := new(processMock.ProcessesBuilder)
	p := process.New(runner, processes)

	exec := executor.New(
		executor.Register(h, p),
	)

	container, err := exec.FromHost("host_1")
	assert.NoError(t, err)
	assert.True(t, container.IsInstalled())
	assert.Equal(t, 3, container.PluginLength())

	transporter, err := container.Get("name_1", 1001, func(rpcType string, port int) caller.Caller {
		return mockCaller
	})

	assert.NoError(t, err)
	assert.Equal(t, "grpc", transporter.Meta.ProtocolType)
}

func TestGetPluginNotFound(t *testing.T) {
	mockCaller := new(callerMock.Caller)
	toml, err := discoverDriver.NewTomlParser().Parse([]byte(tomlContent))
	assert.NoError(t, err)

	md5 := new(hostMock.MD5Checker)
	md5.On("Parse", mock.Anything).Return("d41d8cd98f00b204e9800998ecf8427e", nil)

	h := host.New("host_1", toml, md5)
	runner := new(processMock.Runner)
	processes := new(processMock.ProcessesBuilder)
	p := process.New(runner, processes)

	exec := executor.New(
		executor.Register(h, p),
	)

	container, err := exec.FromHost("host_1")
	assert.NoError(t, err)
	assert.True(t, container.IsInstalled())
	assert.Equal(t, 3, container.PluginLength())

	_, err = container.Get("name_unknown", 1001, func(rpcType string, port int) caller.Caller {
		return mockCaller
	})

	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrPluginNotFound))
}

func TestGetPluginMetaSuccess(t *testing.T) {
	toml, err := discoverDriver.NewTomlParser().Parse([]byte(tomlContent))
	assert.NoError(t, err)

	md5 := new(hostMock.MD5Checker)
	md5.On("Parse", mock.Anything).Return("d41d8cd98f00b204e9800998ecf8427e", nil)

	h := host.New("host_1", toml, md5)
	runner := new(processMock.Runner)
	processes := new(processMock.ProcessesBuilder)
	p := process.New(runner, processes)

	exec := executor.New(
		executor.Register(h, p),
	)

	container, err := exec.FromHost("host_1")
	assert.NoError(t, err)
	assert.True(t, container.IsInstalled())
	assert.Equal(t, 3, container.PluginLength())

	meta, err := container.GetPluginMeta("name_1")
	assert.NoError(t, err)
	assert.Equal(t, "grpc", meta.ProtocolType)
}

func TestGetPluginMetaNotFound(t *testing.T) {
	toml, err := discoverDriver.NewTomlParser().Parse([]byte(tomlContent))
	assert.NoError(t, err)

	md5 := new(hostMock.MD5Checker)
	md5.On("Parse", mock.Anything).Return("d41d8cd98f00b204e9800998ecf8427e", nil)

	h := host.New("host_1", toml, md5)
	runner := new(processMock.Runner)
	processes := new(processMock.ProcessesBuilder)
	p := process.New(runner, processes)

	exec := executor.New(
		executor.Register(h, p),
	)

	container, err := exec.FromHost("host_1")
	assert.NoError(t, err)
	assert.True(t, container.IsInstalled())
	assert.Equal(t, 3, container.PluginLength())

	_, err = container.GetPluginMeta("name_unknown")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrPluginNotFound))
}

func TestAllowedProtocolFailed(t *testing.T) {
	mockCaller := new(callerMock.Caller)
	toml, err := discoverDriver.NewTomlParser().Parse([]byte(tomlContent))
	assert.NoError(t, err)

	md5 := new(hostMock.MD5Checker)
	md5.On("Parse", mock.Anything).Return("d41d8cd98f00b204e9800998ecf8427e", nil)

	h := host.New("host_1", toml, md5)
	runner := new(processMock.Runner)
	processes := new(processMock.ProcessesBuilder)
	p := process.New(runner, processes)

	exec := executor.New(
		executor.Register(h, p),
	)

	container, err := exec.FromHost("host_1")
	assert.NoError(t, err)
	assert.True(t, container.IsInstalled())
	assert.Equal(t, 3, container.PluginLength())

	_, err = container.Get("name_3", 1001, func(rpcType string, port int) caller.Caller {
		return mockCaller
	})

	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrProtocolUnknown))
}
