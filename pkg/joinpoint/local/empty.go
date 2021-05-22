package local

import (
	"context"
	"errors"

	"github.com/anson-xcloud/xdp-demo/pkg/joinpoint"
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

func (p *Provider) Connect(ctx context.Context, addr string) (joinpoint.Transport, []string, error) {
	return &Transport{source: p.source}, nil, nil
}

func (p *Provider) Serve(ctx context.Context, jr joinpoint.Request) {
	jr.ResponseStatus(joinpoint.StatusOK)
}

type Transport struct {
	source chan joinpoint.Request
}

func (t *Transport) Recv(ctx context.Context) (joinpoint.Request, error) {
	select {
	case req, ok := <-t.source:
		if !ok {
			return nil, errors.New("source close")
		}
		return req, nil
	case <-ctx.Done():
		return nil, errors.New("done")
	}
}
