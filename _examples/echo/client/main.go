package main

import (
	"time"

	"github.com/anson-xcloud/xdp-demo/client"
	"github.com/anson-xcloud/xdp-demo/pkg/logger"
)

func main() {
	const mainAppid = "appmain"
	const pluginAppid = "appplugin"

	if err := client.Connect(mainAppid); err != nil {
		logger.Error(err.Error())
		return
	}

	if err := client.Login("user", "pwd"); err != nil {
		logger.Error(err.Error())
		return
	}

	req := client.BuildRequest()
	req.Appid = pluginAppid
	req.Data = []byte("hello")
	t := time.NewTicker(time.Second * 3)
	for range t.C {
		st := time.Now()

		if data, err := client.Get(req); err != nil {
			client.GetLogger().Error("hello err:%s", err)
		} else {
			sec := time.Since(st).Seconds()
			client.GetLogger().Info("hello cost %.3f, msg: %s", sec, data)
		}
	}
}
