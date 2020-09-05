package xdp

import (
	"sync"
)

type Client struct {
	mtx sync.Mutex

	AppID string

	conn *Connection
}

func NewClient() *Client {
	c := new(Client)
	return c
}
