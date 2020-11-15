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

type HandlerRemoteType int

// HandlerRemote handle remote condition
// note: Anonymous donot have HandlerRemoteTypeServer
type HandlerRemote struct {
	Type  HandlerRemoteType
	Appid string
}

const (
	HandlerRemoteTypeBoth   HandlerRemoteType = 0
	HandlerRemoteTypeUser   HandlerRemoteType = 1
	HandlerRemoteTypeServer HandlerRemoteType = 2
)

const (
	HandlerRemoteAppidOwn       = "."
	HandlerRemoteAppidAnonymous = "?"
	HandlerRemoteAppidAll       = "*"
)

var (
	HandlerRemoteAnonymousUser = HandlerRemote{Type: HandlerRemoteTypeUser, Appid: HandlerRemoteAppidAnonymous}
	HandlerRemoteAll           = HandlerRemote{Type: HandlerRemoteTypeBoth, Appid: HandlerRemoteAppidAll}
	HandlerRemoteAllUser       = HandlerRemote{Type: HandlerRemoteTypeUser, Appid: HandlerRemoteAppidAll}
	HandlerRemoteAllServer     = HandlerRemote{Type: HandlerRemoteTypeServer, Appid: HandlerRemoteAppidAll}
	HandlerRemoteOwnUser       = HandlerRemote{Type: HandlerRemoteTypeUser, Appid: HandlerRemoteAppidOwn}
	HandlerRemoteOwnServer     = HandlerRemote{Type: HandlerRemoteTypeServer, Appid: HandlerRemoteAppidOwn}
)

// remoteHandler handler depend on remote
// own, anonymous, other will be selected first, if not found, then get all
type remoteHandler struct {
	own, anonymous Handler
	other          map[string]Handler

	all Handler
}

func newRemoteHandler() *remoteHandler {
	return &remoteHandler{other: make(map[string]Handler)}
}

func (s *remoteHandler) addHandler(appid string, h Handler) {
	switch appid {
	case HandlerRemoteAppidAnonymous:
		s.anonymous = h
	case HandlerRemoteAppidOwn:
		s.own = h
	case HandlerRemoteAppidAll:
		s.all = h
	default:
		s.other[appid] = h
	}
}

func (s *remoteHandler) getHandler(svr Server, req *Request) Handler {
	var h Handler
	switch req.Appid {
	case "":
		h = s.anonymous
	case svr.GetAddr().AppID:
		h = s.own
	default:
		h = s.other[req.Appid]
	}

	if h != nil {
		return h
	}
	return s.all
}

// ServeMux is an XDP request multiplexer.
type ServeMux struct {
	mtx sync.RWMutex

	handlers []map[string]*remoteHandler
}

// NewServeMux create *ServeMux
func NewServeMux() *ServeMux {
	sm := new(ServeMux)
	sm.handlers = make([]map[string]*remoteHandler, HandlerRemoteTypeServer+1)
	for i := HandlerRemoteTypeBoth; i <= HandlerRemoteTypeServer; i++ {
		sm.handlers[i] = make(map[string]*remoteHandler)
	}
	return sm
}

// HandleFunc register handler func
func (s *ServeMux) HandleFunc(remote HandlerRemote, pattern string, h HandlerFunc) {
	s.Handle(remote, pattern, h)
}

// HandleFunc register handler func
func (s *ServeMux) Handle(remote HandlerRemote, pattern string, h Handler) {
	if remote.Type < HandlerRemoteTypeBoth || remote.Type > HandlerRemoteTypeServer {
		panic("invalid remote type")
	}

	s.mtx.Lock()
	defer s.mtx.Unlock()

	hs := s.handlers[remote.Type]
	sh, ok := hs[pattern]
	if !ok {
		sh = newRemoteHandler()
		hs[pattern] = sh
	}
	sh.addHandler(remote.Appid, h)
}

func (s *ServeMux) getHandler(svr Server, req *Request) Handler {
	fn := func(stype HandlerRemoteType) Handler {
		if sh, ok := s.handlers[stype][req.Api]; ok {
			return sh.getHandler(svr, req)
		}
		return nil
	}

	s.mtx.RLock()
	defer s.mtx.RUnlock()

	stype := HandlerRemoteTypeServer
	if req.Sid != "" {
		stype = HandlerRemoteTypeUser
	}
	if h := fn(stype); h != nil {
		return h
	}
	return fn(HandlerRemoteTypeBoth)
}

// Serve implement Handler.Serve
func (s *ServeMux) Serve(svr Server, req *Request) {
	var ec uint32
	defer func() {
		ts := time.Since(req.reqTime).Seconds()
		svr.GetLogger().Debug("[XDP] %s serve %s cost %.3fs, ec(%d)", svr.GetAddr().AppID, req.Api, ts, ec)
	}()

	h := s.getHandler(svr, req)
	if h == nil {
		ec = 100
		svr.ReplyError(req, ec, "")
		return
	}
	h.Serve(svr, req)
}

var defaultServeMux = NewServeMux()

// HandleFunc call defaultServeMux.HandleFunc
func HandleFunc(remote HandlerRemote, pattern string, h HandlerFunc) {
	defaultServeMux.HandleFunc(remote, pattern, h)
}

func Handle(remote HandlerRemote, pattern string, h Handler) {
	defaultServeMux.Handle(remote, pattern, h)
}
