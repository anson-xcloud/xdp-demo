package local

import (
	"context"
	"errors"
	"xcloud/pkg/clientapi/joinpoint"
)

type Provider struct {
	source chan joinpoint.Request
}

func New() *Provider {
	return &Provider{source: make(chan joinpoint.Request)}
}

func (p *Provider) Request(req *Request) error {
	p.source <- req
	return <-req.ch
}

func (p *Provider) Connect(ctx context.Context, addr string) (joinpoint.Point, []string, error) {
	return &Point{source: p.source}, nil, nil
}

func (p *Provider) Serve(ctx context.Context, rw joinpoint.ResponseWriter, jr joinpoint.Request) {
	rw.WriteStatus(joinpoint.StatusOK)
}

type Point struct {
	source chan joinpoint.Request
}

func (p *Point) Recv(ctx context.Context) (joinpoint.Request, error) {
	select {
	case req, ok := <-p.source:
		if !ok {
			return nil, errors.New("source close")
		}
		return req, nil
	case <-ctx.Done():
		return nil, errors.New("done")
	}
}
