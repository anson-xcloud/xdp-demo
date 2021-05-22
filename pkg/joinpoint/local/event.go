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

func (r *Request) String() string {
	return strconv.FormatInt(int64(r.ID), 10)
}

func (r *Request) Response(interface{}) {
	r.ch <- nil
}

func (r *Request) ResponseStatus(st *joinpoint.Status) {
	r.ch <- st
}
