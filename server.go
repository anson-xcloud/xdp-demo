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

func ListenAndServe(url string, handler Handler) {

}

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

type httpResponse struct {
	p *Packet

	sv *Server

	resp HTTPResponse
}

func newHTTPResponse() *httpResponse {
	r := new(httpResponse)
	r.resp.Status = http.StatusOK
	return r
}

func (r *httpResponse) Header() http.Header {
	return http.Header{}
}

func (r *httpResponse) Write(data []byte) (int, error) {
	r.resp.Body = string(data)
	bb := bytes.NewBuffer(nil)
	if _, err := r.resp.WriteTo(bb); err != nil {
		return 0, err
	}

	r.p.Flag |= flagRPCResponse
	r.p.Data = bb.Bytes()
	err := r.sv.conn.write(r.p)
	return 0, err
}

func (r *httpResponse) WriteHeader(statusCode int) {
	r.resp.Status = uint32(statusCode)
}

type Server struct {
	mtx sync.Mutex

	conn *Connection

	AppID, AppSecret, Config string

	Handler Handler
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
	s.conn = conn
	go conn.recv(s.process)

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

	var p Packet
	p.Cmd = svrCmdData
	p.Data = bb.Bytes()
	return s.conn.write(&p)
}

func (s *Server) process(p *Packet) {
	switch p.Cmd {
	case svrCmdData:
		bb := bytes.NewBuffer(p.Data)
		var dt DataTransfer
		if _, err := dt.ReadFrom(bb); err != nil {
			return
		}
		sess := &Session{}
		sess.SessionID = dt.SessionID
		if s.Handler != nil {
			s.Handler.Serve(sess, 0, dt.Data)
		}
	case svrCmdHTTP:
		bb := bytes.NewBuffer(p.Data)
		var dt HTTPRequest
		if _, err := dt.ReadFrom(bb); err != nil {
			return
		}

		req, _ := http.NewRequest("Get", dt.Path, nil)
		res := newHTTPResponse()
		res.p = p
		res.sv = s
		if s.Handler != nil {
			s.Handler.ServeHTTP(res, req)
		}
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

func (s *Server) call(api ServerAPI) error {
	bb := bytes.NewBuffer(nil)
	if _, err := api.WriteTo(bb); err != nil {
		return err
	}

	var p Packet
	p.Cmd = uint32(api.Cmd())
	p.Data = bb.Bytes()
	_, err := s.conn.Call(context.Background(), &p)
	return err
}

func (s *Server) push(api ServerAPI) error {
	bb := bytes.NewBuffer(nil)
	if _, err := api.WriteTo(bb); err != nil {
		return err
	}

	var p Packet
	p.Cmd = uint32(api.Cmd())
	p.Data = bb.Bytes()
	return s.conn.write(&p)
}
