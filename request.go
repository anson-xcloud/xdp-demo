package xdp

import (
	"net/http"
	"sync/atomic"

	"github.com/anson-xcloud/xdp-demo/api"
	"github.com/golang/protobuf/proto"
)

// Request request data info
type Request struct {
	Session *Session
	Path    string
	Data    []byte
}

// HTTPRequest http request info
type HTTPRequest struct {
	Request

	Method string
	Forms  map[string]string
}

// HTTPResponseWriter http response write
// http handler must call HTTPResponseWriter.WriteHander/HTTPResponseWriter.Write at once
type HTTPResponseWriter interface {
	WriteHeader(statusCode int)

	Write(data []byte)
}

type httpResponseWriter struct {
	p *Packet

	sv *Server

	resp api.SessionHTTPNotifyResponse

	writed int32
}

func (r *httpResponseWriter) Write(data []byte) {
	if r.resp.GetCode() == 0 {
		r.resp.Code = uint32(http.StatusOK)
	}
	r.resp.Body = data

	r.write()
}

func (r *httpResponseWriter) WriteHeader(statusCode int) {
	r.resp.Code = uint32(statusCode)

	r.write()
}

func (r *httpResponseWriter) write() error {
	if !atomic.CompareAndSwapInt32(&r.writed, 0, 1) {
		return ErrTwiceWriteHTTPResponse
	}

	data, err := proto.Marshal(&r.resp)
	if err != nil {
		return err
	}

	r.p.Flag |= flagRPCResponse
	r.p.Data = data
	return r.sv.conn.write(r.p)
}
