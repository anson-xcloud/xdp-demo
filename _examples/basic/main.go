package main

import (
	"github.com/anson-xcloud/xdp-demo/pkg/logger"
	"github.com/anson-xcloud/xdp-demo/server"
)

func hello(svr server.Server, req *server.Request) {
	svr.Reply(req, []byte("hello"))
}

func main() {
	server.SetEnv("debug")
	server.HandleFunc(server.HandlerRemoteAll, "", hello)

	addr := "appbasic:appkey"
	addr = "9410:9410"
	if err := server.Serve(addr,
		server.WithConfig(""),
	); err != nil {
		logger.Error("%s", err)
		return
	}
}
