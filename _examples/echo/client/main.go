package main

import (
	"fmt"
	"time"

	"github.com/anson-xcloud/xdp-demo/client"
	"github.com/anson-xcloud/xdp-demo/pkg/logger"
)

func main() {
	const appid = "app2"

	if err := client.Serve(appid); err != nil {
		logger.Error(err.Error())
		return
	}

	if err := client.Login("user", "pwd"); err != nil {
		logger.Error(err.Error())
		return
	}

	req := client.BuildRequest()
	req.Appid = "app1"
	req.Data = []byte("hello")
	for range time.NewTicker(time.Second * 3).C {
		data, err := client.Get(req)
		fmt.Println("echo: ", string(data), err)
	}
}
