package network

import (
	"context"
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/anson-xcloud/xdp-demo/pkg/xlog"
)

type caller struct {
	id uint32
	ch chan *Packet
}

func newCaller() *caller {
	c := new(caller)
	c.ch = make(chan *Packet)
	return c
}

// Connection tcp connection
type Connection struct {
	xlog.Logger

	mtx sync.Mutex

	address string
	nc      net.Conn

	cec error //close error

	recved chan *Packet

	rpcid   uint32
	callers map[uint32]*caller
}

func NewConnection() *Connection {
	conn := new(Connection)
	conn.recved = make(chan *Packet, 1024)
	conn.callers = make(map[uint32]*caller)
	return conn
}

// Connect do connect
func (c *Connection) Connect(addr string) error {
	c.address = addr
	nc, err := net.Dial("tcp", c.address)
	if err != nil {
		return err
	}

	c.nc = nc
	go c.Recv(func(p *Packet) {
		c.recved <- p
	})
	return nil
}

func (c *Connection) Read() *Packet {
	return <-c.recved
}

func (c *Connection) Recv2() <-chan *Packet {
	return c.recved
}

func (c *Connection) Recv(handler func(p *Packet)) error {
	nc := c.nc
	for {
		var p Packet
		if err := p.Read(nc); err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				continue
			}

			c.Errorf("read %s packet error:%s", nc.RemoteAddr(), err)
			return err
		}

		if p.Flag&FlagRPCResponse == 0 {
			handler(&p)
			continue
		}

		c.mtx.Lock()
		if ctx, ok := c.callers[p.ID]; !ok {
			c.mtx.Unlock()
			c.Errorf("%s unexisted rpc return %d", nc.RemoteAddr(), p.ID)
		} else {
			delete(c.callers, p.ID)
			c.mtx.Unlock()
			ctx.ch <- &p
		}
	}
}

// Close do close connection
func (c *Connection) Close(err error) {
	c.cec = err

	nc := c.nc
	if nc != nil {
		nc.Close()
	}
}

// Call do rpc call
func (c *Connection) Call(ctx context.Context, p *Packet) (*Packet, error) {
	caller := newCaller()
	caller.id = atomic.AddUint32(&c.rpcid, 1)
	c.mtx.Lock()
	c.callers[caller.id] = caller
	c.mtx.Unlock()

	p.ID = caller.id
	if err := c.Write(p); err != nil {
		return nil, err
	}

	tctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()
	select {
	case p, ok := <-caller.ch:
		if !ok {
			return nil, errors.New("rpc wait fail")
		}
		if p.Ec != 0 {
			return nil, errors.New("response ec")
		}
		return p, nil
	case <-tctx.Done():
		c.mtx.Lock()
		delete(c.callers, caller.id)
		c.mtx.Unlock()
		return nil, ErrTimeout
	}
}

func (c *Connection) Write(p *Packet) error {
	p.Length = uint32(len(p.Data))
	if err := p.Write(c.nc); err != nil {
		return err
	}
	return nil
}
