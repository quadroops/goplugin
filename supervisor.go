package goplugin

import (
	"errors"
	"fmt"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/quadroops/goplugin/pkg/errs"
	"github.com/quadroops/goplugin/pkg/host"
	"github.com/quadroops/goplugin/pkg/supervisor"
)

const (
	defaultInterval = 5
)

// SupervisorOptionInterval setup interval
func SupervisorOptionInterval(interval int) PluginSupervisorOption {
	return func(s *PluginSupervisor) {
		s.interval = interval
	}
}

// SupervisorOptionCustomDriver used to setup custom driver, by default will
// use current PluginSupervisor as driver
func SupervisorOptionCustomDriver(drv supervisor.Driver) PluginSupervisorOption {
	return func(s *PluginSupervisor) {
		s.driver = drv
	}
}

// Supervisor used to supervisor all available plugins from all hosts
func Supervisor(pluggable *Registry, options ...PluginSupervisorOption) *PluginSupervisor {
	s := &PluginSupervisor{
		pluggable:  pluggable,
		interval:   defaultInterval,
		tickerDone: make(chan bool, 1),
	}

	for _, option := range options {
		option(s)
	}

	return s
}

// Setup build a registry for host and their plugin and their plugin's caller
// and creating new supervisor's instance
func (s *PluginSupervisor) Setup(handlers ...supervisor.OnErrorHandler) error {
	hostPlugins, err := s.pluggable.GetAllPlugins()
	if err != nil {
		return err
	}

	// adding internal error handlers
	handlers = append(handlers, s.AutoRestart)

	s.hostPlugins = hostPlugins
	s.handlers = append(s.handlers, handlers...)

	// make sure to check if driver has been set or not
	// if there are no custom driver has been set, than use
	// current object instance as default driver
	var drv supervisor.Driver
	if s.driver == nil {
		drv = s
	}

	s.runner = supervisor.New(drv, s.handlers...)
	return nil
}

// Start used to start main supervisor process
func (s *PluginSupervisor) Start() *PluginSupervisor {
	s.runner = s.runner.Start()
	return s
}

// Handle used start error handling
func (s *PluginSupervisor) Handle() {
	s.runner.Handle()
}

// Watch implement supervisor.Driver interface
func (s *PluginSupervisor) Watch() <-chan *supervisor.Payload {
	s.ticker = time.NewTicker(time.Duration(s.interval) * time.Second)
	payloadChan := make(chan *supervisor.Payload)

	// put the process in the background
	go func() {
		for {
			select {
			case <-s.tickerDone:
				log.Println("Supervisor stopped...")

				// stop all supervisor's processes
				s.ticker.Stop()

				// stopping infinite loop
				return
			case <-s.ticker.C:
				// get plugin's caller
				// and send the ping request
				// if something went wrong such as header response != 200
				// then trigger an error's event by put the payload variable
				// into payloadChan, this event should be catch
				// by all registered error handlers
				for _, hostPlugin := range s.hostPlugins {
					for plugin := range hostPlugin.Plugins {
						go func(plugin host.PluginName) {
							pluginName := string(plugin)

							// sending ping request
							err := s.checkPlugin(hostPlugin.Host, pluginName)

							// when an error triggered, we need to check if current error
							// allowed to send to channel
							if err != nil && errors.Is(err, errs.ErrEmptyProcesses) {
								payload := supervisor.Payload{
									Host:   hostPlugin.Host,
									Plugin: pluginName,
								}

								payloadChan <- &payload
							}
						}(plugin)
					}
				}
			}
		}
	}()

	return payloadChan
}

// OnError implement supervisor.Driver interface
func (s *PluginSupervisor) OnError(event *supervisor.Payload, handlers ...supervisor.OnErrorHandler) {
	if len(handlers) >= 1 {
		for _, handler := range handlers {
			handler(event)
		}
	}
}

// AutoRestart used to restart plugin's process if cannot be reached or something went wrong
func (s *PluginSupervisor) AutoRestart(payload *supervisor.Payload) {
	log.Println("Starting auto restart...")
	container, err := s.pluggable.GetContainer(payload.Host)
	if err != nil {
		log.Printf("Error getting container: %v", err)
		return
	}

	hostInstance, err := s.pluggable.GetHostPluginInstance(payload.Host)
	if err != nil {
		log.Printf("Error getting host: %v", err)
		return
	}

	pluginConf, err := hostInstance.GetPluginConf(payload.Plugin)
	if err != nil {
		log.Printf("Error getting plugin conf: %v", err)
		return
	}

	log.Printf("Killing plugin's process...")
	process := hostInstance.GetProcessInstance()
	err = process.Kill(payload.Plugin)
	if err != nil {
		log.Println("Killing plugin's process")
		return
	}

	meta, err := container.GetPluginMeta(payload.Plugin)
	if err != nil {
		log.Printf("Error getting meta: %v", err)
		return
	}

	var port int
	if meta.ProtocolType == "rest" {
		port = pluginConf.Protocol.RESTOpts.Port
	} else {
		port = pluginConf.Protocol.GRPCOpts.Port
	}

	log.Println("Restarting plugin's process")
	err = container.Run(payload.Plugin, port)
	if err != nil {
		log.Printf("Error restarting plugin's process: %v", err)
		return
	}
}

// Shutdown should be used on defer's way, it will should be automatically
// shutting down main supervisor's ticker
func (s *PluginSupervisor) Shutdown() {
	s.tickerDone <- true
}

func (s *PluginSupervisor) checkPlugin(host, plugin string) error {
	pid, err := s.pluggable.GetPID(host, plugin)
	if err != nil {
		return err
	}

	proc, err := os.FindProcess(int(pid))
	if err != nil {
		return err
	}

	log.Printf("Proc: %v", proc)
	if proc == nil {
		return errs.ErrEmptyProcesses
	}

	err = proc.Signal(syscall.Signal(0))
	if err != nil {
		log.Printf("Error signal: %v | ProcessID: %v", err, pid)
		return fmt.Errorf("%w: %q", errs.ErrEmptyProcesses, err)
	}

	return nil
}
