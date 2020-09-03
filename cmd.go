package xdp

import (
	"bytes"
	"io"
)

const (
	svrCmdHandshake = 1
	svrCmdData      = 2
)

type HandshakeRequest struct {
	AppID string
	Key   string
}

func (r *HandshakeRequest) GetBody() (io.Reader, error) {
	bb := bytes.NewBuffer(nil)
	bb.WriteString(r.Key)
	// if err := writeString(bb, r.Key); err != nil {
	// 	return nil, err
	// }
	return bb, nil
}

type DataTransfer struct {
	SessionID string
	OpenID    string
	Data      io.Reader
}

func (r *DataTransfer) GetBody() (io.Reader, error) {
	bb := bytes.NewBuffer(nil)
	bb.WriteString(r.SessionID)
	// if err := writeString(bb, r.Key); err != nil {
	// 	return nil, err
	// }
	return bb, nil
}
