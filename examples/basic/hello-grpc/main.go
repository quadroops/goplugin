package main

import (
	"flag"
	"log"

	"github.com/quadroops/goplugin/pkg/caller/driver"

	"github.com/quadroops/goplugin"
)

var (
	addr = flag.String("addr", "localhost", "Setup grpc plugin address")
	port = flag.Int("port", 8181, "Setup custom port")
	msg  = flag.String("msg", "hello world from grpc", "Setup custom msg for exec")
)

func init() {
	flag.Parse()
}

func main() {
	hostAPluginHelloConf := goplugin.PluginConf{
		Protocol: &goplugin.ProtocolOption{
			GRPCOpts: &driver.GrpcOptions{
				Addr: *addr,
				Port: *port,
			},
		},
	}

	mainHost := goplugin.New("main").SetupPlugin(
		goplugin.Map("hello", &hostAPluginHelloConf),
	)

	pluggable, err := goplugin.Register(mainHost).Install()
	if err != nil {
		panic(err)
	}

	defer func() {
		pluggable.KillPlugins()
	}()

	pluginHello, err := pluggable.Get("main", "hello")
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
