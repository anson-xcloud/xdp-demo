package main

import (
	"fmt"
	"time"

	"github.com/anson-xcloud/xdp-demo"
)

func main() {
	xdp.HTTPHandleFunc("/", httpEcho)
	xdp.HandleFunc("", tcpEcho)

	svr := xdp.NewServer()
	if err := svr.Serve("1:test"); err != nil {
		fmt.Println(err)
	}
	time.Sleep(time.Hour)
}

func httpEcho(res xdp.HTTPResponseWriter, req *xdp.HTTPRequest) {
	echo := fmt.Sprintf("%s : %s %v",
		time.Now().Format(time.RFC3339),
		req.Path,
		req.Headers,
	)

	res.Write([]byte(echo))
}

func tcpEcho(req *xdp.Request) {
	fmt.Println("recv %s %s", req.Session.SessionID, string(req.Data))
	// if _, err := con.Write(data); err != nil {
	// 	fmt.Println("write ", err)
	// 	return
	// }
}
