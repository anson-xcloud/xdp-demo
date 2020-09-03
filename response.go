package xdp

import "io"

type Response struct {
	req *Request

	Body io.Reader
}

type ResponseWriter interface {
}
