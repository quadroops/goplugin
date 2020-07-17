package process

import (
	"context"

	"github.com/quadroops/goplugin/internal/utils"
	"github.com/reactivex/rxgo/v2"
)

// ID is an alias for os PID
type ID int

// Plugin used when running a plugin to save their state and process id information
type Plugin struct {
	Kill   context.CancelFunc
	Name   string
	ID     ID
	Stdout *utils.Buffer
	Stderr *utils.Buffer
}

// ProcessesBuilder is main interface to manipulate list of available processes
type ProcessesBuilder interface {
	// Remove used to remove plugin from processes list
	Remove(string) error

	// Add used to adding new process
	Add(Plugin) error	

	// Reset should reset all processes 
	Reset()

	// IsExist used to check if plugin has a process or not
	IsExist(name string) bool

	// Listen used to start observe available streams, if any
	Listen() (rxgo.Observable, error) 

	// Get used to fetch the plugin info from list of processes or from
	// data registry
	Get(string) (Plugin, error) 
}

// RegistryBuilder used to as an interface to build registry of processes
// registry is like an indexing in database, used to prevent system to scanning
// all processes
type RegistryBuilder interface {
	Register(string, Plugin)
	IsExist(string) bool
	Delete(string)
	Get(string) (Plugin, error) 
	Reset()
}

// Runner used as main interface to start new subprocess
type Runner interface {
	Run(toWait int, name, execCommand string, args ...string) (<-chan Plugin, error)
}
