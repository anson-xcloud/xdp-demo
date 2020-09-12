package xdp

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/anson-xcloud/xdp-demo/api"
)

type Address struct {
	AppID, AppSecret string
}

func ParseAddress(addr string) (*Address, error) {
	sl := strings.Split(addr, ":")
	if len(sl) != 2 {
		return nil, errors.New("url format error")
	}

	return &Address{AppID: sl[0], AppSecret: sl[1]}, nil
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

type Server interface {
	HandleFunc(pattern string, h HandlerFunc)
	HTTPHandleFunc(pattern string, hh HTTPHandlerFunc)
	Handle(pattern string, h Handler)
	HTTPHandle(pattern string, hh HTTPHandler)

	Serve(addr string) error

	Send(sess *Session, data []byte) error

	// Ping() error
}

type xdpServer struct {
	mtx sync.RWMutex

	Config string

	addr *Address
	conn *Connection

	hs  map[string]Handler
	hhs map[string]HTTPHandler
}

func NewServer() Server {
	svr := new(xdpServer)

	svr.hs = make(map[string]Handler)
	svr.hhs = make(map[string]HTTPHandler)
	return svr
}

func (s *xdpServer) HandleFunc(pattern string, h HandlerFunc) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.hs[pattern] = h
}

func (s *xdpServer) HTTPHandleFunc(pattern string, hh HTTPHandlerFunc) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.hhs[pattern] = hh
}

func (s *xdpServer) Handle(pattern string, h Handler) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.hs[pattern] = h
}

func (s *xdpServer) HTTPHandle(pattern string, hh HTTPHandler) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.hhs[pattern] = hh
}

func (s *xdpServer) getHTTPHandler(req *HTTPRequest) HTTPHandler {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	hh := s.hhs[req.Path]
	return hh
}

func (s *xdpServer) getHandler(req *Request) Handler {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	h := s.hs[req.Path]
	return h
}

func (s *xdpServer) Serve(addr string) error {
	ad, err := ParseAddress(addr)
	if err != nil {
		return err
	}
	s.addr = ad

	ap, err := s.getAccessPoint()
	if err != nil {
		return err
	}

	conn := newConnection("tcp", ap.Addr)
	if err := conn.Connect(); err != nil {
		return err
	}
	s.conn = conn
	go conn.recv(s.process)

	if err = s.call(&api.HandshakeRequest{
		AppID:     s.addr.AppID,
		AccessKey: ap.AccessKey,
	}); err != nil {
		return err
	}

	return nil
}

func (s *xdpServer) Send(sess *Session, data []byte) error {
	var dt api.DataTransfer
	dt.SessionID = sess.SessionID
	dt.Data = data

	bb := bytes.NewBuffer(nil)
	if _, err := dt.WriteTo(bb); err != nil {
		return err
	}

	var p Packet
	p.Cmd = api.CmdData
	p.Data = bb.Bytes()
	return s.conn.write(&p)
}

func (s *xdpServer) process(p *Packet) {
	switch p.Cmd {
	case api.CmdData:
		s.handleData(p)
	case api.CmdHTTP:
		s.handleHTTP(p)
	}
}

func (s *xdpServer) signUrl(vals url.Values) {
	md5str := fmt.Sprintf("%s%s", vals.Encode(), s.addr.AppSecret)
	m := md5.New()
	token := hex.EncodeToString(m.Sum([]byte(md5str)))
	vals.Set("token", token)
}

type AccessPoint struct {
	Addr      string `json:"addr"`
	AccessKey string `json:"access_key"`
}

func (s *xdpServer) getAccessPoint() (*AccessPoint, error) {
	values := make(url.Values)
	values.Set("appid", s.addr.AppID)
	s.signUrl(values)
	url := fmt.Sprintf("%s%s?%s", XCloudAddr, APIAccessPoint, values.Encode())

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	type AccessPointResult struct {
		Status int    `json:"status"`
		Msg    string `json:"msg"`
		AccessPoint
	}
	var ret AccessPointResult
	if err := json.Unmarshal(data, &ret); err != nil {
		return nil, err
	}
	if ret.Status != 0 {
		return nil, fmt.Errorf("response error(%d):%s", ret.Status, ret.Msg)
	}
	return &ret.AccessPoint, nil
}

func (s *xdpServer) call(sa api.ServerAPI) error {
	bb := bytes.NewBuffer(nil)
	if _, err := sa.WriteTo(bb); err != nil {
		return err
	}

	var p Packet
	p.Cmd = uint32(sa.Cmd())
	p.Data = bb.Bytes()
	_, err := s.conn.Call(context.Background(), &p)
	return err
}

func (s *xdpServer) push(sa api.ServerAPI) error {
	bb := bytes.NewBuffer(nil)
	if _, err := sa.WriteTo(bb); err != nil {
		return err
	}

	var p Packet
	p.Cmd = uint32(sa.Cmd())
	p.Data = bb.Bytes()
	return s.conn.write(&p)
}

func (s *xdpServer) handleData(p *Packet) {
	bb := bytes.NewBuffer(p.Data)
	var dt api.DataTransfer
	if _, err := dt.ReadFrom(bb); err != nil {
		return
	}

	sess := &Session{}
	sess.SessionID = dt.SessionID
	sess.sv = s

	var req Request
	req.Session = sess
	if h := s.getHandler(&req); h != nil {
		h.Serve(&req)
	}
}

func (s *xdpServer) handleHTTP(p *Packet) {
	bb := bytes.NewBuffer(p.Data)
	var dt api.HTTPRequest
	if _, err := dt.ReadFrom(bb); err != nil {
		return
	}

	res := &httpResponseWriter{}
	res.p = p
	res.sv = s
	res.writed = 0

	var req HTTPRequest
	req.Path = dt.Path
	if h := s.getHTTPHandler(&req); h != nil {
		h.ServeHTTP(res, &req)
	} else {
		res.WriteHeader(http.StatusNotFound)
	}
}
