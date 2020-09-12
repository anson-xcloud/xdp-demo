package main

import (
	"fmt"
	"time"

	"github.com/anson-xcloud/xdp-demo"
)

func main() {
	svr := xdp.NewServer()
	svr.HTTPHandleFunc("/get", httpEcho)
	svr.HandleFunc("", tcpEcho)
	if err := svr.Serve("1:test"); err != nil {
		fmt.Println(err)
	}
	time.Sleep(time.Hour)
}

func httpEcho(res xdp.HTTPResponseWriter, req *xdp.HTTPRequest) {
	echo := fmt.Sprintf("echo %s", time.Now().Format(time.RFC3339))

	res.Write([]byte(echo))
}

func tcpEcho(req *xdp.Request) {
	fmt.Println("recv %s %s", req.Session.SessionID, string(req.Data))
	// if _, err := con.Write(data); err != nil {
	// 	fmt.Println("write ", err)
	// 	return
	// }
}
