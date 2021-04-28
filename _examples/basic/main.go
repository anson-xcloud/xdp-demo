package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/anson-xcloud/xdp-demo/pkg/joinpoint"
	"github.com/anson-xcloud/xdp-demo/xcloud"
)

type Response struct {
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
}

func hello(ctx context.Context, req *xcloud.Request) {
	var resp Response
	resp.Time = time.Now()
	resp.Message = "hello world"
	data, _ := json.Marshal(&resp)
	req.Write(data)
}

func main() {
	// server.SetEnv("dev")
	xcloud.HandleFunc(xcloud.HandlerRemoteAll, "", hello)

	// addr := "appbasic:appkey"
	// if err := server.Serve(addr,
	// 	server.WithConfig(""),
	// 	server.WithOnceTry(true),
	// 	// server.WithUID("1"),
	// ); err != nil {
	// 	logger.Error("%s", err)
	// 	return
	// }
	c := xcloud.DefaultConfig()
	c.Env = xcloud.Env(xcloud.EnvDebugDiscription)
	xc := xcloud.New(c)

	if err := joinpoint.Join(context.Background(), &joinpoint.Config{
		Addr:     "appbasic:appkey",
		Provider: xc,
	}); err != nil {
		return
	}

	time.Sleep(time.Hour * 24)
}
