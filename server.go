package xdp

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
)

type ServerInterface interface {
	Ping() error

	Serve() error

	Write(*Session, io.Reader)
}

type Session struct {
	Addr      string
	OpenID    string
	SessionID string
}

type Server struct {
	mtx sync.Mutex

	conn *Connection

	AppID, AppSecret, Config string

	handler Handler
}

func NewServer() *Server {
	svr := new(Server)
	return svr
}

func (s *Server) Serve() error {
	ap, err := s.getAccessPoint()
	if err != nil {
		return err
	}

	conn := newConnection("tcp", ap.Addr)
	if err := conn.Connect(); err != nil {
		return err
	}
	go conn.recv()
	s.conn = conn

	if err = s.call(&HandshakeRequest{
		AppID:     s.AppID,
		AccessKey: ap.AccessKey,
	}); err != nil {
		return err
	}

	return nil
}

func (s *Server) Send(sess *Session, data []byte) error {
	var dt DataTransfer
	dt.SessionID = sess.SessionID
	dt.Data = data

	bb := bytes.NewBuffer(nil)
	if _, err := dt.WriteTo(bb); err != nil {
		return err
	}

	var req Request
	req.Cmd = svrCmdData
	req.Body = bb
	return s.conn.Push(&req)
}

func (s *Server) process() {
	select {
	case p := <-s.conn.Fetch():
		_ = p
		// var dt DataTransfer
		// if err := dt.Unmarshal(p.Data); err != nil {
		// 	return nil, nil //err
		// }

		// s := &session{}
		// s.id = dt.SessionID
		// s.srv.OnSessionData(s, p.Data)
	}
}

func (s *Server) signUrl(vals url.Values) {
	md5str := fmt.Sprintf("%s%s", vals.Encode(), s.AppSecret)
	m := md5.New()
	token := hex.EncodeToString(m.Sum([]byte(md5str)))
	vals.Set("token", token)
}

type AccessPoint struct {
	Addr      string `json:"addr"`
	AccessKey string `json:"access_key"`
}

func (s *Server) getAccessPoint() (*AccessPoint, error) {
	values := make(url.Values)
	values.Set("appid", s.AppID)
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
	var ret AccessPoint
	if err := json.Unmarshal(data, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}

func (s *Server) call(api ServerAPI) error {
	bb := bytes.NewBuffer(nil)
	if _, err := api.WriteTo(bb); err != nil {
		return err
	}

	var req Request
	req.Cmd = api.Cmd()
	req.Body = bb
	_, err := s.conn.Call(context.Background(), &req)
	return err
}

func (s *Server) push(api ServerAPI) error {
	bb := bytes.NewBuffer(nil)
	if _, err := api.WriteTo(bb); err != nil {
		return err
	}

	var req Request
	req.Cmd = api.Cmd()
	req.Body = bb
	return s.conn.Push(&req)
}
