package main

import (
	"golang.org/x/sync/errgroup"
)

func main() {
	// appone 是插件app、  apptwo是主app
	eg := new(errgroup.Group)
	eg.Go(appOne)
	eg.Go(appTwo)
	eg.Wait()
}
