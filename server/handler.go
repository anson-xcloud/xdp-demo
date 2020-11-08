package server

import (
	"sync"
	"time"
)

type Handler interface {
	Serve(Server, *Request)

	// ServeConnect()
	// ServeClose()
}

type HandlerFunc func(Server, *Request)

func (h HandlerFunc) Serve(svr Server, req *Request) {
	h(svr, req)
}

// ServeMux is an XDP request multiplexer.
type ServeMux struct {
	mtx sync.RWMutex

	hs map[string]Handler
}

// NewServeMux create *ServeMux
func NewServeMux() *ServeMux {
	sm := new(ServeMux)
	sm.hs = make(map[string]Handler)
	return sm
}

// HandleFunc register handler func
func (s *ServeMux) HandleFunc(pattern string, h HandlerFunc) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.hs[pattern] = h
}

// HandleFunc register handler func
func (s *ServeMux) Handle(pattern string, h Handler) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.hs[pattern] = h
}

func (s *ServeMux) getHandler(api string) Handler {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	h := s.hs[api]
	return h
}

// Serve implement Handler.Serve
func (s *ServeMux) Serve(svr Server, req *Request) {
	defer func() {
		ms := time.Since(req.reqTime).Milliseconds()
		svr.GetLogger().Debug("[XDP] serve %s cost %dms", req.Api, ms)
	}()

	h := s.getHandler(req.Api)
	if h == nil {
		svr.ReplyError(req, 1, "")
		return
	}
	h.Serve(svr, req)
}

var defaultServeMux = NewServeMux()

// HandleFunc call defaultServeMux.HandleFunc
func HandleFunc(pattern string, h HandlerFunc) {
	defaultServeMux.HandleFunc(pattern, h)
}

func Handle(pattern string, h Handler) {
	defaultServeMux.Handle(pattern, h)
}
