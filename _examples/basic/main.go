package main

import (
	"context"
	"encoding/json"
	"flag"
	"time"

	"github.com/anson-xcloud/xdp-demo/pkg/joinpoint"
	"github.com/anson-xcloud/xdp-demo/xcloud"
)

type Response struct {
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
}

func JSON(req *xcloud.Request, obj interface{}) {
	data, _ := json.Marshal(obj)
	req.Response(data)
}

func hello(ctx context.Context, req *xcloud.Request) {
	var resp Response
	resp.Time = time.Now()
	resp.Message = "hello world"
	JSON(req, &resp)
}

func main() {
	env := flag.String("env", "release", "xcloud env")
	flag.Parse()

	xcloud.SetEnv(*env)
	xcloud.HandleFunc(xcloud.HandlerRemoteAll, "", hello)

	ctx, cancel := context.WithTimeout(context.Background(), time.Hour*24)
	defer cancel()
	terminal, err := joinpoint.Join(ctx, &joinpoint.Config{
		ServerAddr: "appbasic:tokenbasic",
		Provider:   xcloud.Default(),
	})
	if err != nil {
		return
	}

	<-terminal.Context().Done()
}
