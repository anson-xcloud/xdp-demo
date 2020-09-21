package main

import (
	"github.com/anson-xcloud/xdp-demo"
)

func hello(res xdp.ResponseWriter, req *xdp.Request) {
	res.Write([]byte("hello"))
}

func main() {
	xdp.HandleFunc("", hello)

	svr := xdp.NewServer()
	if err := svr.Serve("1:test"); err != nil {
		svr.Logger().Error("%s", err)
		return
	}
}
