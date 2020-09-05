package xdp

import (
	"io"
	"io/ioutil"
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
	n, err := writeString(w, r.Key)
	return int64(n), err
}

type RegisterRequest struct {
	Config string
}

func (r *RegisterRequest) Cmd() int {
	return svrCmdRegister
}

func (r *RegisterRequest) WriteTo(w io.Writer) (int64, error) {
	n, err := writeString(w, r.Config)
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
	var n, total int
	var err error
	if n, err = writeString(w, r.SessionID); err != nil {
		return 0, err
	}
	total += n

	if n, err = writeString(w, r.OpenID); err != nil {
		return 0, err
	}
	total += n

	data, err := ioutil.ReadAll(r.Data)
	if err != nil {
		return 0, err
	}
	if n, err = w.Write(data); err != nil {
		return 0, err
	}
	total += n

	return int64(total), err
}
