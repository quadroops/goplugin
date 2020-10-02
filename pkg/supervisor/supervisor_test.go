package supervisor_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/quadroops/goplugin/pkg/supervisor"
	"github.com/quadroops/goplugin/pkg/supervisor/mocks"
)

func makePayload(host, plugin string) *supervisor.Payload {
	return &supervisor.Payload{
		Host:   host,
		Plugin: plugin,
	}
}

func makeChan(payload *supervisor.Payload) <-chan *supervisor.Payload {
	c := make(chan *supervisor.Payload)
	go func() {
		c <- payload
	}()

	return c
}

func TestListenSuccess(t *testing.T) {
	payload := makePayload("test-host", "test-plugin")
	c := makeChan(payload)

	d := new(mocks.Driver)
	d.On("Watch").Once().Return(c)
	d.On("OnError", mock.Anything, mock.Anything, mock.Anything).Once().Run(func(args mock.Arguments) {
		for _, arg := range args {
			handler, ok := arg.(supervisor.OnErrorHandler)
			if ok {
				handler(payload)
			}
		}
	})

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

	s := supervisor.New(d, handler1, handler2)
	s.Start().Listen()
	time.Sleep(1 * time.Second)
}
