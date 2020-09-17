package main

import (
	"flag"
	"log"

	"github.com/quadroops/goplugin"
)

var (
	msg = flag.String("msg", "world", "Setup message to exec")
)

func init() {
	flag.Parse()
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

	respExec, err := pluginHello.Exec("hello.msg", []byte(*msg))
	if err != nil {
		panic(err)
	}

	log.Printf("Response exec: %s", respExec)
}
