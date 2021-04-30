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

	apipb "github.com/anson-xcloud/xdp-demo/api"
	"github.com/anson-xcloud/xdp-demo/pkg/joinpoint"
	"github.com/anson-xcloud/xdp-demo/pkg/network"
	"github.com/anson-xcloud/xdp-demo/pkg/xlog"
	"github.com/anson-xcloud/xdp-demo/statuscode"
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
	data, err := transport.call(ctx, "serivce.register", &apipb.ServiceRegisterRequest{
		Id:    ap.ID,
		Rid:   int32(x.Rid),
		Token: ap.Token,
		// Config: x.opts.Config,
	})
	if err != nil {
		conn.Close(err)
		return nil, nil, err
	}

	var resp apipb.ServiceRegisterResponse
	if err := proto.Unmarshal(data, &resp); err != nil {
		conn.Close(err)
		return nil, nil, err
	}
	x.Rid = int(resp.Rid)
	x.transport = transport
	return x.transport, nil, nil
}

func (x *XCloud) Serve(ctx context.Context, rw joinpoint.ResponseWriter, jr joinpoint.Request) {
	req := jr.(*Request)

	h := x.sm.Get(req)
	if h == nil {
		rw.WriteStatus(joinpoint.NewStatus(100, ""))
		return
	}
	h.Serve(ctx, rw, req)
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

	var notify apipb.Message
	if err := proto.Unmarshal(p.Data, &notify); err != nil {
		return nil, err
	}

	var req Request
	req.t = t
	req.Remote = (*Remote)(notify.Remote)
	req.Data = (*Data)(notify.Data)
	req.reqTime = time.Now()
	req.pid = p.ID
	if req.Remote == nil {
		req.Remote = &Remote{}
	}
	if req.Data == nil {
		req.Data = &Data{}
	}
	return &req, nil
}

func (t *Transport) writePacket(packet *network.Packet) error {
	return t.conn.Write(packet)
}

func pack(cmd string, pm proto.Message) (*network.Packet, error) {
	bs, err := proto.Marshal(pm)
	if err != nil {
		return nil, err
	}

	var p apipb.Packet
	p.Cmd = cmd
	// p.Version=1
	p.Data = bs
	pbs, err := proto.Marshal(&p)
	if err != nil {
		return nil, err
	}

	var np network.Packet
	// p.Cmd = uint32(cmd)
	np.Data = pbs
	return &np, nil
}

func (t *Transport) call(ctx context.Context, cmd string, pm proto.Message) ([]byte, error) {
	np, err := pack(cmd, pm)
	if err != nil {
		return nil, err
	}

	rp, err := t.conn.Call(ctx, np)
	if err != nil {
		return nil, err
	}

	var rd apipb.RawData
	if err := proto.Unmarshal(rp.Data, &rd); err != nil {
		return nil, err
	}
	return rd.Data, nil
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

func (x *XCloud) Post(ctx context.Context, remote *Remote, data *Data) error {
	if !IsValidRemote(remote) {
		return ErrInvalidRemote
	}

	pbs := (*apipb.Remote)(remote)
	// if !x.isApiAllow(data.Api, pbs) {
	// 	return ErrApiNowAllowed
	// }

	return x.transport.write(ctx, "xdp.post", &apipb.Message{
		Remote: pbs,
		Data:   (*apipb.Data)(data),
	})
}

// MultiPost multi send data to session at once
func (x *XCloud) MultiPost(ctx context.Context, remotes RemoteSlice, data *Data) error {
	for _, remote := range remotes {
		if !IsValidRemote((*Remote)(remote)) {
			return ErrInvalidRemote
		}
	}

	// if !x.isApiAllow(data.Api, remotes...) {
	// 	return ErrApiNowAllowed
	// }

	return x.transport.write(ctx, "xdp.multipost", &apipb.MultiMessage{
		Remotes: ([]*apipb.Remote)(remotes),
		Data:    (*apipb.Data)(data),
	})
}

func (x *XCloud) Get(ctx context.Context, appid string, data *apipb.Data) ([]byte, error) {
	remote := &apipb.Remote{Appid: appid}
	// if !x.isApiAllow(data.Api, pbs) {
	// 	return nil, ErrApiNowAllowed
	// }

	return x.transport.call(ctx, "xdp.get", &apipb.Message{
		Remote: remote,
		Data:   data,
		// Config: x.opts.Config,
	})
}
