package client

import "github.com/anson-xcloud/xdp-demo/pkg/logger"

var defaultClient Client

func init() {
	defaultClient = NewHttpClient()
}

type Request struct {
	Api     string
	Appid   string
	Headers map[string]string
	Data    []byte
}

func BuildRequest() *Request {
	return &Request{Headers: make(map[string]string)}
}

func GetLogger() logger.Logger {
	return logger.Default
}

type Client interface {
	Connect(appid string) error

	Get(req *Request) ([]byte, error)

	// Post(req *ClientRequest) error

	// for test
	Login(user, pwd string) error
}

func Connect(appid string) error {
	return defaultClient.Connect(appid)
}

func Get(req *Request) ([]byte, error) {
	return defaultClient.Get(req)
}

func Login(user, pwd string) error {
	return defaultClient.Login(user, pwd)
}
