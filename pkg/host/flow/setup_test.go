package flow_test

import (
	"context"
	"errors"
	"testing"

	"github.com/quadroops/goplugin/pkg/errs"
	"github.com/quadroops/goplugin/pkg/discover/driver"
	"github.com/quadroops/goplugin/pkg/host/flow"
	"github.com/reactivex/rxgo/v2"
	"github.com/stretchr/testify/assert"
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
		md5 = "d194b7bad208c2ddfa0ef597fd4abcc5"
		exec = "/path/to/exec"
    	exec_file = "/path/to/exec"
		exec_time = 10
		rpc_type = "rest"
		rpc_addr = "8080"
		
		[plugins.name_3]
		author = "author_3|author_3@gmail.com"
		md5 = "d194b7bad208c2ddfa0ef597fd4abcc5"
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

func TestFilterByHostnameSuccess(t *testing.T) {
	parser := driver.NewTomlParser()
	conf, err := parser.Parse([]byte(tomlContent))
	assert.NoError(t, err)

	f := flow.NewSetup("host_1", conf)
	exist := f.FilterByHostname(conf)
	assert.True(t, exist)
}

func TestFilterByHostnameInvalidType(t *testing.T) {
	parser := driver.NewTomlParser()
	conf, err := parser.Parse([]byte(tomlContent))
	assert.NoError(t, err)

	f := flow.NewSetup("host_1", conf)
	exist := f.FilterByHostname(1)
	assert.False(t, exist)
}

func TestFilterByHostnameNotExist(t *testing.T) {
	parser := driver.NewTomlParser()
	conf, err := parser.Parse([]byte(tomlContent))
	assert.NoError(t, err)

	f := flow.NewSetup("unknown", conf)
	exist := f.FilterByHostname(conf)
	assert.False(t, exist)
}

func TestFilterByPluginSuccess(t *testing.T) {
	parser := driver.NewTomlParser()
	conf, err := parser.Parse([]byte(tomlContent))
	assert.NoError(t, err)

	f := flow.NewSetup("host_3", conf)
	exist := f.FilterByPlugin("name_2")
	assert.True(t, exist)
}

func TestFilterByPluginInvalidType(t *testing.T) {
	parser := driver.NewTomlParser()
	conf, err := parser.Parse([]byte(tomlContent))
	assert.NoError(t, err)

	f := flow.NewSetup("host_3", conf)
	exist := f.FilterByPlugin(0)
	assert.False(t, exist)
}

func TestFilterByPluginUnknown(t *testing.T) {
	parser := driver.NewTomlParser()
	conf, err := parser.Parse([]byte(tomlContent))
	assert.NoError(t, err)

	f := flow.NewSetup("host_3", conf)
	exist := f.FilterByPlugin("unknown")
	assert.False(t, exist)
}

func TestMapToPluginsSuccess(t *testing.T) {
	parser := driver.NewTomlParser()
	conf, err := parser.Parse([]byte(tomlContent))
	assert.NoError(t, err)

	f := flow.NewSetup("host_1", conf)
	out, err := f.MapToPlugins(context.Background(), conf)
	assert.NoError(t, err)

	plugins, ok := out.([]string)
	assert.True(t, ok)
	assert.Len(t, plugins, 3)
}

func TestMapToPluginsUnknownHost(t *testing.T) {
	parser := driver.NewTomlParser()
	conf, err := parser.Parse([]byte(tomlContent))
	assert.NoError(t, err)

	f := flow.NewSetup("host_10", conf)
	_, err = f.MapToPlugins(context.Background(), conf)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrNoPlugins))
}

func TestMapToPluginsInvalidType(t *testing.T) {
	parser := driver.NewTomlParser()
	conf, err := parser.Parse([]byte(tomlContent))
	assert.NoError(t, err)

	f := flow.NewSetup("host_1", conf)
	_, err = f.MapToPlugins(context.Background(), 1)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrCastInterface))
}

func TestMapToPluginsEmpty(t *testing.T) {
	parser := driver.NewTomlParser()
	conf, err := parser.Parse([]byte(tomlContent))
	assert.NoError(t, err)

	f := flow.NewSetup("host_2", conf)
	out, err := f.MapToPlugins(context.Background(), conf)
	assert.NoError(t, err)

	plugins, ok := out.([]string)
	assert.True(t, ok)
	assert.Len(t, plugins, 0)
}

func TestFlatToNewObservableErrorCast(t *testing.T) {
	parser := driver.NewTomlParser()
	conf, err := parser.Parse([]byte(tomlContent))
	assert.NoError(t, err)

	f := flow.NewSetup("host_1", conf)
	observable := f.FlatToNewObservable(rxgo.Of(1))
	rxgo.Assert(context.Background(), t, observable, 
		rxgo.HasError(errs.ErrCastInterface),
	)
}

func TestFlatToNewObservableErrorEmtpyPlugins(t *testing.T) {
	parser := driver.NewTomlParser()
	conf, err := parser.Parse([]byte(tomlContent))
	assert.NoError(t, err)

	f := flow.NewSetup("host_1", conf)

	plugins := []string{}
	observable := f.FlatToNewObservable(rxgo.Of(plugins))
	rxgo.Assert(context.Background(), t, observable, 
		rxgo.HasError(errs.ErrNoPlugins),
	)
}

func TestFlatToNewObservableNoErrors(t *testing.T) {
	parser := driver.NewTomlParser()
	conf, err := parser.Parse([]byte(tomlContent))
	assert.NoError(t, err)

	f := flow.NewSetup("host_1", conf)

	plugins := []string{"test_1"}
	observable := f.FlatToNewObservable(rxgo.Of(plugins))
	rxgo.Assert(context.Background(), t, observable, 
		rxgo.HasNoError(),
		rxgo.IsNotEmpty(),
		rxgo.HasItems("test_1"),
	)
}