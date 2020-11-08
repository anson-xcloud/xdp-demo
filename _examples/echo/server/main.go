package main

import (
	"golang.org/x/sync/errgroup"
)

func main() {
	eg := new(errgroup.Group)
	eg.Go(hostApp)
	eg.Go(pluginApp)
	eg.Wait()
}
