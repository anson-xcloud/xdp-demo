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

type HandlerSourceType int

// HandlerSource handle source condition
// note: Anonymous donot have HandlerSourceTypeServer
type HandlerSource struct {
	Type  HandlerSourceType
	Appid string
}

const (
	HandlerSourceTypeBoth   HandlerSourceType = 0
	HandlerSourceTypeUser   HandlerSourceType = 1
	HandlerSourceTypeServer HandlerSourceType = 2
)

const (
	HandlerSourceAppidOwn       = "."
	HandlerSourceAppidAnonymous = "?"
	HandlerSourceAppidAll       = "*"
)

var (
	HandlerSourceAnonymousUser = HandlerSource{Type: HandlerSourceTypeUser, Appid: HandlerSourceAppidAnonymous}
	HandlerSourceAll           = HandlerSource{Type: HandlerSourceTypeBoth, Appid: HandlerSourceAppidAll}
	HandlerSourceAllUser       = HandlerSource{Type: HandlerSourceTypeUser, Appid: HandlerSourceAppidAll}
	HandlerSourceAllServer     = HandlerSource{Type: HandlerSourceTypeServer, Appid: HandlerSourceAppidAll}
	HandlerSourceOwnUser       = HandlerSource{Type: HandlerSourceTypeUser, Appid: HandlerSourceAppidOwn}
	HandlerSourceOwnServer     = HandlerSource{Type: HandlerSourceTypeServer, Appid: HandlerSourceAppidOwn}
)

// sourceHandler handler depend on source
// own, anonymous, other will be selected first, if not found, then get all
type sourceHandler struct {
	own, anonymous Handler
	other          map[string]Handler

	all Handler
}

func newSourceHandler() *sourceHandler {
	return &sourceHandler{other: make(map[string]Handler)}
}

func (s *sourceHandler) addHandler(appid string, h Handler) {
	switch appid {
	case HandlerSourceAppidAnonymous:
		s.anonymous = h
	case HandlerSourceAppidOwn:
		s.own = h
	case HandlerSourceAppidAll:
		s.all = h
	default:
		s.other[appid] = h
	}
}

func (s *sourceHandler) getHandler(svr Server, req *Request) Handler {
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

	handlers []map[string]*sourceHandler
}

// NewServeMux create *ServeMux
func NewServeMux() *ServeMux {
	sm := new(ServeMux)
	sm.handlers = make([]map[string]*sourceHandler, HandlerSourceTypeServer+1)
	for i := HandlerSourceTypeBoth; i <= HandlerSourceTypeServer; i++ {
		sm.handlers[i] = make(map[string]*sourceHandler)
	}
	return sm
}

// HandleFunc register handler func
func (s *ServeMux) HandleFunc(source HandlerSource, pattern string, h HandlerFunc) {
	s.Handle(source, pattern, h)
}

// HandleFunc register handler func
func (s *ServeMux) Handle(source HandlerSource, pattern string, h Handler) {
	if source.Type < HandlerSourceTypeBoth || source.Type > HandlerSourceTypeServer {
		panic("invalid source type")
	}

	s.mtx.Lock()
	defer s.mtx.Unlock()

	hs := s.handlers[source.Type]
	sh, ok := hs[pattern]
	if !ok {
		sh = newSourceHandler()
		hs[pattern] = sh
	}
	sh.addHandler(source.Appid, h)
}

func (s *ServeMux) getHandler(svr Server, req *Request) Handler {
	fn := func(stype HandlerSourceType) Handler {
		if sh, ok := s.handlers[stype][req.Api]; ok {
			return sh.getHandler(svr, req)
		}
		return nil
	}

	s.mtx.RLock()
	defer s.mtx.RUnlock()

	stype := HandlerSourceTypeServer
	if req.Sid != "" {
		stype = HandlerSourceTypeUser
	}
	if h := fn(stype); h != nil {
		return h
	}
	return fn(HandlerSourceTypeBoth)
}

// Serve implement Handler.Serve
func (s *ServeMux) Serve(svr Server, req *Request) {
	defer func() {
		ms := time.Since(req.reqTime).Milliseconds()
		svr.GetLogger().Debug("[XDP] serve %s cost %dms", req.Api, ms)
	}()

	h := s.getHandler(svr, req)
	if h == nil {
		svr.ReplyError(req, 1, "")
		return
	}
	h.Serve(svr, req)
}

var defaultServeMux = NewServeMux()

// HandleFunc call defaultServeMux.HandleFunc
func HandleFunc(source HandlerSource, pattern string, h HandlerFunc) {
	defaultServeMux.HandleFunc(source, pattern, h)
}

func Handle(source HandlerSource, pattern string, h Handler) {
	defaultServeMux.Handle(source, pattern, h)
}
