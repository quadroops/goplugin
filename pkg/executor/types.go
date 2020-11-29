package executor

import (
	"github.com/quadroops/goplugin/pkg/host"
	"github.com/quadroops/goplugin/pkg/process"
	"github.com/reactivex/rxgo/v2"
)

// Registry used when register new host and their isolated process
type Registry struct {
	Host    *host.Builder
	Process *process.Instance
}

// Exec used main object of executor
type Exec struct {
	processes []rxgo.Supplier
	options   *Options
}

// Options for now used only to store retry timeout which will
// be used by Container
type Options struct {
	RetryTimeout int
}

// Container used as main object to start host's processes
type Container struct {
	installed    bool
	retryTimeout int
	Registry     *Registry
	plugins      host.Plugins
}
