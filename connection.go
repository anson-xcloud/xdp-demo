package xdp

import (
	"context"
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

const flagRPCResponse = 0x1

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
	Logger

	mtx sync.Mutex

	address string
	nc      net.Conn

	rpcid   uint32
	callers map[uint32]*caller
}

func newConnection(address string) *Connection {
	conn := new(Connection)
	conn.address = address
	conn.callers = make(map[uint32]*caller)
	return conn
}

// Connect do connect
func (c *Connection) Connect() error {
	nc, err := net.Dial("tcp", c.address)
	if err != nil {
		return err
	}

	c.nc = nc
	return nil
}

func (c *Connection) recv(handler func(p *Packet)) {
	nc := c.nc
	for {
		var p Packet
		if err := p.Read(nc); err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				continue
			}
			c.Error("read %s packet error:%s", c.nc.RemoteAddr(), err)
			break
		}

		if p.Flag&flagRPCResponse != 0 {
			c.mtx.Lock()
			if ctx, ok := c.callers[p.ID]; !ok {
				c.mtx.Unlock()
				c.Error("%s unexisted rpc return %d", c.nc.RemoteAddr(), p.ID)
			} else {
				delete(c.callers, p.ID)
				c.mtx.Unlock()
				ctx.ch <- &p
			}
		} else {
			go handler(&p)
		}
	}
}

// Close do close connection
func (c *Connection) Close() {
	nc := c.nc
	if nc != nil {
		c.nc = nil
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
	if err := c.write(p); err != nil {
		return nil, err
	}

	tctx, _ := context.WithTimeout(ctx, time.Second*30)
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

func (c *Connection) write(p *Packet) error {
	p.Length = uint32(len(p.Data))
	if err := p.Write(c.nc); err != nil {
		return err
	}
	return nil
}
