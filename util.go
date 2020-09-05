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
