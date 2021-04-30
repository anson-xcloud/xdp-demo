package main

import (
	"github.com/anson-xcloud/xdp-demo/xcloud"
	"golang.org/x/sync/errgroup"
)

func main() {
	xcloud.SetEnv("dev")

	eg := new(errgroup.Group)
	eg.Go(appMain)   // 主app
	eg.Go(appPlugin) // 插件app
	eg.Wait()
}
