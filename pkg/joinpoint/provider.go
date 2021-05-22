package joinpoint

import (
	"context"
)

type Transport interface {
	Recv(ctx context.Context) (Request, error)
}

type Request interface {
	String() string

	Response(interface{})

	ResponseStatus(st *Status)
}

type Provider interface {
	Connect(ctx context.Context, addr string) (Transport, []string, error)

	Serve(context.Context, Request)
}
