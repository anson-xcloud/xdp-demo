package main

import (
	"context"
	"fmt"

	"github.com/anson-xcloud/xdp-demo/pkg/joinpoint"
	"github.com/anson-xcloud/xdp-demo/pkg/xlog"
	"github.com/anson-xcloud/xdp-demo/xcloud"
)

const appidPlugin = "appplugin"

var xcPlugin *xcloud.XCloud

func appPlugin() error {
	c := xcloud.DefaultConfig()
	c.Handler = xcloud.NewServeMux()
	c.Logger = xlog.Default.With("app", "plugin")
	c.Handler.HandleFunc(xcloud.HandlerRemoteAllUser, "echo", echo)
	c.Handler.HandleFunc(xcloud.HandlerRemoteAllServer, "echo", echoServer)
	xcPlugin, _ = xcloud.New(c)

	return joinpoint.Join(context.Background(), &joinpoint.Config{
		ServerAddr: "appplugin:",
		Provider:   xcPlugin,
	}, joinpoint.WithLogger(c.Logger))
}

func echo(ctx context.Context, rw joinpoint.ResponseWriter, jr joinpoint.Request) {
	req := jr.(*xcloud.Request)

	echo := fmt.Sprintf("%s too, guys", req.Data.Data)
	rw.Write([]byte(echo))
	notify(ctx, req)
}

func echoServer(ctx context.Context, rw joinpoint.ResponseWriter, jr joinpoint.Request) {
	req := jr.(*xcloud.Request)

	echo := fmt.Sprintf("%s too, bots", req.Data.Data)
	rw.Write([]byte(echo))
}

func notify(ctx context.Context, req *xcloud.Request) {
	if req.Appid == "" || req.Appid == appidPlugin {
		return
	}

	xcPlugin.Post(ctx, &xcloud.Remote{Appid: req.Appid}, &xcloud.Data{Api: ""})
}
