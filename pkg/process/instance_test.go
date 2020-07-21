package process_test

import (
	"context"
	"errors"
	"testing"

	"github.com/quadroops/goplugin/pkg/errs"
	"github.com/quadroops/goplugin/pkg/process"
	"github.com/quadroops/goplugin/pkg/process/mocks"
	"github.com/reactivex/rxgo/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestRunSuccess(t *testing.T) {
	runner := new(mocks.Runner)
	runner.On("Run", 1, "test", "test", 1001).Once().Return(createMockChanPlugin(createMockPlugin("test")), nil)

	processes := new(mocks.ProcessesBuilder)
	processes.On("IsExist", "test").Once().Return(false)

	p := process.New(runner, processes)
	ch, err := p.Run(1, "test", "test", 1001)
	assert.NoError(t, err)

	plugin := <-ch
	assert.Equal(t, plugin.Name, "test")
}

func TestRunErrorExist(t *testing.T) {
	runner := new(mocks.Runner)
	processes := new(mocks.ProcessesBuilder)
	processes.On("IsExist", "test").Once().Return(true)

	p := process.New(runner, processes)
	_, err := p.Run(1, "test", "test", 1001)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrPluginStarted))
}

func TestWatchSuccess(t *testing.T) {
	plugin := createMockPlugin("test")
	pluginCh := createMockChanPlugin(plugin)

	runner := new(mocks.Runner)
	runner.On("Run", 1, "test", "test", 1001).Once().Return(pluginCh, nil)

	processes := new(mocks.ProcessesBuilder)
	processes.On("IsExist", "test").Once().Return(false)
	processes.On("Add", mock.Anything).Once().Return(nil)

	p := process.New(runner, processes)
	p.Watch(p.Run(1, "test", "test", 1001))
}

func TestKillSuccess(t *testing.T) {
	plugin := createMockChanPlugin(createMockPlugin("test"))
	pl := <-plugin

	runner := new(mocks.Runner)
	processes := new(mocks.ProcessesBuilder)
	processes.On("Get", "test").Once().Return(pl, nil)
	processes.On("Remove", "test").Once().Return(nil)

	p := process.New(runner, processes)

	err := p.Kill("test")
	assert.NoError(t, err)
}

func TestKillError(t *testing.T) {
	runner := new(mocks.Runner)
	processes := new(mocks.ProcessesBuilder)
	processes.On("Get", "test").Once().Return(process.Plugin{}, errs.ErrPluginNotFound)

	p := process.New(runner, processes)
	err := p.Kill("test")
	assert.Error(t, err)
}

func TestKillAllSuccess(t *testing.T) {
	plugin := createMockPlugin("test")
	obs := rxgo.Start([]rxgo.Supplier{func(_ context.Context) rxgo.Item {
		return rxgo.Of(plugin)
	}})

	runner := new(mocks.Runner)
	processes := new(mocks.ProcessesBuilder)
	processes.On("Listen").Once().Return(obs, nil)
	processes.On("Reset").Once().Return(nil)

	p := process.New(runner, processes)
	errs := p.KillAll()
	assert.Len(t, errs, 0)
}

func TestKillAllError(t *testing.T) {
	runner := new(mocks.Runner)
	processes := new(mocks.ProcessesBuilder)
	processes.On("Listen").Once().Return(nil, errors.New("error"))
	processes.On("Reset").Once().Return(nil)

	p := process.New(runner, processes)
	errs := p.KillAll()
	assert.Len(t, errs, 1)
}
