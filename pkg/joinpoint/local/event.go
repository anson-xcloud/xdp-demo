package local

import (
	"sync"
	"xcloud/pkg/clientapi/joinpoint"
)

var (
	evId uint64
	mtx  sync.Mutex
)

type Request struct {
	ID uint64

	ch chan error
}

func NewEvent() *Request {
	mtx.Lock()
	evId++
	id := evId
	mtx.Unlock()

	return &Request{ID: id, ch: make(chan error)}
}

func (e *Request) GetResponseWriter() joinpoint.ResponseWriter {
	return &ResponseWriter{ev: e}
}

type ResponseWriter struct {
	ev *Request
}

func (p *ResponseWriter) Write(interface{}) {
	p.ev.ch <- nil
}

func (p *ResponseWriter) WriteStatus(st *joinpoint.Status) {
	p.ev.ch <- st
}
