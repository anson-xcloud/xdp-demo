package local

import (
	"strconv"
	"sync"

	"github.com/anson-xcloud/xdp-demo/pkg/joinpoint"
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

func (r *Request) Discription() string {
	return strconv.FormatInt(int64(r.ID), 10)
}

func (r *Request) GetResponseWriter() joinpoint.ResponseWriter {
	return &ResponseWriter{ev: r}
}

type ResponseWriter struct {
	ev *Request
}

func (r *ResponseWriter) Write(interface{}) {
	r.ev.ch <- nil
}

func (r *ResponseWriter) WriteStatus(st *joinpoint.Status) {
	r.ev.ch <- st
}
