package main

import "github.com/anson-xcloud/xdp-demo/server"

var count int

func pluginApp() error {
	sm := server.NewServeMux()
	sm.HandleFunc("/", echo)

	svr := server.NewServer(server.WithHandler(sm))
	return svr.Serve("app2:key2")
}

func onEchoUser(svr server.Server, req *server.Request) {
	count++
	svr.GetLogger().Debug("user total echo count: %d", count)
}
