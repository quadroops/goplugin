package main

import (
	"log"

	"github.com/quadroops/goplugin"
)

func main() {
	mainHost := goplugin.New("main")
	defer func() {
		mainHost.GetProcessInstance().KillAll()
	}()

	pluggable := goplugin.Register(mainHost)
	err := pluggable.Setup()
	if err != nil {
		panic(err)
	}

	plugins, err := pluggable.FromHost("main")
	if err != nil {
		panic(err)
	}

	err = plugins.Run("hello", 8081)
	if err != nil {
		panic(err)
	}

	pluginHello, err := plugins.Get("hello", 8081, goplugin.BuildProtocol(&goplugin.ProtocolOption{
		RestAddr: "http://localhost",
	}))
	if err != nil {
		panic(err)
	}

	resp, err := pluginHello.Ping()
	if err != nil {
		panic(err)
	}

	log.Printf("Response ping: %s", resp)
}
