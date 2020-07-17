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
}

// Container used as main object to start host's processes
type Container struct {
	Registry  *Registry
	installed bool
	plugins   host.Plugins
}
