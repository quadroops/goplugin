package flow_test

import (
	"errors"
	"testing"

	"github.com/quadroops/goplugin/pkg/host/flow"
	"github.com/quadroops/goplugin/pkg/host/flow/mocks"
	"github.com/stretchr/testify/assert"
)

func TestInstallSuccess(t *testing.T) {
	md5Checker := new(mocks.MD5CheckerProxy)
	install := flow.NewInstall(md5Checker)

	plugin := flow.Plugin{
		Name: "test",
		Registry: flow.RegistryProxy{
			ExecFile: "./tmp/test",
		},
	}

	assertTrue := install.FilterByExecFile(plugin)
	assert.True(t, assertTrue)
}

func TestInstallDifferentCast(t *testing.T) {
	md5Checker := new(mocks.MD5CheckerProxy)
	install := flow.NewInstall(md5Checker)
	
	assertFalse := install.FilterByExecFile("wrong interface")
	assert.False(t, assertFalse)
}

func TestFilterMD5Success(t *testing.T) {
	md5Checker := new(mocks.MD5CheckerProxy)
	md5Checker.On("Parse", "./tmp/test").Once().Return("test", nil)

	install := flow.NewInstall(md5Checker)
	plugin := flow.Plugin{
		Name: "test",
		Registry: flow.RegistryProxy{
			ExecFile: "./tmp/test",
			MD5Sum: "test",
		},
	}
	
	assertTrue := install.FilterByMD5(plugin)
	assert.True(t, assertTrue)
}

func TestFilterMD5ErrorParse(t *testing.T) {
	md5Checker := new(mocks.MD5CheckerProxy)
	md5Checker.On("Parse", "./tmp/test").Once().Return("", errors.New("test"))

	install := flow.NewInstall(md5Checker)
	plugin := flow.Plugin{
		Name: "test",
		Registry: flow.RegistryProxy{
			ExecFile: "./tmp/test",
			MD5Sum: "test",
		},
	}
	
	assertFalse := install.FilterByMD5(plugin)
	assert.False(t, assertFalse)
}

func TestFilterMD5WrongInterface(t *testing.T) {
	md5Checker := new(mocks.MD5CheckerProxy)
	install := flow.NewInstall(md5Checker)
	
	assertFalse := install.FilterByMD5("wrong interface")
	assert.False(t, assertFalse)
}