package joinpoint

import (
	"context"
)

type Transport interface {
	Recv(ctx context.Context) (Request, error)
}

type Request interface {
	GetResponseWriter() ResponseWriter
}

type ResponseWriter interface {
	Write(interface{})

	WriteStatus(st *Status)
}

type Provider interface {
	Connect(ctx context.Context, addr string) (Transport, []string, error)

	Serve(context.Context, ResponseWriter, Request)
}
