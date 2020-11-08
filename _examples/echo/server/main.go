package main

import (
	"golang.org/x/sync/errgroup"
)

func main() {
	eg := new(errgroup.Group)
	eg.Go(app)
	eg.Go(plugin)
	eg.Wait()
}
