package xdp

import (
	"bytes"
	"errors"
	"net/http"
	"sync/atomic"

	"github.com/anson-xcloud/xdp-demo/api"
)

type Request struct {
	Session *Session
	Path    string
	Data    []byte
}

type HTTPRequest struct {
	Request

	Method  string
	Headers map[string]string
}

type HTTPResponseWriter interface {
	WriteHeader(statusCode int)

	Write(data []byte)
}

type httpResponseWriter struct {
	p *Packet

	sv *xdpServer

	api.HTTPResponse

	writed int32
}

func (r *httpResponseWriter) Write(data []byte) {
	if r.Status == 0 {
		r.Status = uint32(http.StatusOK)
	}

	r.Body = string(data)

	r.write()
}

func (r *httpResponseWriter) WriteHeader(statusCode int) {
	r.Status = uint32(statusCode)

	r.write()
}

func (r *httpResponseWriter) write() error {
	if !atomic.CompareAndSwapInt32(&r.writed, 0, 1) {
		return errors.New("rewrite response")
	}

	bb := bytes.NewBuffer(nil)
	if _, err := r.HTTPResponse.WriteTo(bb); err != nil {
		return err
	}
	r.p.Flag |= flagRPCResponse
	r.p.Data = bb.Bytes()
	return r.sv.conn.write(r.p)
}
