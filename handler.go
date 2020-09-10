package xdp

import (
	"net/http"
	"sync"
)

type Request struct {
	Session *Session
	Path    string
	Data    []byte
}

type HTTPRequest struct {
	Request

	Method  string
	Headers map[string]string
}

type HTTPResponseWriter interface {
	WriteHeader(statusCode int)

	Write(data []byte)
}

type Handler interface {
	Serve(*Request)
}

type HTTPHandler interface {
	ServeHTTP(HTTPResponseWriter, *HTTPRequest)
}

type HandlerFunc func(*Request)

func (f HandlerFunc) Serve(req *Request) {
	f(req)
}

type HTTPHandlerFunc func(HTTPResponseWriter, *HTTPRequest)

func (f HTTPHandlerFunc) ServeHTTP(res HTTPResponseWriter, req *HTTPRequest) {
	f(res, req)
}

type ServeMux struct {
	mtx sync.RWMutex

	hs  map[string]Handler
	hhs map[string]HTTPHandler
}

func NewServeMux() *ServeMux {
	sm := new(ServeMux)
	sm.hs = make(map[string]Handler)
	sm.hhs = make(map[string]HTTPHandler)
	return sm
}

func (m *ServeMux) HandleFunc(pattern string, h HandlerFunc) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.hs[pattern] = h
}

func (m *ServeMux) HTTPHandleFunc(pattern string, hh HTTPHandlerFunc) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.hhs[pattern] = hh
}

func (m *ServeMux) ServeHTTP(res HTTPResponseWriter, req *HTTPRequest) {
	m.mtx.RLock()
	hh, ok := m.hhs[req.Path]
	if !ok {
		m.mtx.RUnlock()

		res.WriteHeader(http.StatusNotFound)
		return
	}
	m.mtx.RUnlock()

	hh.ServeHTTP(res, req)
}

func (m *ServeMux) Serve(req *Request) {
	m.mtx.RLock()
	h, ok := m.hs[req.Path]
	if !ok {
		m.mtx.RUnlock()
		return
	}
	m.mtx.RUnlock()

	h.Serve(req)
}
