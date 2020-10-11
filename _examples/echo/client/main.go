package main

import (
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

	cli.Login("user", "pwd")
	// for range time.NewTicker(time.Second).C {
	// 	data, err := xdp.Get(appid, "", headers)
	// 	fmt.Println("echo: \n", string(data), err)
	// }
}
