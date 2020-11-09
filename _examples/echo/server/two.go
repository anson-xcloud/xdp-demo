package main

import "github.com/anson-xcloud/xdp-demo/server"

var count int

// apptwo has install plugin support by appone
func appTwo() error {
	appOneServer := server.HandlerSource{Type: server.HandlerSourceTypeServer, Appid: appidOne}

	sm := server.NewServeMux()
	sm.HandleFunc(appOneServer, "", onApp1Echo)

	svr := server.NewServer(server.WithHandler(sm))
	return svr.Serve("app2:key2")
}

func onApp1Echo(svr server.Server, req *server.Request) {
	count++
	svr.GetLogger().Debug("user total echo count: %d", count)
}
