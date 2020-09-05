package xdp

import (
	"io"
)

const (
	svrCmdHandshake = 1
	svrCmdRegister  = 2
	svrCmdData      = 3
)

type ServerAPI interface {
	io.WriterTo

	Cmd() int
}

type HandshakeRequest struct {
	AppID string
	Key   string
}

func (r *HandshakeRequest) Cmd() int {
	return svrCmdHandshake
}

func (r *HandshakeRequest) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write([]byte(r.Key))
	return int64(n), err
}

type RegisterRequest struct {
	Config string
}

func (r *RegisterRequest) Cmd() int {
	return svrCmdRegister
}

func (r *RegisterRequest) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write([]byte(r.Config))
	return int64(n), err
}

type DataTransfer struct {
	SessionID string
	OpenID    string
	Data      io.Reader
}

func (r *DataTransfer) Cmd() int {
	return svrCmdData
}

func (r *DataTransfer) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write([]byte(r.SessionID))
	return int64(n), err
}
