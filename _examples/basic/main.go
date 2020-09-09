package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/anson-xcloud/xdp-demo"
)

type Handler struct {
}

func (h *Handler) Serve(sess *xdp.Session, cmd uint32, data []byte) {

}

func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("hello http"))
}

func main() {
	svr := xdp.NewServer()
	svr.AppID = "1"
	svr.AppSecret = "test"
	svr.Handler = new(Handler)

	if err := svr.Serve(); err != nil {
		fmt.Println(err)
		return
	}

	if err := svr.Send(&xdp.Session{}, []byte("hello")); err != nil {
		fmt.Println(err)
	}
	time.Sleep(time.Hour)
}
