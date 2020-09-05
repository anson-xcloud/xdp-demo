package xdp

import "fmt"

func ExampleServer() {
	svr := NewServer()
	svr.AppID = "1"
	svr.AppSecret = "test"

	if err := svr.Serve(); err != nil {
		fmt.Println(err)
	}

	// Output:
}

func ExampleClient() {
	svr := NewClient()
	svr.AppID = "1"

	// Output:
}
