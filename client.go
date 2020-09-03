package xdp

import (
	"sync"
)

type Client struct {
	mtx sync.Mutex

	conn *Connection
}
