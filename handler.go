package xdp

import "net/http"

type Handler interface {
	Serve(sess *Session, cmd uint32, data []byte)

	ServeHTTP(http.ResponseWriter, *http.Request)
}
