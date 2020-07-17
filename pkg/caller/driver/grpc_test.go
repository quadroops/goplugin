package driver_test

import (
	"context"
	"encoding/hex"
	"errors"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/quadroops/goplugin/pkg/errs"
	"github.com/quadroops/goplugin/pkg/caller/driver"
	pbPlugin "github.com/quadroops/goplugin/proto/plugin"
	"github.com/quadroops/goplugin/proto/plugin/mocks"
	"github.com/stretchr/testify/assert"
)

func TestPingSuccess(t *testing.T) {
	client := new(mocks.PluginClient)
	client.On("Ping", context.Background(), &empty.Empty{}).Once().Return(&pbPlugin.PingResponse{
		Status: "success",
		Data: &pbPlugin.Data{
			Response: "pong",
		},
	}, nil)

	rpc := driver.NewGRPC(func() (pbPlugin.PluginClient, error) {
		return client, nil
	})

	resp, err := rpc.Ping()
	assert.NoError(t, err)
	assert.Equal(t, "pong", resp)
}

func TestPingErrorClient(t *testing.T) {
	rpc := driver.NewGRPC(func() (pbPlugin.PluginClient, error) {
		return nil, errors.New("error conn") 
	})

	resp, err := rpc.Ping()
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrProtocolGRPCConnection))
	assert.Empty(t, resp)
}

func TestPingErrorResponse(t *testing.T) {
	client := new(mocks.PluginClient)
	client.On("Ping", context.Background(), &empty.Empty{}).Once().Return(&pbPlugin.PingResponse{}, errors.New("err response"))

	rpc := driver.NewGRPC(func() (pbPlugin.PluginClient, error) {
		return client, nil
	})

	resp, err := rpc.Ping()
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrProtocolGRPCResponse))
	assert.Empty(t, resp)
}

func TestGRPCExecSuccess(t *testing.T) {
	client := new(mocks.PluginClient)
	client.On("Exec", context.Background(), &pbPlugin.ExecRequest{
		Command: "test.command",
		Payload: hex.EncodeToString([]byte("hello")),
	}).Once().Return(&pbPlugin.ExecResponse{
		Status: "success",
		Data: &pbPlugin.Data{
			Response: hex.EncodeToString([]byte("world")),
		},
	}, nil)
	
	rpc := driver.NewGRPC(func() (pbPlugin.PluginClient, error) {
		return client, nil
	})

	resp, err := rpc.Exec("test.command", []byte("hello"))
	assert.NoError(t, err)
	assert.Equal(t, []byte("world"), resp)
}

func TestGRPCErrorClient(t *testing.T) {	
	rpc := driver.NewGRPC(func() (pbPlugin.PluginClient, error) {
		return nil, errors.New("error conn") 
	})

	resp, err := rpc.Exec("test.command", []byte("hello"))
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrProtocolGRPCConnection))
	assert.Empty(t, resp)
}

func TestGRPCErrorResponse(t *testing.T) {
	client := new(mocks.PluginClient)
	client.On("Exec", context.Background(), &pbPlugin.ExecRequest{
		Command: "test.command",
		Payload: hex.EncodeToString([]byte("hello")),
	}).Once().Return(&pbPlugin.ExecResponse{}, errors.New("error response"))

	rpc := driver.NewGRPC(func() (pbPlugin.PluginClient, error) {
		return client, nil
	})

	resp, err := rpc.Exec("test.command", []byte("hello"))
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrProtocolGRPCResponse))
	assert.Empty(t, resp)
}