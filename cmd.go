package xdp

import (
	"encoding/binary"
	"io"
)

const (
	svrCmdHandshake = 1
	svrCmdData      = 2
	svrCmdHTTP      = 3
)

type ServerAPI interface {
	io.WriterTo

	Cmd() int
}

type HandshakeRequest struct {
	AppID     string
	AccessKey string
	Config    string
}

func (r *HandshakeRequest) Cmd() int {
	return svrCmdHandshake
}

func (r *HandshakeRequest) WriteTo(w io.Writer) (int64, error) {
	var n, total int
	var err error

	if n, err = writeString(w, r.AppID); err != nil {
		return 0, err
	}
	total += n

	if n, err = writeString(w, r.AccessKey); err != nil {
		return 0, err
	}
	total += n

	if n, err = writeString(w, r.Config); err != nil {
		return 0, err
	}
	total += n
	return int64(total), nil
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

	return int64(total), nil
}

type HTTPRequest struct {
	Path string
}

func (r *HTTPRequest) ReadFrom(rd io.Reader) (n int64, err error) {
	if r.Path, err = readString(rd); err != nil {
		return 0, err
	}
	return 0, nil
}

type HTTPResponse struct {
	Status uint32
	Body   string
}

func (r *HTTPResponse) Cmd() int {
	return svrCmdHTTP
}

func (r *HTTPResponse) WriteTo(w io.Writer) (int64, error) {
	if err := binary.Write(w, endian, &r.Status); err != nil {
		return 0, err
	}
	if _, err := writeString(w, r.Body); err != nil {
		return 0, err
	}
	return 0, nil
}
