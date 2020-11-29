package caller_test

import (
	"errors"
	"testing"

	"github.com/quadroops/goplugin/pkg/errs"

	"github.com/quadroops/goplugin/pkg/caller"
	"github.com/quadroops/goplugin/pkg/caller/mocks"
	"github.com/quadroops/goplugin/pkg/executor"
	"github.com/quadroops/goplugin/pkg/host"
	"github.com/quadroops/goplugin/pkg/process"

	discoverDriver "github.com/quadroops/goplugin/pkg/discover/driver"
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
		rpc_type = "grpc"
		rpc_addr = "8080"
		
		[plugins.name_2]
		author = "author_2|author_2@gmail.com"
		md5 = "d41d8cd98f00b204e9800998ecf8427e"
		exec = "./tmp/test"
    	exec_file = "./tmp/test"
		exec_time = 10
		rpc_type = "rest"
		rpc_addr = "8080"
		
		[plugins.name_3]
		author = "author_3|author_3@gmail.com"
		md5 = "d41d8cd98f00b204e9800998ecf8427e"
		exec = "./tmp/test"
    	exec_file = "./tmp/test"
		exec_time = 20
		rpc_type = "nano"
		rpc_addr = "8080"

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

func TestPingSuccess(t *testing.T) {
	toml, err := discoverDriver.NewTomlParser().Parse([]byte(tomlContent))
	assert.NoError(t, err)

	md5 := new(hostMock.MD5Checker)
	md5.On("Parse", mock.Anything).Return("d41d8cd98f00b204e9800998ecf8427e", nil)

	h := host.New("host_1", toml, md5)
	assert.NotNil(t, h)

	runner := new(processMock.Runner)
	processes := new(processMock.ProcessesBuilder)
	p := process.New(runner, processes)

	exec := executor.New(
		&executor.Options{
			RetryTimeout: 3,
		},
		executor.Register(h, p),
	)

	container, err := exec.FromHost("host_1")
	assert.NoError(t, err)

	mockCaller := new(mocks.Caller)
	mockCaller.On("Ping").Once().Return("pong", nil)

	meta, err := container.GetPluginMeta("name_1")
	assert.NoError(t, err)

	plugin := caller.New(meta, mockCaller, 3)
	resp, err := plugin.Ping()
	assert.NoError(t, err)
	assert.Equal(t, "pong", resp)
}

func TestPingError(t *testing.T) {
	toml, err := discoverDriver.NewTomlParser().Parse([]byte(tomlContent))
	assert.NoError(t, err)

	md5 := new(hostMock.MD5Checker)
	md5.On("Parse", mock.Anything).Return("d41d8cd98f00b204e9800998ecf8427e", nil)

	h := host.New("host_1", toml, md5)
	assert.NotNil(t, h)

	runner := new(processMock.Runner)
	processes := new(processMock.ProcessesBuilder)
	p := process.New(runner, processes)

	exec := executor.New(
		&executor.Options{
			RetryTimeout: 3,
		},
		executor.Register(h, p),
	)

	container, err := exec.FromHost("host_1")
	assert.NoError(t, err)

	mockCaller := new(mocks.Caller)
	mockCaller.On("Ping").Once().Return("", errs.ErrPluginPing)

	meta, err := container.GetPluginMeta("name_1")
	assert.NoError(t, err)

	plugin := caller.New(meta, mockCaller, 3)
	_, err = plugin.Ping()
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrPluginPing))
}

func TestExecSuccess(t *testing.T) {
	toml, err := discoverDriver.NewTomlParser().Parse([]byte(tomlContent))
	assert.NoError(t, err)

	md5 := new(hostMock.MD5Checker)
	md5.On("Parse", mock.Anything).Return("d41d8cd98f00b204e9800998ecf8427e", nil)

	h := host.New("host_1", toml, md5)
	assert.NotNil(t, h)

	runner := new(processMock.Runner)
	processes := new(processMock.ProcessesBuilder)
	p := process.New(runner, processes)

	exec := executor.New(
		&executor.Options{
			RetryTimeout: 3,
		},
		executor.Register(h, p),
	)

	container, err := exec.FromHost("host_1")
	assert.NoError(t, err)

	mockCaller := new(mocks.Caller)
	mockCaller.On("Exec", "test.action", []byte("hello")).Once().Return([]byte("world"), nil)

	meta, err := container.GetPluginMeta("name_1")
	assert.NoError(t, err)

	plugin := caller.New(meta, mockCaller, 3)
	resp, err := plugin.Exec("test.action", []byte("hello"))
	assert.NoError(t, err)
	assert.Equal(t, []byte("world"), resp)
}

func TestExecError(t *testing.T) {
	toml, err := discoverDriver.NewTomlParser().Parse([]byte(tomlContent))
	assert.NoError(t, err)

	md5 := new(hostMock.MD5Checker)
	md5.On("Parse", mock.Anything).Return("d41d8cd98f00b204e9800998ecf8427e", nil)

	h := host.New("host_1", toml, md5)
	assert.NotNil(t, h)

	runner := new(processMock.Runner)
	processes := new(processMock.ProcessesBuilder)
	p := process.New(runner, processes)

	exec := executor.New(
		&executor.Options{
			RetryTimeout: 3,
		},
		executor.Register(h, p),
	)

	container, err := exec.FromHost("host_1")
	assert.NoError(t, err)

	mockCaller := new(mocks.Caller)
	mockCaller.On("Exec", "test.action", []byte("hello")).Once().Return(nil, errs.ErrPluginExec)

	meta, err := container.GetPluginMeta("name_1")
	assert.NoError(t, err)

	plugin := caller.New(meta, mockCaller, 3)
	_, err = plugin.Exec("test.action", []byte("hello"))
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrPluginExec))
}
