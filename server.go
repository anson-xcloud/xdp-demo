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
	"sync"
)

type Server struct {
	mtx sync.Mutex

	conn *Connection

	AppID, AppSecret string

	handler Handler
}

func NewServer() *Server {
	svr := new(Server)
	return svr
}

func (s *Server) Serve() error {
	if err := s.Init(); err != nil {
		return err
	}

	return nil
}

func (s *Server) Init() error {
	ap, err := s.getAccessPoint()
	if err != nil {
		return err
	}

	err = s.handshake(ap)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Write(p []byte) error {
	// func (s *Server) Write(sess Session, p []byte) error {
	// var dt DataTransfer
	// dt.SessionID = sess.ID()
	// dt.Data = p
	// data, err := dt.Marshal()
	// if err != nil {
	// 	return err
	// }

	// var pk Packet
	// pk.Cmd = CmdData
	// pk.Length = uint32(len(p))
	// pk.Data = data
	// return s.cli.write(&pk)
	return nil
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
	Addr string `json:"addr"`
	Key  string `json:"key"`
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

func (s *Server) handshake(ap *AccessPoint) error {
	conn := newConnection("tcp", ap.Addr)
	if err := conn.Connect(); err != nil {
		return err
	}

	var reqbody HandshakeRequest
	reqbody.AppID = s.AppID
	reqbody.Key = ap.Key

	var req Request
	req.Cmd = svrCmdHandshake
	_, err := conn.Call(context.Background(), &req)
	if err != nil {
		conn.Close()
		return err
	}
	s.conn = conn
	return nil
}
