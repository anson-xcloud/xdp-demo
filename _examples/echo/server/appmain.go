package main

import (
	"context"
	"time"

	"github.com/anson-xcloud/xdp-demo/pkg/joinpoint"
	"github.com/anson-xcloud/xdp-demo/pkg/xlog"
	"github.com/anson-xcloud/xdp-demo/xcloud"
)

var count int
var xcMainLogger xlog.Logger

// appMain has install plugin support by appPlugin
func appMain() error {
	appPluginServer := xcloud.HandlerRemote{Type: xcloud.HandlerRemoteTypeServer, Appid: appidPlugin}
	xcMainLogger = xlog.Default.With("app", "main")

	c := xcloud.DefaultConfig()
	c.Logger = xcMainLogger
	c.Handler = xcloud.NewServeMux()
	c.Handler.HandleFunc(appPluginServer, "on_user_echo", onPluginUserEcho)
	xc, _ := xcloud.New(c)

	if err := joinpoint.Join(context.Background(), &joinpoint.Config{
		ServerAddr: "appmain:",
		Provider:   xc,
	}, joinpoint.WithLogger(xcMainLogger)); err != nil {
		return err
	}

	t := time.NewTicker(5 * time.Second)
	for range t.C {
		st := time.Now()
		if bdata, err := xc.Get(context.Background(), appidPlugin, "echo", nil, []byte("hello")); err != nil {
			xcMainLogger.Errorf("hello to plugin err:%s", err)
		} else {
			sec := time.Since(st).Seconds()
			xcMainLogger.Infof("hello plugin cost %.3f, msg: %s", sec, bdata)
		}
	}
	return nil
}

func onPluginUserEcho(ctx context.Context, req *xcloud.Request) {
	count++
	xcMainLogger.Debugf("user total echo count: %d", count)
}
