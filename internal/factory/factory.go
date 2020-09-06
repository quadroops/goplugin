package factory

import (
	"github.com/quadroops/goplugin/pkg/host"
	driverHost "github.com/quadroops/goplugin/pkg/host/driver"

	"github.com/quadroops/goplugin/pkg/discover"
	driverDiscover "github.com/quadroops/goplugin/pkg/discover/driver"

	"github.com/quadroops/goplugin/pkg/process"
	driverProcess "github.com/quadroops/goplugin/pkg/process/driver"
)

// DefaultConfigChecker .
func DefaultConfigChecker() *discover.ConfigChecker {
	defaultChecker := driverDiscover.NewDefaultChecker()
	osChecker := driverDiscover.NewOsChecker()
	configChecker := discover.NewConfigChecker(osChecker, defaultChecker)
	return configChecker
}

// DefaultConfigParser .
func DefaultConfigParser() *discover.ConfigParser {
	tomlParser := driverDiscover.NewTomlParser()
	fileReader := driverDiscover.NewFileReader()
	return discover.NewConfigParser(tomlParser, fileReader)
}

// DefaultProcessInstance .
func DefaultProcessInstance() *process.Instance {
	subprocess := driverProcess.NewSubProcess()
	registry := driverProcess.NewRegistry()
	processes := driverProcess.NewProcesses(registry)
	return process.New(subprocess, processes)
}

// DefaultHostIdentityChecker .
func DefaultHostIdentityChecker() host.IdentityChecker {
	return driverHost.NewMd5Check()
}
