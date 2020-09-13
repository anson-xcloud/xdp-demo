package xdp

import (
	"net/http"
	"sync"
)

// Handler tcp/udp raw data handler
type Handler interface {
	ServeConnect(*Session)

	Serve(*Request)

	ServeClose(*Session)
}

// HTTPHandler http data handler
type HTTPHandler interface {
	ServeHTTP(HTTPResponseWriter, *HTTPRequest)
}

type handlerServeFunc func(*Request)

func (f handlerServeFunc) Serve(req *Request) {
	f(req)
}

// HTTPHandlerFunc http handler for func
type HTTPHandlerFunc func(HTTPResponseWriter, *HTTPRequest)

func (f HTTPHandlerFunc) ServeHTTP(res HTTPResponseWriter, req *HTTPRequest) {
	f(res, req)
}

// ServeMux support multi handler based on path
// path support syntax
type ServeMux struct {
	mtx sync.RWMutex

	onConnect, onClose func(*Session)
	hs                 map[string]handlerServeFunc
	hhs                map[string]HTTPHandler
}

// NewServeMux create *ServeMux
func NewServeMux() *ServeMux {
	sm := new(ServeMux)
	sm.hs = make(map[string]handlerServeFunc)
	sm.hhs = make(map[string]HTTPHandler)
	return sm
}

// HandleFunc register handler func
func (s *ServeMux) HandleFunc(pattern string, h handlerServeFunc) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.hs[pattern] = h
}

// HTTPHandleFunc register http handler func
func (s *ServeMux) HTTPHandleFunc(pattern string, hh HTTPHandlerFunc) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.hhs[pattern] = hh
}

// HTTPHandle register http handler
func (s *ServeMux) HTTPHandle(pattern string, hh HTTPHandler) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.hhs[pattern] = hh
}

// HandleConnect register onConnect
func (s *ServeMux) HandleConnect(fn func(*Session)) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.onConnect = fn
}

// HandleClose register onClose
func (s *ServeMux) HandleClose(fn func(*Session)) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.onClose = fn
}

func (s *ServeMux) getHTTPHandler(req *HTTPRequest) HTTPHandler {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	hh := s.hhs[req.Path]
	return hh
}

func (s *ServeMux) getHandler(req *Request) handlerServeFunc {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	h := s.hs[req.Path]
	return h
}

// ServeConnect implement Handler.ServeConnect
func (s *ServeMux) ServeConnect(sess *Session) {
	s.mtx.RLock()
	fn := s.onConnect
	s.mtx.RUnlock()

	if fn != nil {
		fn(sess)
	}
}

// Serve implement Handler.Serve
func (s *ServeMux) Serve(req *Request) {
	if h := s.getHandler(req); h != nil {
		h.Serve(req)
	}
}

// ServeClose implement Handler.ServeClose
func (s *ServeMux) ServeClose(sess *Session) {
	s.mtx.RLock()
	fn := s.onClose
	s.mtx.RUnlock()

	if fn != nil {
		fn(sess)
	}
}

// ServeHTTP implement HTTPHandler.ServeHTTP
func (s *ServeMux) ServeHTTP(res HTTPResponseWriter, req *HTTPRequest) {
	if h := s.getHTTPHandler(req); h != nil {
		h.ServeHTTP(res, req)
	} else {
		res.WriteHeader(http.StatusNotFound)
	}
}

var defaultServeMux = NewServeMux()

// HandleFunc call defaultServeMux.HandleFunc
func HandleFunc(pattern string, h handlerServeFunc) {
	defaultServeMux.HandleFunc(pattern, h)
}

// HTTPHandleFunc call defaultServeMux.HTTPHandleFunc
func HTTPHandleFunc(pattern string, hh HTTPHandlerFunc) {
	defaultServeMux.HTTPHandleFunc(pattern, hh)
}

// HTTPHandle call defaultServeMux.HTTPHandle
func HTTPHandle(pattern string, hh HTTPHandler) {
	defaultServeMux.HTTPHandle(pattern, hh)
}

// HandleConnect call defaultServeMux.HandleConnect
func HandleConnect(fn func(*Session)) {
	defaultServeMux.HandleConnect(fn)
}

// HandleClose call defaultServeMux.HandleClose
func HandleClose(fn func(*Session)) {
	defaultServeMux.HandleClose(fn)
}
