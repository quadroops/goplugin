package driver

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/quadroops/goplugin/pkg/errs"
	"github.com/quadroops/goplugin/pkg/caller"

	pbPlugin "github.com/quadroops/goplugin/proto/plugin"
	"google.golang.org/protobuf/types/known/emptypb"
)

// MakeClient used for building grpc connection without any options
// for now we only suppport for insecure connection
type MakeClient func() (pbPlugin.PluginClient, error)

type grpcProtocol struct {
	client MakeClient
}

// NewGRPC return new instance that implement Caller specifically
// for grpc's protocol
func NewGRPC(client MakeClient) caller.Caller {
	return &grpcProtocol{client}
}

func (g *grpcProtocol) Ping() (string, error) {
	client, err := g.client()
	if err != nil {
		return "", fmt.Errorf("%w: %q", errs.ErrProtocolGRPCConnection, err) 
	}

	resp, err := client.Ping(context.Background(), &emptypb.Empty{})
	if err != nil {
		return "", fmt.Errorf("%w: %q", errs.ErrProtocolGRPCResponse, err)
	}

	return resp.GetData().GetResponse(), nil
}

func (g *grpcProtocol) Exec(cmdName string, payload []byte) ([]byte, error) {
	client, err := g.client()
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
