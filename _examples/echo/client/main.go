package main

import (
	"fmt"
	"net/url"
	"time"

	"github.com/anson-xcloud/xdp-demo"
)

func main() {
	const appid = "1"

	headers := make(url.Values)
	headers.Set("msg", "hello")

	cli, _ := xdp.Login(appid, "user", "pwd")

	time.Sleep(time.Second)

	for range time.NewTicker(time.Second * 3).C {
		data, err := cli.Get("", headers)
		fmt.Println("echo: ", string(data), err)
	}
}
