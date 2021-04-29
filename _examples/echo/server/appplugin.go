package main

import (
	"context"
	"fmt"

	"github.com/anson-xcloud/xdp-demo/pkg/joinpoint"
	"github.com/anson-xcloud/xdp-demo/server"
	"github.com/anson-xcloud/xdp-demo/xcloud"
)

const appidPlugin = "appplugin"

func appPlugin() error {
	c := xcloud.DefaultConfig()
	c.Handler = xcloud.NewServeMux()
	c.Handler.HandleFunc(xcloud.HandlerRemoteAllUser, "", echo)
	c.Handler.HandleFunc(xcloud.HandlerRemoteAllServer, "", echoServer)
	xc := xcloud.New(c)

	return joinpoint.Join(context.Background(), &joinpoint.Config{
		Addr:     ":key1",
		Provider: xc,
	})
}

func echo(ctx context.Context, rw joinpoint.ResponseWriter, jr joinpoint.Request) {
	echo := fmt.Sprintf("%s too, guys", jr.Data.Data)
	rw.Write([]byte(echo))
	notify(svr, jr)
}

func echoServer(ctx context.Context, rw joinpoint.ResponseWriter, jr joinpoint.Request) {
	echo := fmt.Sprintf("%s too, bots", jr.Data.Data)
	rw.Write([]byte(echo))
}

func notify(rw joinpoint.ResponseWriter, jr joinpoint.Request) {
	if req.Appid == appidPlugin {
		return
	}

	svr.Send(&server.Remote{Appid: req.Appid}, &server.Data{Api: ""})
}
