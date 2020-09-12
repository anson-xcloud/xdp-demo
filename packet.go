package xdp

import (
	"encoding/binary"
	"io"
)

var packetbo = binary.LittleEndian

// TODO
// adjust like http header
// Proto      string // XDP/1
// ProtoMajor int
// ProtoMinor int
type Packet struct {
	Length uint32
	ID     uint32
	Flag   uint32
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
