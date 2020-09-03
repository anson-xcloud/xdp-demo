package xdp

type Handler interface {
	ServeXDP(ResponseWriter, *Request)
}
