package xcloud

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
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
	Rid int // runtime id

	env EnvConfig

	logger xlog.Logger

	sm *ServeMux
}

func New(c *Config) *XCloud {
	xc := &XCloud{logger: c.Logger}

	xc.env = c.Env
	xc.sm = defaultServeMux

	return xc
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
	// x.ID = ap.ID

	conn := network.NewConnection()
	// conn.Logger = x.opts.Logger
	if err := conn.Connect(ap.Addr); err != nil {
		return nil, nil, err
	}

	data, err := call(conn, "serivce.register", &apipb.ServiceRegisterRequest{
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
	// x.conn = conn
	return &Transport{conn: conn}, nil, nil
	// x.opts.Logger.Info("start serve xdp app %s(%d) ... ", x.addr.AppID, x.Rid)

	// return conn.Recv(x.process)
}

func (x *XCloud) Serve(ctx context.Context, rw joinpoint.ResponseWriter, jr joinpoint.Request) {
	req := jr.(*Request)

	req.rw = rw
	x.sm.Serve(ctx, req)
	// x.opts.Handler.Serve(x, &req)
}

type Transport struct {
	conn *network.Connection
}

func (t *Transport) Recv(ctx context.Context) (joinpoint.Request, error) {
	var p *network.Packet
	select {
	case p = <-t.conn.Recv2():
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

func call(conn *network.Connection, cmd string, pm proto.Message) ([]byte, error) {
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
	rp, err := conn.Call(context.Background(), &np)
	if err != nil {
		return nil, err
	}
	return rp.Data, nil
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
	url := fmt.Sprintf("%s%s?%s", x.env.XcloudAddr, APIAccessPoint, values.Encode())

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http code: %d", resp.StatusCode)
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
