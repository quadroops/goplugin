package main

import (
	"flag"
	"log"

	plugin "github.com/quadroops/goplugin/proto/plugin"
	"google.golang.org/grpc"

	"github.com/quadroops/goplugin"
)

var (
	addr = flag.String("addr", "", "Setup grpc plugin address")
	port = flag.Int("port", 8181, "Setup custom port")
	msg  = flag.String("msg", "hello world from grpc", "Setup custom msg for exec")
)

func init() {
	flag.Parse()
	if *addr == "" {
		panic("Need to define plugin grpc address")
	}
}

func main() {
	mainHost := goplugin.New("main")
	defer func() {
		mainHost.GetProcessInstance().KillAll()
	}()

	pluggable, err := goplugin.Register(mainHost)
	if err != nil {
		panic(err)
	}

	err = pluggable.Setup()
	if err != nil {
		panic(err)
	}

	plugins, err := pluggable.FromHost("main")
	if err != nil {
		panic(err)
	}

	err = plugins.Run("hello", *port)
	if err != nil {
		panic(err)
	}

	// used for driver.MakeClient
	client := func() (plugin.PluginClient, error) {
		conn, err := grpc.Dial(*addr, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}

		client := plugin.NewPluginClient(conn)
		return client, nil
	}

	pluginHello, err := plugins.Get("hello", *port, goplugin.BuildProtocol(&goplugin.ProtocolOption{
		GrpcClient: client,
	}))

	if err != nil {
		panic(err)
	}

	resp, err := pluginHello.Ping()
	if err != nil {
		panic(err)
	}

	log.Printf("Response ping: %s", resp)

	respExec, err := pluginHello.Exec("hello.msg", []byte(*msg))
	if err != nil {
		panic(err)
	}

	log.Printf("Response exec: %s", respExec)
}
