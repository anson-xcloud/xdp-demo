package main

import (
	"golang.org/x/sync/errgroup"
)

func main() {
	eg := new(errgroup.Group)
	eg.Go(appOne)
	eg.Go(appTwo)
	eg.Wait()
}
