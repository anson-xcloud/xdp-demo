package joinpoint

import (
	"context"
)

type Transport interface {
	Recv(ctx context.Context) (Request, error)

	Send(ctx context.Context, data interface{}) error

	Get(ctx context.Context, data interface{}) ([]byte, error)
}

type Request interface {
	String() string

	Response(data interface{})

	ResponseStatus(st *Status)
}

type Provider interface {
	Connect(ctx context.Context, addr string) (Transport, []string, error)

	Serve(context.Context, Request)
}
