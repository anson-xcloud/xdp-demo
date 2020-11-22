package server

import (
	"container/list"
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
	handlerRemoteTypeUserBitsize = iota
	handlerRemoteTypeServerBitsize
	handlerRemoteTypeXcloudBitsize
)

const (
	HandlerRemoteTypeUser         HandlerRemoteType = 1 << handlerRemoteTypeUserBitsize
	HandlerRemoteTypeServer       HandlerRemoteType = 1 << handlerRemoteTypeServerBitsize
	HandlerRemoteTypeXcloud       HandlerRemoteType = 1 << handlerRemoteTypeXcloudBitsize
	HandlerRemoteTypeUserOrServer HandlerRemoteType = HandlerRemoteTypeUser | HandlerRemoteTypeServer
	HandlerRemoteTypeAll          HandlerRemoteType = HandlerRemoteTypeUserOrServer | HandlerRemoteTypeXcloud
)

const (
	HandlerRemoteAppidOwn       = "."
	HandlerRemoteAppidAnonymous = "?"
	HandlerRemoteAppidAll       = "*"
)

var (
	HandlerRemoteAnonymousUser = HandlerRemote{Type: HandlerRemoteTypeUser, Appid: HandlerRemoteAppidAnonymous}
	HandlerRemoteAll           = HandlerRemote{Type: HandlerRemoteTypeUserOrServer, Appid: HandlerRemoteAppidAll}
	HandlerRemoteAllUser       = HandlerRemote{Type: HandlerRemoteTypeUser, Appid: HandlerRemoteAppidAll}
	HandlerRemoteAllServer     = HandlerRemote{Type: HandlerRemoteTypeServer, Appid: HandlerRemoteAppidAll}
	HandlerRemoteOwnUser       = HandlerRemote{Type: HandlerRemoteTypeUser, Appid: HandlerRemoteAppidOwn}
	HandlerRemoteOwnServer     = HandlerRemote{Type: HandlerRemoteTypeServer, Appid: HandlerRemoteAppidOwn}
	HandlerRemoteXcloud        = HandlerRemote{Type: HandlerRemoteTypeXcloud}
)

// typedHandler handler depend on remote type
// own, anonymous, other will be selected first, if not found, then get all
type typedHandler struct {
	typ HandlerRemoteType

	xcloud Handler

	own, anonymous, all Handler
	apps                map[string]Handler
}

func newRemoteHandler(remote HandlerRemote, h Handler) *typedHandler {
	t := &typedHandler{typ: remote.Type, apps: make(map[string]Handler)}

	switch remote.Appid {
	case HandlerRemoteAppidAnonymous:
		t.anonymous = h
	case HandlerRemoteAppidOwn:
		t.own = h
	case HandlerRemoteAppidAll:
		t.all = h
	default:
		t.apps[remote.Appid] = h
	}

	if t.typ&HandlerRemoteTypeXcloud != 0 {
		t.xcloud = h
	}

	return t
}

func (t *typedHandler) getHandler(svr Server, typ HandlerRemoteType, appid string) Handler {
	if typ == HandlerRemoteTypeXcloud {
		return t.xcloud
	}

	var h Handler
	switch appid {
	case "":
		h = t.anonymous
	case svr.GetAddr().AppID:
		h = t.own
	default:
		h = t.apps[appid]
	}
	if h != nil {
		return h
	}
	return t.all
}

// ServeMux is an XDP request multiplexer.
type ServeMux struct {
	mtx sync.RWMutex

	handlers map[string]*list.List //*typedHandler
}

// NewServeMux create *ServeMux
func NewServeMux() *ServeMux {
	sm := new(ServeMux)
	sm.handlers = make(map[string]*list.List)
	return sm
}

// HandleFunc register handler func
func (s *ServeMux) HandleFunc(remote HandlerRemote, api string, h HandlerFunc) {
	s.Handle(remote, api, h)
}

// HandleFunc register handler func
func (s *ServeMux) Handle(remote HandlerRemote, api string, h Handler) {
	if remote.Type < HandlerRemoteTypeUser || remote.Type > HandlerRemoteTypeAll {
		panic("invalid remote type")
	}

	s.mtx.Lock()
	defer s.mtx.Unlock()
	hs, ok := s.handlers[api]
	if !ok {
		hs = list.New()
		s.handlers[api] = hs
	}
	hs.PushBack(newRemoteHandler(remote, h))
}

func (s *ServeMux) getHandler(svr Server, req *Request) Handler {
	var typ HandlerRemoteType
	if req.Sid != "" {
		typ = HandlerRemoteTypeUser
	} else if req.Appid != "" {
		typ = HandlerRemoteTypeServer
	} else {
		typ = HandlerRemoteTypeXcloud
	}

	s.mtx.RLock()
	defer s.mtx.RUnlock()
	hs, ok := s.handlers[req.Api]
	if !ok {
		return nil
	}

	for it := hs.Front(); it != nil; it = it.Next() {
		th := it.Value.(*typedHandler)
		if h := th.getHandler(svr, typ, req.Appid); h != nil {
			return h
		}
	}
	return nil
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
func HandleFunc(remote HandlerRemote, api string, h HandlerFunc) {
	defaultServeMux.HandleFunc(remote, api, h)
}

func Handle(remote HandlerRemote, api string, h Handler) {
	defaultServeMux.Handle(remote, api, h)
}
