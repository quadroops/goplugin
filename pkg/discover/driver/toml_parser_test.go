package driver_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/quadroops/goplugin/pkg/discover/driver"
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
		md5 = "d194b7bad208c2ddfa0ef597fd4abcc5"
		exec = "/path/to/exec"
		exec_time = 5
		rpc_type = "grpc"
		
		[plugins.name_2]
		author = "author_2|author_2@gmail.com"
		md5 = "d194b7bad208c2ddfa0ef597fd4abcc5"
		exec = "/path/to/exec"
		exec_time = 10
		rpc_type = "rest"
		
		[plugins.name_3]
		author = "author_3|author_3@gmail.com"
		md5 = "d194b7bad208c2ddfa0ef597fd4abcc5"
		exec = "/path/to/exec"
		exec_time = 20
		rpc_type = "nano"

	# Used as service registries
	# A service is an application that consume / using plugins
	# They are not allowed to access a plugin that not registered to service's plugin registries
	[hosts]

		[hosts.host_1]
		plugins = ["name_1", "name_2"]
		
		[hosts.host_2]
		plugins = ["name_3"]
	`

	tomlWrongContent = `
	title = "testing"
	[test]
	test1 = "value1"
	`
	
	tomlInvalidContent = `
	title = "testing
	[test]
	test1 = "value1"
	`
)

func TestParseSuccess(t *testing.T) {
	parser := driver.NewTomlParser()
	conf, err := parser.Parse([]byte(tomlContent))

	assert.NoError(t, err)
	assert.Equal(t, conf.Meta.Version, "1.0.0")
	assert.Equal(t, conf.Settings.Debug, true)
	assert.Len(t, conf.Plugins, 3)
	assert.Len(t, conf.Hosts, 2)

	_, plugin1Exist := conf.Plugins["name_1"]
	assert.True(t, plugin1Exist)	
	
	_, plugin2Exist := conf.Plugins["name_2"]
	assert.True(t, plugin2Exist)	
	
	_, plugin3Exist := conf.Plugins["name_3"]
	assert.True(t, plugin3Exist)	

	_, host1Exist := conf.Hosts["host_1"]
	assert.True(t, host1Exist)
	
	_, host2Exist := conf.Hosts["host_2"]
	assert.True(t, host2Exist)
}

func TestParseInvalidContent(t *testing.T) {
	parser := driver.NewTomlParser()
	conf, err := parser.Parse([]byte(tomlWrongContent))
	assert.NoError(t, err)
	assert.Empty(t, conf.Meta.Author)
}

func TestParseInvalidToml(t *testing.T) {
	parser := driver.NewTomlParser()
	_, err := parser.Parse([]byte(tomlInvalidContent))
	assert.Error(t, err)
}