package server

import (
	"context"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

const flagRPCResponse = 0x1

var (
	packetbo = binary.LittleEndian
)

// Packet xcloud packet define
// TODO
// adjust like http header
// Proto      string // XDP/1
// ProtoMajor int
// ProtoMinor int
type Packet struct {
	Length uint32
	Flag   uint32
	ID     uint32
	Ec     uint32
	Cmd    uint32
	Data   []byte
}

func (p *Packet) Write(writer io.Writer) error {
	if err := binary.Write(writer, packetbo, &p.Length); err != nil {
		return err
	}
	if err := binary.Write(writer, packetbo, &p.ID); err != nil {
		return err
	}
	if err := binary.Write(writer, packetbo, &p.Flag); err != nil {
		return err
	}
	if err := binary.Write(writer, packetbo, &p.Ec); err != nil {
		return err
	}
	if err := binary.Write(writer, packetbo, &p.Cmd); err != nil {
		return err
	}
	// TODO
	// 暂时忽略n
	_, err := writer.Write(p.Data)
	return err
}

func (p *Packet) Read(reader io.Reader) error {
	if err := binary.Read(reader, packetbo, &p.Length); err != nil {
		return err
	}
	if err := binary.Read(reader, packetbo, &p.ID); err != nil {
		return err
	}
	if err := binary.Read(reader, packetbo, &p.Flag); err != nil {
		return err
	}
	if err := binary.Read(reader, packetbo, &p.Ec); err != nil {
		return err
	}
	if err := binary.Read(reader, packetbo, &p.Cmd); err != nil {
		return err
	}

	p.Data = make([]byte, p.Length)
	_, err := io.ReadFull(reader, p.Data)
	return err
}

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

	cec error //close error

	rpcid   uint32
	callers map[uint32]*caller
}

func newConnection() *Connection {
	conn := new(Connection)
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
	return nil
}

func (c *Connection) recv(handler func(p *Packet)) error {
	nc := c.nc
	for {
		var p Packet
		if err := p.Read(nc); err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				continue
			}

			c.Error("read %s packet error:%s", nc.RemoteAddr(), err)
			return err
		}

		if p.Flag&flagRPCResponse == 0 {
			handler(&p)
			continue
		}

		c.mtx.Lock()
		if ctx, ok := c.callers[p.ID]; !ok {
			c.mtx.Unlock()
			c.Error("%s unexisted rpc return %d", nc.RemoteAddr(), p.ID)
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
	if err := c.write(p); err != nil {
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

func (c *Connection) write(p *Packet) error {
	p.Length = uint32(len(p.Data))
	if err := p.Write(c.nc); err != nil {
		return err
	}
	return nil
}
