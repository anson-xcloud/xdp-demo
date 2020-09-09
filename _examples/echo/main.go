package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/anson-xcloud/xdp-demo"
)

func main() {
	go httpServe()
	go xdpServe()
	time.Sleep(time.Hour)
}

func handleEcho(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("echo"))
}

func httpServe() {
	http.HandleFunc("/", handleEcho)
	if err := http.ListenAndServe(":8000", nil); err != nil {
		fmt.Println(err)
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
