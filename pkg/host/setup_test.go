package host_test

import (
	"errors"
	"testing"

	"github.com/quadroops/goplugin/pkg/discover/driver"
	"github.com/quadroops/goplugin/pkg/errs"
	"github.com/quadroops/goplugin/pkg/host"
	"github.com/quadroops/goplugin/pkg/host/mocks"
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
		protocol_type = "grpc"
		
		[plugins.name_2]
		author = "author_2|author_2@gmail.com"
		md5 = "d194b7bad208c2ddfa0ef597fd4abcc5"
		exec = "/path/to/exec"
    	exec_file = "/path/to/exec"
		exec_time = 10
		protocol_type = "grpc"
		
		[plugins.name_3]
		author = "author_3|author_3@gmail.com"
		md5 = "d194b7bad208c2ddfa0ef597fd4abcc5"
		exec = "./tmp/test"
    	exec_file = "./tmp/test"
		exec_time = 20
		protocol_type = "grpc"

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

func TestSetupSuccess(t *testing.T) {
	parser := driver.NewTomlParser()
	conf, err := parser.Parse([]byte(tomlContent))
	assert.NoError(t, err)
	assert.Len(t, conf.Plugins, 3)
	assert.Len(t, conf.Hosts, 3)

	md5 := new(mocks.MD5Checker)
	h := host.New("host_1", conf, md5)
	assert.NotNil(t, h)

	hmap := h.Setup()
	plugin, exist := hmap[host.PluginName("name_1")]
	assert.True(t, exist)
	assert.Equal(t, "./tmp/test", plugin.ExecPath)
	assert.Equal(t, "grpc", plugin.ProtocolType)
	assert.Equal(t, 5, plugin.ExecTime)
	assert.Equal(t, "d41d8cd98f00b204e9800998ecf8427e", plugin.MD5Sum)
}

func TestSetupEmptyConfig(t *testing.T) {
	md5 := new(mocks.MD5Checker)
	h := host.New("host_1", nil, md5)
	hmap := h.Setup()
	assert.Nil(t, hmap)
}

func TestSetupUnknownHostname(t *testing.T) {
	parser := driver.NewTomlParser()
	conf, err := parser.Parse([]byte(tomlContent))
	assert.NoError(t, err)
	assert.Len(t, conf.Plugins, 3)
	assert.Len(t, conf.Hosts, 3)

	md5 := new(mocks.MD5Checker)
	h := host.New("unknown", conf, md5)
	hmap := h.Setup()
	assert.Len(t, hmap, 0)
}

func TestSetupDoesntHavePlugins(t *testing.T) {
	parser := driver.NewTomlParser()
	conf, err := parser.Parse([]byte(tomlContent))
	assert.NoError(t, err)
	assert.Len(t, conf.Plugins, 3)
	assert.Len(t, conf.Hosts, 3)

	md5 := new(mocks.MD5Checker)
	h := host.New("host_2", conf, md5)
	hmap := h.Setup()
	assert.Len(t, hmap, 0)
}

func TestInstallSuccess(t *testing.T) {
	parser := driver.NewTomlParser()
	conf, err := parser.Parse([]byte(tomlContent))
	assert.NoError(t, err)
	assert.Len(t, conf.Plugins, 3)
	assert.Len(t, conf.Hosts, 3)

	md5Drv := new(mocks.MD5Checker)
	md5Drv.On("Parse", "./tmp/test").Return("d41d8cd98f00b204e9800998ecf8427e", nil)

	h := host.New("host_1", conf, md5Drv)
	hmap := h.Setup()
	assert.NotNil(t, hmap)

	hp, err := h.Install(hmap)
	assert.NoError(t, err)
	assert.NotNil(t, hp)

	plugins, exist := hp[host.Name("host_1")]
	assert.True(t, exist)
	assert.Len(t, plugins, 1)
	md5Drv.AssertCalled(t, "Parse", "./tmp/test")
}

func TestInstallEmptyPlugins(t *testing.T) {
	md5Drv := new(mocks.MD5Checker)
	h := host.New("host_1", nil, md5Drv)
	_, err := h.Install(nil)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrEmptyPlugins))
	md5Drv.AssertNotCalled(t, "Parse", mock.Anything)
}

func TestInstallNoPlugins(t *testing.T) {
	parser := driver.NewTomlParser()
	conf, err := parser.Parse([]byte(tomlContent))
	assert.NoError(t, err)
	assert.Len(t, conf.Plugins, 3)
	assert.Len(t, conf.Hosts, 3)

	md5Drv := new(mocks.MD5Checker)
	md5Drv.On("Parse", mock.Anything).Return("", errors.New("error"))

	h := host.New("host_3", conf, md5Drv)
	assert.NotNil(t, h)

	_, err = h.Install(h.Setup())
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrNoPlugins))
	md5Drv.AssertCalled(t, "Parse", "./tmp/test")
}
