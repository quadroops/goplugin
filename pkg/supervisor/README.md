# pkg/supervisor

Package: `github.com/quadroops/goplugin/pkg/supervisor`

**Overview**

This package provide an abstraction and functionality for watching plugin's
processes.  With current abstraction, we should be able to provide a mechanism to restart plugin's process or at least do something when plugin
cannot be reached.

## Types 

```go

// Payload used as main data when some plugin from some host indicated as error / cannot be reached
type Payload struct {
	Host   string
	Plugin string
}

// Driver used as main interface to run supervisor activities
type Driver interface {
	Watch() <-chan *Payload
	OnError(event *Payload, handlers ...OnErrorHandler)
}

// OnErrorHandler used as main type for handling plugin's error
type OnErrorHandler func(payload *Payload)

```

## Usages

```go
import (
	"github.com/quadroops/goplugin/pkg/supervisor"
)

handler1 := func(p *supervisor.Payload) {
    assert.NotEmpty(t, p.Host)
    assert.NotEmpty(t, p.Plugin)

    assert.Equal(t, p.Host, "test-host")
    assert.Equal(t, p.Plugin, "test-plugin")
}

handler2 := func(p *supervisor.Payload) {
    assert.Equal(t, p.Host, "test-host")
    assert.Equal(t, p.Plugin, "test-plugin")
}

// d is a driver implement `Driver` interface
// handler1 / handler2 is a function which implement `OnErrorHandler` signature
s := supervisor.New(d, handler1, handler2)

// run supervisor
s.Start().Handle()
```