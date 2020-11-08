package main

import "github.com/anson-xcloud/xdp-demo/server"

func plugin() error {
	svr := server.NewServer()
	return svr.Serve("app1:key1")
}
