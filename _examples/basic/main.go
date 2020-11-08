package main

import (
	"github.com/anson-xcloud/xdp-demo/server"
)

func hello(svr server.Server, req *server.Request) {
	svr.Reply(req, []byte("hello"))
}

func main() {
	server.HandleFunc("", hello)

	if err := server.Serve("appid:appkey"); err != nil {
		server.Error("%s", err)
		return
	}
}
