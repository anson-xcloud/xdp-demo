package main

import (
	"github.com/anson-xcloud/xdp-demo/server"
	"golang.org/x/sync/errgroup"
)

func main() {
	server.SetEnv("debug")

	eg := new(errgroup.Group)
	eg.Go(appMain)   // 主app
	eg.Go(appPlugin) // 插件app
	eg.Wait()
}
