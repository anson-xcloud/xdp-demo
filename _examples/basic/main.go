package main

import (
	"github.com/anson-xcloud/xdp-demo/pkg/logger"
	"github.com/anson-xcloud/xdp-demo/server"
)

func hello(svr server.Server, req *server.Request) {
	svr.Reply(req, []byte("hello"))
}

func main() {
	server.SetEnv("dev")
	server.HandleFunc(server.HandlerRemoteAll, "", hello)

	addr := "appbasic:appkey"
	if err := server.Serve(addr,
		server.WithConfig(""),
		server.WithOnceTry(true),
		// server.WithUID("1"),
	); err != nil {
		logger.Error("%s", err)
		return
	}
}
