package client

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

type Client interface {
	Serve(appid string) error

	Get(req *Request) ([]byte, error)

	// Post(req *ClientRequest) error

	// for test
	Login(user, pwd string) error
}

func Serve(appid string) error {
	return defaultClient.Serve(appid)
}

func Get(req *Request) ([]byte, error) {
	return defaultClient.Get(req)
}

func Login(user, pwd string) error {
	return defaultClient.Login(user, pwd)
}
