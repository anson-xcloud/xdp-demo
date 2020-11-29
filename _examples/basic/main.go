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

	if err := server.Serve("appbasic:appkey",
		server.WithConfig(""),
	); err != nil {
		logger.Error("%s", err)
		return
	}
}
