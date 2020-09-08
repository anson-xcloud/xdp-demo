package xdp

import (
	"encoding/binary"
	"io"
)

func writeString(writer io.Writer, s string) (n int, err error) {
	c := uint32(len(s))
	if err = binary.Write(writer, endian, &c); err != nil {
		return
	}
	n += 4

	if _, err = writer.Write([]byte(s)); err != nil {
		return
	}
	n += len(s)

	return
}

func readString(r io.Reader) (string, error) {
	var c uint32
	if err := binary.Read(r, endian, &c); err != nil {
		return "", err
	}
	if c == 0 {
		return "", nil
	}

	buf := make([]byte, c)
	if _, err := r.Read(buf); err != nil {
		return "", err
	}

	return string(buf), nil
}
