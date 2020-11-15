package main

import (
	"golang.org/x/sync/errgroup"
)

func main() {
	eg := new(errgroup.Group)
	eg.Go(appMain)   // 主app
	eg.Go(appPlugin) // 插件app
	eg.Wait()
}
