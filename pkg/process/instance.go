package process

import (
	"fmt"

	"github.com/quadroops/goplugin/pkg/errs"
)

// Instance used to setup new process instance
type Instance struct {
	runner    Runner
	processes ProcessesBuilder
}

// New used to create new process instance
func New(runner Runner, processes ProcessesBuilder) *Instance {
	return &Instance{
		runner:    runner,
		processes: processes,
	}
}

// Run used to start new subprocess
func (i *Instance) Run(toWait int, name, command string, args ...string) (<-chan Plugin, error) {
	if i.processes.IsExist(name) {
		return nil, fmt.Errorf("%w", errs.ErrPluginStarted)
	}

	return i.runner.Run(toWait, name, command, args...)
}

// Watch used to put current running plugin to list of rxgo.Item
func (i *Instance) Watch(p <-chan Plugin, err error) {
	if p != nil && err == nil {
		plugin := <-p
		i.processes.Add(plugin)
	}
}

// Kill used to kill individual plugin's process
func (i *Instance) Kill(name string) error {
	plugin, err := i.processes.Get(name)
	if err != nil {
		return err
	}

	if plugin.Name != "" {
		plugin.Kill()
	}

	i.processes.Remove(name)
	return nil
}

// KillAll used to kill all available plugin's processes
func (i *Instance) KillAll() []error {
	var errors []error
	observer, err := i.processes.Listen()
	if err != nil {
		errors = append(errors, err)
		return errors
	}

	<-observer.DoOnNext(func(val interface{}) {
		plugin, ok := val.(Plugin)
		if !ok {
			errors = append(errors, errs.ErrCastInterface)
		} else {
			plugin.Kill()
		}
	})

	// after kill all processes, we need to make sure
	// current processes is empty
	i.processes.Reset()
	return errors
}
