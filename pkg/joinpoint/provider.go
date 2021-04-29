package joinpoint

import (
	"context"
)

type Transport interface {
	Recv(ctx context.Context) (Request, error)
}

type Request interface {
	Discription() string

	GetResponseWriter() ResponseWriter
}

type ResponseWriter interface {
	Write(interface{})

	WriteStatus(st *Status)
}

type Handler interface {
	Serve(context.Context, ResponseWriter, Request)
}

type Provider interface {
	Handler

	Connect(ctx context.Context, addr string) (Transport, []string, error)
}

type HandlerFunc func(context.Context, ResponseWriter, Request)

func (h HandlerFunc) Serve(ctx context.Context, rw ResponseWriter, req Request) {
	h(ctx, rw, req)
}
