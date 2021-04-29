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

func hello(ctx context.Context, rw joinpoint.ResponseWriter, jr joinpoint.Request) {
	var resp Response
	resp.Time = time.Now()
	resp.Message = "hello world"
	data, _ := json.Marshal(&resp)
	rw.Write(data)
}

func main() {
	xcloud.SetEnv("dev")
	xcloud.HandleFunc(xcloud.HandlerRemoteAll, "", hello)

	if err := joinpoint.Join(context.Background(), &joinpoint.Config{
		Addr:     "appbasic:appkey",
		Provider: xcloud.Default,
	}); err != nil {
		return
	}

	time.Sleep(time.Hour * 24)
}
