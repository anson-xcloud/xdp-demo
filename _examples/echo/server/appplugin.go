package main

import (
	"context"
	"fmt"

	"github.com/anson-xcloud/xdp-demo/api"
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

	_, err := joinpoint.Join(context.Background(), &joinpoint.Config{
		ServerAddr: "appplugin:",
		Provider:   xcPlugin,
	}, joinpoint.WithLogger(c.Logger))
	return err
}

func echo(ctx context.Context, req *xcloud.Request) {
	echo := fmt.Sprintf("%s too, guys", req.GetBody())
	req.Response([]byte(echo))
	notify(ctx, req)
}

func echoServer(ctx context.Context, req *xcloud.Request) {
	echo := fmt.Sprintf("%s too, bots", req.GetBody())
	req.Response([]byte(echo))
}

func notify(ctx context.Context, req *xcloud.Request) {
	// if req.Source.Appid == "appmain" {
	// 	return
	// }

	xcPlugin.Post(ctx, &api.Peer{Appid: "appmain"}, "on_user_echo", nil, nil)
}
