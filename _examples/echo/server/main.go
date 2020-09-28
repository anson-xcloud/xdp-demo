package main

import (
	"fmt"
	"time"

	"github.com/anson-xcloud/xdp-demo"
)

func main() {
	xdp.HandleFunc("/", echo)

	svr := xdp.NewServer()
	if err := svr.Serve("1:test"); err != nil {
		svr.Logger().Error("%s", err)
	}
}

func echo(res xdp.ResponseWriter, req *xdp.Request) {
	echo := fmt.Sprintf("%s : %s %v",
		time.Now().Format(time.RFC3339),
		req.Api,
		req.Headers,
	)

	var sid = "TODO: nil session"
	if req.Session != nil {
		sid = req.Session.SessionID
	}
	fmt.Println("recv %s %s", sid, string(req.Body))
	res.Write([]byte(echo))
}
