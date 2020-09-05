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
	AppID     string
	AccessKey string
}

func (r *HandshakeRequest) Cmd() int {
	return svrCmdHandshake
}

func (r *HandshakeRequest) WriteTo(w io.Writer) (int64, error) {
	n, err := writeString(w, r.AccessKey)
	return int64(n), err
}

type DataTransfer struct {
	SessionID string
	OpenID    string
	Data      []byte
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

	if n, err = w.Write(r.Data); err != nil {
		return 0, err
	}
	total += n

	return int64(total), err
}
