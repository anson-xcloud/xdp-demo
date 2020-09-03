package xdp

import (
	"io"
)

type Request struct {
	Proto      string // XDP/1
	ProtoMajor int
	ProtoMinor int

	Cmd int

	Plugin string

	Body io.Reader
}
