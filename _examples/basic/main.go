package main

import (
	"fmt"
	"time"

	"github.com/anson-xcloud/xdp-demo"
)

func main() {
	svr := xdp.NewServer()
	svr.AppID = "1"
	svr.AppSecret = "test"

	if err := svr.Serve(); err != nil {
		fmt.Println(err)
	}

	if err := svr.Send(&xdp.Session{}, []byte("hello")); err != nil {
		fmt.Println(err)
	}
	time.Sleep(time.Hour)
}
