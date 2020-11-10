package main

import (
	"fmt"
	"time"

	"github.com/anson-xcloud/xdp-demo/server"
)

const appidOne = "app1"

func appOne() error {
	sm := server.NewServeMux()
	sm.HandleFunc(server.HandlerRemoteAllUser, "", echo)

	svr := server.NewServer(server.WithHandler(sm))
	return svr.Serve(appidOne + ":key1")
}

func echo(svr server.Server, req *server.Request) {
	echo := fmt.Sprintf("%s : %s %v",
		time.Now().Format(time.RFC3339),
		req.Api,
		req.Headers,
	)

	fmt.Printf("recv %s %s %v\n", req.Sid, string(req.Data.Data), req.Headers)
	svr.Reply(req, []byte(echo))
	notify(svr, req)
}

func notify(svr server.Server, req *server.Request) {
	if req.Appid == appidOne {
		return
	}

	svr.Send(&server.Remote{Appid: req.Appid}, &server.Data{Api: ""})
}
