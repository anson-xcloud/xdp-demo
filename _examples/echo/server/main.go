package main

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/anson-xcloud/xdp-demo"
)

func main() {
	go httpServe()
	go tcpServe()
	go xdpServe()
	time.Sleep(time.Hour)
}

func httpEcho(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("echo"))
}

func tcpEcho(con net.Conn, data []byte) {
	if _, err := con.Write(buf[:n]); err != nil {
		fmt.Println("write ", err)
		return
	}
}

func httpServe() {
	http.HandleFunc("/", httpEcho)
	if err := http.ListenAndServe(":8000", nil); err != nil {
		fmt.Println(err)
	}
}

func tcpServe() {
	lis, err := net.Listen("tcp", ":8001")
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		con, err := lis.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		go func() {
			var buf [1024]byte
			n, err := con.Read(buf[:0])
			if err != nil {
				fmt.Println("read ", err)
				return
			}

			tcpEcho(con, buf[:n])
		}()
	}
}

func xdpServe() {
	sv := xdp.NewServer()
	svr.AppID = "1"
	svr.AppSecret = "test"
	// svr.Handler = new(Handler)
	if err := sv.Serve(); err != nil {
		fmt.Println(err)
	}
}
