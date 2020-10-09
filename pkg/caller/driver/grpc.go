package driver

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/quadroops/goplugin/pkg/errs"
	pbPlugin "github.com/quadroops/goplugin/proto/plugin"
)

// GrpcClientConnector used to create client connection
type GrpcClientConnector func(addr string, port int) (pbPlugin.PluginClient, error)

// DefaultGrpcClientConnector used as default client object to create grpc
// connection
func DefaultGrpcClientConnector(addr string, port int) (pbPlugin.PluginClient, error) {
	endpoint := fmt.Sprintf("%s:%d", addr, port)
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := pbPlugin.NewPluginClient(conn)
	return client, nil
}

// GrpcOptions used to save grpc options
type GrpcOptions struct {
	Addr      string
	Port      int
	Connector GrpcClientConnector
}

// GrpcObj used as main grpc struct object
type GrpcObj struct {
	opt *GrpcOptions
}

// NewGRPC return new instance that implement Caller specifically
// for grpc's protocol
func NewGRPC(opt *GrpcOptions) *GrpcObj {
	// if caller not giving connector config
	// we're need to use default connector
	if opt.Connector == nil {
		opt.Connector = DefaultGrpcClientConnector
	}

	return &GrpcObj{opt}
}

// Ping implement caller.Caller ping method
func (g *GrpcObj) Ping() (string, error) {
	client, err := g.opt.Connector(g.opt.Addr, g.opt.Port)
	if err != nil {
		return "", fmt.Errorf("%w: %q", errs.ErrProtocolGRPCConnection, err)
	}

	resp, err := client.Ping(context.Background(), &emptypb.Empty{})
	if err != nil {
		return "", fmt.Errorf("%w: %q", errs.ErrPluginPing, err)
	}

	return resp.GetData().GetResponse(), nil
}

// Exec implement caller.Caller exec method
func (g *GrpcObj) Exec(cmdName string, payload []byte) ([]byte, error) {
	client, err := g.opt.Connector(g.opt.Addr, g.opt.Port)
	if err != nil {
		return nil, fmt.Errorf("%w: %q", errs.ErrProtocolGRPCConnection, err)
	}

	resp, err := client.Exec(context.Background(), &pbPlugin.ExecRequest{
		Command: cmdName,
		Payload: payload,
	})

	if err != nil {
		return nil, fmt.Errorf("%w: %q", errs.ErrPluginExec, err)
	}

	return resp.GetData().GetResponse(), nil
}
