package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/anson-xcloud/xdp-demo/pkg/joinpoint"
	"github.com/anson-xcloud/xdp-demo/pkg/logger"
	"github.com/anson-xcloud/xdp-demo/server"
)

type Response struct {
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
}

func hello(svr server.Server, req *server.Request) {
	var resp Response
	resp.Time = time.Now()
	resp.Message = "hello world"
	data, _ := json.Marshal(&resp)
	svr.Reply(req, data)
}

func main() {
	server.SetEnv("dev")
	server.HandleFunc(server.HandlerRemoteAll, "", hello)

	addr := "appbasic:appkey"
	if err := server.Serve(addr,
		server.WithConfig(""),
		server.WithOnceTry(true),
		// server.WithUID("1"),
	); err != nil {
		logger.Error("%s", err)
		return
	}

	if err := joinpoint.Join(context.Background(), nil); err != nil {
		return
	}

	time.Sleep(time.Hour * 24)
}
