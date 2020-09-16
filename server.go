package xdp

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/anson-xcloud/xdp-demo/api"
	"github.com/golang/protobuf/proto"
)

// Address for app address token
// format is   appid:appsecret
type Address struct {
	AppID, AppSecret string
}

// ParseAddress parse address string to *Address
func ParseAddress(addr string) (*Address, error) {
	sl := strings.Split(addr, ":")
	if len(sl) != 2 {
		return nil, ErrAddressFormat
	}

	return &Address{AppID: sl[0], AppSecret: sl[1]}, nil
}

// Server for app server
type Server struct {
	opts *Options

	addr *Address

	conn *Connection

	sessMgr *sessionManager
}

// NewServer create server
func NewServer(opt ...Option) *Server {
	var opts = defaultOptions
	for _, o := range opt {
		o(&opts)
	}

	svr := new(Server)
	svr.opts = &opts
	svr.sessMgr = newSessionManager()
	return svr
}

// Logger get server Logger
func (s *Server) Logger() Logger {
	return s.opts.Logger
}

// Serve start serve at addr, addr must be valid *Address
func (s *Server) Serve(addr string) error {
	ad, err := ParseAddress(addr)
	if err != nil {
		return err
	}
	s.addr = ad

	ap, err := s.getAccessPoint()
	if err != nil {
		return err
	}

	conn := newConnection()
	conn.Logger = s.opts.Logger
	if err := conn.Connect(ap.Addr); err != nil {
		return err
	}
	s.conn = conn

	go func() {
		if err = s.call(api.Cmd_Handshake, &api.HandshakeRequest{
			AppID:     s.addr.AppID,
			AccessKey: ap.AccessKey,
		}); err != nil {
			conn.Close(err)
			return
		}
	}()
	return conn.recv(s.process)
}

// Send send data to session client
func (s *Server) Send(sess *Session, data []byte) error {
	var sd api.SessionDataBiNotify
	sd.SessionID = sess.SessionID
	sd.Data = data
	return s.push(api.Cmd_SessionData, &sd)
}

// MultiSend multi send data to session at once
func (s *Server) MultiSend(sids []string, data []byte) error {
	var sd api.MultiSessionDataRequest
	sd.SessionIDs = sids
	sd.Data = data
	return s.push(api.Cmd_MultiSessionData, &sd)
}

func (s *Server) process(p *Packet) {
	switch p.Cmd {
	case uint32(api.Cmd_SessionConnect):
		s.handleData(p)
	case uint32(api.Cmd_SessionClose):
		s.handleData(p)
	case uint32(api.Cmd_SessionData):
		s.handleData(p)
	case uint32(api.Cmd_SessionHTTP):
		s.handleHTTP(p)
	}
}

func (s *Server) signURL(vals url.Values) {
	md5str := fmt.Sprintf("%s%s", vals.Encode(), s.addr.AppSecret)
	m := md5.New()
	token := hex.EncodeToString(m.Sum([]byte(md5str)))
	vals.Set("token", token)
}

// AccessPoint xcloud return access_point info
type AccessPoint struct {
	Addr      string `json:"addr"`
	AccessKey string `json:"access_key"`
}

func (s *Server) getAccessPoint() (*AccessPoint, error) {
	values := make(url.Values)
	values.Set("appid", s.addr.AppID)
	s.signURL(values)
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

func (s *Server) call(cmd api.Cmd, sa proto.Message) error {
	data, err := proto.Marshal(sa)
	if err != nil {
		return err
	}

	var p Packet
	p.Cmd = uint32(cmd)
	p.Data = data
	_, err = s.conn.Call(context.Background(), &p)
	return err
}

func (s *Server) push(cmd api.Cmd, sa proto.Message) error {
	data, err := proto.Marshal(sa)
	if err != nil {
		return err
	}

	var p Packet
	p.Cmd = uint32(cmd)
	p.Data = data
	return s.conn.write(&p)
}

func (s *Server) handleData(p *Packet) {
	var dt api.SessionDataBiNotify
	if err := proto.Unmarshal(p.Data, &dt); err != nil {
		s.opts.Logger.Debug("unmarshal handleData error:%s", err)
		return
	}

	var req Request
	req.Session = s.sessMgr.Get(dt.SessionID)
	req.Data = dt.Data
	s.opts.Handler.Serve(&req)
}

func (s *Server) handleSessionConnect(p *Packet) {
	var dt api.SessionConnectNotify
	if err := proto.Unmarshal(p.Data, &dt); err != nil {
		s.opts.Logger.Debug("unmarshal handleSessionConnect error:%s", err)
		return
	}

	var sess Session
	sess.SessionID = dt.SessionID
	// sess.Addr=dt.a
	sess.sv = s
	s.sessMgr.Add(&sess)
	s.opts.Handler.ServeConnect(&sess)
}

func (s *Server) handleSessionClose(p *Packet) {
	var dt api.SessionCloseNotify
	if err := proto.Unmarshal(p.Data, &dt); err != nil {
		s.opts.Logger.Debug("unmarshal SessionCloseNotify error:%s", err)
		return
	}

	sess := s.sessMgr.Get(dt.SessionID)
	s.sessMgr.Del(sess)
	s.opts.Handler.ServeClose(sess)
}

func (s *Server) handleHTTP(p *Packet) {
	var dt api.SessionHTTPNotify
	if err := proto.Unmarshal(p.Data, &dt); err != nil {
		s.opts.Logger.Debug("unmarshal SessionHTTPNotify error:%s", err)
		return
	}

	var req HTTPRequest
	req.Session = s.sessMgr.Get(dt.SessionID)
	req.Method = dt.Method
	req.Path = dt.Path
	req.Forms = dt.Forms
	req.Data = dt.Body
	req.reqTime = time.Now()

	res := &httpResponseWriter{}
	res.p = p
	res.sv = s
	res.writed = 0
	res.req = &req
	s.opts.HTTPHandler.ServeHTTP(res, &req)
}
