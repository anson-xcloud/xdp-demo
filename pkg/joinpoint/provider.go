package joinpoint

import (
	"context"
)

type Point interface {
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
	Connect(ctx context.Context, addr string) (Point, []string, error)

	Serve(context.Context, ResponseWriter, Request)
}
