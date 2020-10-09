package main

import (
	"flag"
	"log"
	"time"

	"github.com/quadroops/goplugin"
	"github.com/quadroops/goplugin/pkg/caller/driver"
)

var (
	msg      = flag.String("msg", "world", "Setup message to exec")
	port     = flag.Int("port", 8080, "Setup custom port")
	interval = flag.Int("interval", 2, "Setup custom interval")
	stopper  = flag.Int("stopper", 30, "Setup custom interval")
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
	pluggable, err := goplugin.Register(mainHost).Install()
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

	pluginProcess, err := pluggable.GetPID("main", "hello")
	if err != nil {
		log.Printf("Error get PID: %v", err)
		return
	}
	log.Printf("Hello processID: %d", int(pluginProcess))

	ticker := time.NewTicker(time.Duration(*interval) * time.Second)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				log.Println("Process done...")
				ticker.Stop()
				return
			case t := <-ticker.C:
				log.Println("=======================================")
				resp, err := pluginHello.Ping()
				if err != nil {
					log.Printf("Error ping: %v", err)
					panic(err)
				}
				log.Printf("Response ping: %s, at: %v", resp, t)

				respExec, err := pluginHello.Exec("hello.msg", []byte(*msg))
				if err != nil {
					log.Printf("Error exec: %v", err)
					panic(err)
				}

				log.Printf("Response exec: %s, at: %v", respExec, t)
			}
		}
	}()

	// start supervisor
	watcher := goplugin.Supervisor(pluggable)
	err = watcher.Setup()
	if err != nil {
		panic(err)
	}
	defer watcher.Shutdown()
	watcher.Start().Handle()

	time.Sleep(time.Duration(*stopper) * time.Second)
	done <- true
}
