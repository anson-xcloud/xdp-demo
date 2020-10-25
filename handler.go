package xdp

import (
	"encoding/json"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/anson-xcloud/xdp-demo/api"
	"github.com/golang/protobuf/proto"
)

// Handler tcp/udp raw data handler
type Handler interface {
	ServeXDP(ResponseWriter, *Request)
}

type StreamHandler interface {
	Handler

	ServeXDPConnect(*Session)

	ServeXDPClose(*Session)
}

type handlerServeFunc func(ResponseWriter, *Request)

func (f handlerServeFunc) Serve(res ResponseWriter, req *Request) {
	f(res, req)
}

// Request request data info
type Request struct {
	Session *Session
	Api     string
	Headers map[string]string
	Body    []byte

	reqTime time.Time
	packet  *Packet
}

// ResponseWriter response write
type ResponseWriter interface {
	WriteStatus(statusCode int)

	Write(data interface{}) error
}

type responseWriter struct {
	sv *Server

	resp api.SessionOnRecvNotifyResponse

	writed int32

	req *Request
}

func (r *responseWriter) Write(data interface{}) error {
	if r.resp.Status == 0 {
		r.resp.Status = uint32(http.StatusOK)
	}

	switch rd := data.(type) {
	case []byte:
		r.resp.Body = rd
	default:
		bd, err := json.Marshal(data)
		if err != nil {
			return err
		}
		r.resp.Body = bd
	}

	return r.write()
}

func (r *responseWriter) WriteStatus(statusCode int) {
	r.resp.Status = uint32(statusCode)

	r.write()
}

func (r *responseWriter) write() error {
	if !atomic.CompareAndSwapInt32(&r.writed, 0, 1) {
		return ErrTwiceWriteHTTPResponse
	}

	defer func() {
		r.sv.Logger().Info("http %s status %v cost %0.3fs", r.req.Api,
			r.resp.Status, time.Since(r.req.reqTime).Seconds())
	}()

	data, err := proto.Marshal(&r.resp)
	if err != nil {
		return err
	}

	var p Packet
	p.ID = r.req.packet.ID
	p.Cmd = r.req.packet.Cmd
	if p.ID != 0 {
		p.Flag |= flagRPCResponse
	}
	p.Data = data
	return r.sv.conn.write(&p)
}

// ServeMux support multi handler based on path
// path support syntax
type ServeMux struct {
	mtx sync.RWMutex

	onConnect, onClose func(*Session)
	hs                 map[string]handlerServeFunc
}

// NewServeMux create *ServeMux
func NewServeMux() *ServeMux {
	sm := new(ServeMux)
	sm.hs = make(map[string]handlerServeFunc)
	return sm
}

// HandleFunc register handler func
func (s *ServeMux) HandleFunc(pattern string, h handlerServeFunc) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.hs[pattern] = h
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

func (s *ServeMux) getHandler(req *Request) handlerServeFunc {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	h := s.hs[req.Api]
	return h
}

// ServeConnect implement Handler.ServeConnect
func (s *ServeMux) ServeXDPConnect(sess *Session) {
	s.mtx.RLock()
	fn := s.onConnect
	s.mtx.RUnlock()

	if fn != nil {
		fn(sess)
	}
}

// Serve implement Handler.Serve
func (s *ServeMux) ServeXDP(res ResponseWriter, req *Request) {
	if h := s.getHandler(req); h != nil {
		h.Serve(res, req)
		return
	}

	if req.packet.ID != 0 {
		res.WriteStatus(http.StatusNotFound)
	}
}

// ServeClose implement Handler.ServeClose
func (s *ServeMux) ServeXDPClose(sess *Session) {
	s.mtx.RLock()
	fn := s.onClose
	s.mtx.RUnlock()

	if fn != nil {
		fn(sess)
	}
}

var defaultServeMux = NewServeMux()

// HandleFunc call defaultServeMux.HandleFunc
func HandleFunc(pattern string, h handlerServeFunc) {
	defaultServeMux.HandleFunc(pattern, h)
}

// HandleConnect call defaultServeMux.HandleConnect
func HandleConnect(fn func(*Session)) {
	defaultServeMux.HandleConnect(fn)
}

// HandleClose call defaultServeMux.HandleClose
func HandleClose(fn func(*Session)) {
	defaultServeMux.HandleClose(fn)
}
