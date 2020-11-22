package server

import (
	"time"

	apipb "github.com/anson-xcloud/xdp-demo/api"
)

type Remote apipb.Remote
type RemoteSlice []*apipb.Remote
type Data apipb.Data

func IsValidRemote(remote *Remote) bool {
	return remote.Sid != "" || remote.Appid != ""
}

type Request struct {
	*Remote

	*Data

	pid uint32

	reqTime time.Time
}

func (r *Request) GetHeader(key string) string {
	v := r.Data.Headers[key]
	return v
}
