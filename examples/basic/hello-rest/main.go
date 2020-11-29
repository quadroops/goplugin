package main

import (
	"flag"
	"log"

	"github.com/quadroops/goplugin"
	"github.com/quadroops/goplugin/pkg/caller/driver"
)

var (
	msg  = flag.String("msg", "world", "Setup message to exec")
	port = flag.Int("port", 8080, "Setup custom port")
)

func init() {
	flag.Parse()
}

func main() {
	hostAPluginHelloConf := goplugin.PluginConf{
		Protocol: &goplugin.ProtocolOption{
			RESTOpts: &driver.RESTOptions{
				Addr:    "localhost",
				Port:    *port,
				Timeout: 10, // in seconds
			},
		},
	}

	mainHost := goplugin.New("main").SetupPlugin(goplugin.Map("hello", &hostAPluginHelloConf))
	pluggable, err := goplugin.Register(mainHost).Install(&goplugin.InstallationOptions{
		RetryTimeoutCaller: 3,
	})
	if err != nil {
		panic(err)
	}

	defer func() {
		pluggable.KillPlugins()
	}()

	pluginHello, err := pluggable.GetCaller("main", "hello")
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
