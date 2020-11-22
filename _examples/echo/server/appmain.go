package main

import (
	"time"

	"github.com/anson-xcloud/xdp-demo/server"
	"github.com/oklog/oklog/pkg/group"
)

var count int

// appMain has install plugin support by appPlugin
func appMain() error {
	appPluginServer := server.HandlerRemote{Type: server.HandlerRemoteTypeServer, Appid: appidPlugin}

	sm := server.NewServeMux()
	sm.HandleFunc(appPluginServer, "", onApp1Echo)

	svr := server.NewServer(server.WithHandler(sm))

	var gg group.Group
	gg.Add(func() error { return svr.Serve("appmain:") }, func(error) {})
	gg.Add(func() error {
		t := time.NewTicker(5 * time.Second)
		for range t.C {
			st := time.Now()
			if bdata, err := svr.Get(appidPlugin, &server.Data{Data: []byte("hello")}); err != nil {
				svr.GetLogger().Error("hello to plugin err:%s", err)
			} else {
				sec := time.Since(st).Seconds()
				svr.GetLogger().Info("hello plugin cost %.3f, msg: %s", sec, bdata)
			}
		}
		return nil
	}, func(error) {})
	return gg.Run()
}

func onApp1Echo(svr server.Server, req *server.Request) {
	count++
	svr.GetLogger().Debug("user total echo count: %d", count)
}
