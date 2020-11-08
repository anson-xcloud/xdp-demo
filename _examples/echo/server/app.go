package main

import (
	"fmt"
	"time"

	"github.com/anson-xcloud/xdp-demo/server"
)

func app() error {
	server.HandleFunc("/", echo)

	svr := server.NewServer()
	return svr.Serve("app1:key1")
}

func echo(svr server.Server, req *server.Request) {
	echo := fmt.Sprintf("%s : %s %v",
		time.Now().Format(time.RFC3339),
		req.Api,
		req.Headers,
	)

	fmt.Printf("recv %s %s %v\n", req.Sid, string(req.Data.Data), req.Headers)
	svr.Reply(req, []byte(echo))
}
