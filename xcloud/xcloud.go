package xcloud

import (
	"container/list"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/anson-xcloud/xdp-demo/pkg/joinpoint"
	"github.com/anson-xcloud/xdp-demo/pkg/network"
	"github.com/anson-xcloud/xdp-demo/pkg/xlog"
	"github.com/anson-xcloud/xdp-demo/statuscode"
	"github.com/anson-xcloud/xdp-demo/xcloud/apis"
	"google.golang.org/protobuf/proto"
)

type XCloud struct {
	serverAddrs *list.List

	logger xlog.Logger

	sm *ServeMux

	transport *Transport

	Rid int // runtime id
}

func Default() *XCloud {
	xc, _ := New(DefaultConfig())
	return xc
}

func New(c *Config) (*XCloud, error) {
	if len(c.Env.XcloudAddrs) == 0 {
		return nil, errors.New("no server support")
	}

	addrs := list.New()
	for _, addr := range c.Env.XcloudAddrs {
		addrs.PushBack(addr)
	}

	return &XCloud{
		serverAddrs: addrs,
		logger:      c.Logger,
		sm:          c.Handler,
	}, nil
}

func (x *XCloud) Connect(ctx context.Context, addr string) (joinpoint.Transport, []string, error) {
	xaddr, err := ParseAddress(addr)
	if err != nil {
		return nil, nil, err
	}

	ap, err := x.getAccessPoint(xaddr)
	if err != nil {
		return nil, nil, err
	}

	conn := network.NewConnection()
	conn.Logger = x.logger
	if err := conn.Connect(ap.Addr); err != nil {
		return nil, nil, err
	}

	transport := &Transport{conn: conn}
	cresp, err := transport.call(ctx, "serivce.register", &apis.ServiceRegisterRequest{
		Id:    ap.ID,
		Rid:   int32(x.Rid),
		Token: ap.Token,
		// Config: x.opts.Config,
	})
	if err != nil {
		conn.Close(err)
		return nil, nil, err
	}

	var resp apis.ServiceRegisterResponse
	if err := proto.Unmarshal(cresp.Body, &resp); err != nil {
		conn.Close(err)
		return nil, nil, err
	}
	x.Rid = int(resp.Rid)
	x.transport = transport
	return x.transport, nil, nil
}

func (x *XCloud) Serve(ctx context.Context, rw joinpoint.ResponseWriter, jr joinpoint.Request) {
	req := jr.(*Request)
	req.rw = rw

	h := x.sm.Get(req)
	if h == nil {
		rw.WriteStatus(joinpoint.NewStatus(100, ""))
		return
	}
	h.Serve(ctx, req)
}

type Transport struct {
	conn *network.Connection
}

func (t *Transport) Recv(ctx context.Context) (joinpoint.Request, error) {
	var p *network.Packet
	select {
	case p = <-t.conn.Recv2():
	case err := <-t.conn.Error():
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	var msg apis.Request
	if err := proto.Unmarshal(p.Data, &msg); err != nil {
		return nil, err
	}

	var req Request
	req.Request = &msg
	req.t = t
	req.reqTime = time.Now()
	req.pid = p.ID
	if req.Source == nil {
		req.Source = &apis.Peer{}
	}
	return &req, nil
}

func (t *Transport) writePacket(packet *network.Packet) error {
	return t.conn.Write(packet)
}

func pack(api string, pm proto.Message) (*network.Packet, error) {
	bs, err := proto.Marshal(pm)
	if err != nil {
		return nil, err
	}

	var p apis.Request
	p.Api = api
	p.Body = bs
	pbs, err := proto.Marshal(&p)
	if err != nil {
		return nil, err
	}

	var np network.Packet
	np.Data = pbs
	return &np, nil
}

func (t *Transport) call(ctx context.Context, cmd string, pm proto.Message) (*apis.Response, error) {
	np, err := pack(cmd, pm)
	if err != nil {
		return nil, err
	}

	rp, err := t.conn.Call(ctx, np)
	if err != nil {
		return nil, err
	}

	var resp apis.Response
	if err := proto.Unmarshal(rp.Data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (t *Transport) write(ctx context.Context, cmd string, pm proto.Message) error {
	np, err := pack(cmd, pm)
	if err != nil {
		return err
	}

	return t.writePacket(np)
}

func signURL(sec string, vals url.Values) {
	md5str := fmt.Sprintf("%s%s", vals.Encode(), sec)
	m := md5.New()
	token := hex.EncodeToString(m.Sum([]byte(md5str)))
	vals.Set("token", token)
}

// AccessPoint xcloud return access_point info
type AccessPoint struct {
	ID    string `json:"id"`
	Addr  string `json:"addr"`
	Token string `json:"token"`
}

func (x *XCloud) getAccessPoint(addr *Address) (*AccessPoint, error) {
	values := make(url.Values)
	values.Set("appid", addr.AppID)
	values.Set("timestamp", strconv.FormatInt(time.Now().Unix(), 10))
	signURL(addr.AppSecret, values)

	it := x.serverAddrs.Front()
	x.serverAddrs.Remove(it)
	x.serverAddrs.PushBack(it.Value)
	url := fmt.Sprintf("%s%s?%s", it.Value, APIAccessPoint, values.Encode())

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("getAccessPoint errcode %d", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ret struct {
		statuscode.Response
		AccessPoint
	}
	if err := json.Unmarshal(data, &ret); err != nil {
		return nil, err
	}
	if ret.Code != statuscode.CodeOK {
		return nil, &ret.Response
	}
	return &ret.AccessPoint, nil
}

func (x *XCloud) Post(ctx context.Context, target *apis.Peer, api string, headers map[string]string, body []byte) error {
	// if !x.isApiAllow(data.Api, pbs) {
	// 	return ErrApiNowAllowed
	// }

	return x.transport.write(ctx, "xdp.post", &apis.Request{
		Api:     api,
		Target:  target,
		Headers: headers,
		Body:    body,
	})
}

// MultiPost multi send data to session at once
// func (x *XCloud) MultiPost(ctx context.Context, remotes RemoteSlice, data *Data) error {
// 	for _, remote := range remotes {
// 		if !isValidRemote((*Remote)(remote)) {
// 			return ErrInvalidRemote
// 		}
// 	}

// 	// if !x.isApiAllow(data.Api, remotes...) {
// 	// 	return ErrApiNowAllowed
// 	// }

// 	return x.transport.write(ctx, "xdp.multipost", &apis.MultiMessage{
// 		Remotes: ([]*apis.Remote)(remotes),
// 		Data:    (*apis.Data)(data),
// 	})
// }

func (x *XCloud) Get(ctx context.Context, appid, api string, headers map[string]string, body []byte) ([]byte, error) {
	// if !x.isApiAllow(data.Api, pbs) {
	// 	return nil, ErrApiNowAllowed
	// }

	resp, err := x.transport.call(ctx, "xdp.get", &apis.Request{
		Api:     api,
		Target:  &apis.Peer{Appid: appid},
		Headers: headers,
		Body:    body,
	})
	if err != nil {
		return nil, err
	}
	return resp.Body, err
}
