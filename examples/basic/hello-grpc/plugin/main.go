package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/golang/protobuf/ptypes/empty"
	plugin "github.com/quadroops/goplugin/proto/plugin"
	"google.golang.org/grpc"
)

var (
	commands = []string{"hello.world", "hello.msg"}
	port     = flag.String("port", "8080", "Set custom port")
)

type handler struct {
	plugin.UnimplementedPluginServer
}

func (h *handler) Ping(ctx context.Context, req *empty.Empty) (*plugin.PingResponse, error) {
	return &plugin.PingResponse{
		Status: "success",
		Data: &plugin.Data{
			Response: "pong",
		},
	}, nil
}

func (h *handler) Exec(ctx context.Context, req *plugin.ExecRequest) (*plugin.ExecResponse, error) {
	var command string
	for _, c := range commands {
		if c == req.GetCommand() {
			command = c
			break
		}
	}

	if command == "" {
		return nil, errors.New("Unknown command")
	}

	return &plugin.ExecResponse{
		Status: "success",
		Data: &plugin.Data{
			Response: req.GetPayload(),
		},
	}, nil
}

func init() {
	flag.Parse()
}

func main() {
	log.Printf("Running grpc servcer at port: %s", *port)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", *port))
	if err != nil {
		panic(err)
	}

	srv := grpc.NewServer()
	plugin.RegisterPluginServer(srv, &handler{})
	log.Fatalln(srv.Serve(lis))
}
