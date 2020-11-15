package main

import (
	"fmt"

	"github.com/anson-xcloud/xdp-demo/server"
)

const appidPlugin = "appplugin"

func appPlugin() error {
	sm := server.NewServeMux()
	sm.HandleFunc(server.HandlerRemoteAllUser, "", echo)
	sm.HandleFunc(server.HandlerRemoteAllServer, "", echoServer)

	svr := server.NewServer(server.WithHandler(sm))
	return svr.Serve(appidPlugin + ":key1")
}

func echo(svr server.Server, req *server.Request) {
	echo := fmt.Sprintf("%s too, guys", req.Data.Data)
	svr.Reply(req, []byte(echo))
	notify(svr, req)
}

func echoServer(svr server.Server, req *server.Request) {
	echo := fmt.Sprintf("%s too, bots", req.Data.Data)
	svr.Reply(req, []byte(echo))
}

func notify(svr server.Server, req *server.Request) {
	if req.Appid == appidPlugin {
		return
	}

	svr.Send(&server.Remote{Appid: req.Appid}, &server.Data{Api: ""})
}
