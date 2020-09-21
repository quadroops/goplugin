package driver

import (
	"context"
	"encoding/hex"
	"fmt"

	"google.golang.org/grpc"

	"github.com/quadroops/goplugin/pkg/errs"

	pbPlugin "github.com/quadroops/goplugin/proto/plugin"
	"google.golang.org/protobuf/types/known/emptypb"
)

// GrpcClientConnector used to create client connection
type GrpcClientConnector func(addr string) (pbPlugin.PluginClient, error)

// DefaultGrpcClientConnector used as default client object to create grpc
// connection
func DefaultGrpcClientConnector(addr string) (pbPlugin.PluginClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := pbPlugin.NewPluginClient(conn)
	return client, nil
}

// GrpcOptions used to save grpc options
type GrpcOptions struct {
	Addr      string
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
	client, err := g.opt.Connector(g.opt.Addr)
	if err != nil {
		return "", fmt.Errorf("%w: %q", errs.ErrProtocolGRPCConnection, err)
	}

	resp, err := client.Ping(context.Background(), &emptypb.Empty{})
	if err != nil {
		return "", fmt.Errorf("%w: %q", errs.ErrProtocolGRPCResponse, err)
	}

	return resp.GetData().GetResponse(), nil
}

// Exec implement caller.Caller exec method
func (g *GrpcObj) Exec(cmdName string, payload []byte) ([]byte, error) {
	client, err := g.opt.Connector(g.opt.Addr)
	if err != nil {
		return nil, fmt.Errorf("%w: %q", errs.ErrProtocolGRPCConnection, err)
	}

	resp, err := client.Exec(context.Background(), &pbPlugin.ExecRequest{
		Command: cmdName,
		Payload: hex.EncodeToString(payload),
	})

	if err != nil {
		return nil, fmt.Errorf("%w: %q", errs.ErrProtocolGRPCResponse, err)
	}

	return hex.DecodeString(resp.GetData().GetResponse())
}
