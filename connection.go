package xdp

import (
	"context"
	"errors"
	"net"
	"sync"
	"sync/atomic"

	"github.com/anson-xcloud/xdp-demo/api"
)

const flagRPCResponse = 0x1

type caller struct {
	id uint32
	ch chan *api.Packet
}

func newCaller() *caller {
	c := new(caller)
	c.ch = make(chan *api.Packet)
	return c
}

type Connection struct {
	Logger

	mtx sync.Mutex

	network, address string
	nc               net.Conn

	rpcid   uint32
	callers map[uint32]*caller
}

func newConnection(network, address string) *Connection {
	conn := new(Connection)
	conn.network, conn.address = network, address
	conn.callers = make(map[uint32]*caller)
	return conn
}

func (c *Connection) Connect() error {
	nc, err := net.Dial(c.network, c.address)
	if err != nil {
		return err
	}

	c.nc = nc
	return nil
}

func (c *Connection) recv(handler func(p *api.Packet)) {
	nc := c.nc
	for {
		var p api.Packet
		if err := p.Read(nc); err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				continue
			}
			c.Errorf("read %s packet error:%s", c.nc.RemoteAddr(), err)
			break
		}

		if p.Flag&flagRPCResponse != 0 {
			c.mtx.Lock()
			if ctx, ok := c.callers[p.ID]; !ok {
				c.mtx.Unlock()
				c.Errorf("%s unexisted rpc return %d", c.nc.RemoteAddr(), p.ID)
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

func (c *Connection) Close() {
	nc := c.nc
	if nc != nil {
		c.nc = nil
		nc.Close()
	}
}

func (c *Connection) Call(ctx context.Context, p *api.Packet) (*api.Packet, error) {
	caller := newCaller()
	caller.id = atomic.AddUint32(&c.rpcid, 1)
	c.mtx.Lock()
	c.callers[caller.id] = caller
	c.mtx.Unlock()

	p.ID = caller.id
	if err := c.write(p); err != nil {
		return nil, err
	}

	select {
	case p, ok := <-caller.ch:
		if !ok {
			return nil, errors.New("rpc wait fail")
		}
		if p.Ec != 0 {
			return nil, errors.New("response ec")
		}
		return p, nil
	case <-ctx.Done():
		c.mtx.Lock()
		delete(c.callers, caller.id)
		c.mtx.Unlock()
		return nil, errors.New("timeout")
	}
}

func (c *Connection) write(p *api.Packet) error {
	p.Length = uint32(len(p.Data))
	if err := p.Write(c.nc); err != nil {
		return err
	}
	return nil
}
