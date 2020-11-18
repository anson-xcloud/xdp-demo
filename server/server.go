package server

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	apipb "github.com/anson-xcloud/xdp-demo/api"
	"github.com/anson-xcloud/xdp-demo/config"
	"github.com/anson-xcloud/xdp-demo/pkg/logger"
	"github.com/anson-xcloud/xdp-demo/pkg/network"
	"github.com/golang/protobuf/proto"
)

var defaultSvr Server

func init() {
	defaultSvr = NewServer()
}

const (
	XdpGet  = "GET"
	XdpPost = "POST"
)

type Remote apipb.Remote
type RemoteSlice []*apipb.Remote
type Data apipb.Data

func IsValidRemote(remote *Remote) bool {
	return remote.Sid != "" || remote.Appid != ""
}

// Address for app address token
// format is   appid:appsecret
type Address struct {
	AppID, AppSecret string
}

func (a *Address) String() string {
	return fmt.Sprintf("%s:%s", a.AppID, a.AppSecret)
}

// ParseAddress parse address string to *Address
func ParseAddress(addr string) (*Address, error) {
	sl := strings.Split(addr, ":")
	if len(sl) != 2 {
		return nil, ErrAddressFormat
	}

	return &Address{AppID: sl[0], AppSecret: sl[1]}, nil
}

type Request struct {
	*Remote

	*Data

	pid uint32

	reqTime time.Time
}

type Server interface {
	GetLogger() logger.Logger

	GetAddr() *Address

	// Serve block run server until error or shutdown
	Serve(addr string) error

	// Stop()

	// RefreshHostServerAuthority(appid string, setting HostSetting)

	// Reply reply request message data, only valid when method equal XdpGet
	Reply(req *Request, data []byte)

	// ReplyError reply request message data with error, only valid when method equal XdpGet
	ReplyError(req *Request, ec uint32, msg string)

	// Send send data to target, support any target,
	// when send to host user/server, need host setting allow
	Send(remote *Remote, data *Data) error

	// Send send data to multi target, support any target,
	// when send to any host user/server, need host setting allow
	MultiSend(remotes RemoteSlice, data *Data) error

	// Get will wait until return
	// note: you can only get from host server / plugin server
	Get(appid string, data *Data) ([]byte, error)
}

func SetEnv(env string) {
	config.SetEnv(env)
}

func Serve(addr string) error {
	return defaultSvr.Serve(addr)
}

func Send(appid, api string, data []byte) {
	defaultSvr.Send(&Remote{Appid: appid}, &Data{Api: api, Data: data})
}

func Get(appid string, api string, data []byte) ([]byte, error) {
	return defaultSvr.Get(appid, &Data{Api: api, Data: data})
}

func GetLogger() logger.Logger {
	return defaultSvr.GetLogger()
}

func ReplyJson(svr Server, req *Request, data interface{}) error {
	bdata, err := json.Marshal(data)
	if err != nil {
		return err
	}
	svr.Reply(req, bdata)
	return nil
}

// xdpServer for app server
type xdpServer struct {
	sync.RWMutex

	opts *Options

	addr *Address

	// TODO check nil
	conn *network.Connection

	// settings about who host current app as plugin
	hosts map[string]*HostSetting
}

// NewServer create server
func NewServer(opt ...Option) Server {
	var opts = defaultOptions
	for _, o := range opt {
		o(&opts)
	}

	xs := new(xdpServer)
	xs.opts = &opts
	return xs
}

// Logger implement server.GetLogger
func (x *xdpServer) GetLogger() logger.Logger {
	return x.opts.Logger
}

// GetAddr implement server.GetAddr
func (x *xdpServer) GetAddr() *Address {
	return x.addr
}

// Serve start serve at addr
func (x *xdpServer) Serve(addr string) error {
	address, err := ParseAddress(addr)
	if err != nil {
		return err
	}

	x.addr = address
	ap, err := x.getAccessPoint()
	if err != nil {
		return err
	}

	conn := network.NewConnection()
	conn.Logger = x.opts.Logger
	if err := conn.Connect(ap.Addr); err != nil {
		return err
	}

	go func() {
		if _, err = call(conn, apipb.Cmd_CmdHandshake, &apipb.HandshakeRequest{
			AppID:     x.addr.AppID,
			AccessKey: ap.AccessKey,
		}); err != nil {
			conn.Close(err)
			return
		}
		x.conn = conn

		x.opts.Logger.Info("start serve xdp app %s ... ", x.addr.AppID)
	}()
	return conn.Recv(x.process)
}

// Reply implement Server.Reply
func (x *xdpServer) Reply(req *Request, data []byte) {
	if req.pid == 0 {
		return
	}

	var p network.Packet
	p.ID = req.pid
	p.Flag |= network.FlagRPCResponse
	p.Data = data
	x.writePacket(&p)
}

func (x *xdpServer) ReplyError(req *Request, ec uint32, msg string) {
	if req.pid == 0 {
		return
	}

	var p network.Packet
	p.ID = req.pid
	p.Flag |= network.FlagRPCResponse
	p.Ec = ec
	// p.EcMsg=msg
	x.writePacket(&p)
}

func (x *xdpServer) Send(remote *Remote, data *Data) error {
	if !IsValidRemote(remote) {
		return ErrInvalidRemote
	}

	pbs := (*apipb.Remote)(remote)
	if !x.isApiAllow(data.Api, pbs) {
		return ErrApiNowAllowed
	}

	var m apipb.Message
	m.Remote = pbs
	m.Data = (*apipb.Data)(data)
	return x.write(apipb.Cmd_CmdSend, &m)
}

// MultiSend multi send data to session at once
func (x *xdpServer) MultiSend(remotes RemoteSlice, data *Data) error {
	for _, remote := range remotes {
		if !IsValidRemote((*Remote)(remote)) {
			return ErrInvalidRemote
		}
	}

	if !x.isApiAllow(data.Api, remotes...) {
		return ErrApiNowAllowed
	}

	var m apipb.MultiMessage
	m.Remotes = ([]*apipb.Remote)(remotes)
	m.Data = (*apipb.Data)(data)
	return x.write(apipb.Cmd_CmdMultiSend, &m)
}

func (x *xdpServer) Get(appid string, data *Data) ([]byte, error) {
	pbs := &apipb.Remote{Appid: appid}
	if !x.isApiAllow(data.Api, pbs) {
		return nil, ErrApiNowAllowed
	}

	var m apipb.Message
	m.Remote = pbs
	m.Data = (*apipb.Data)(data)
	return x.call(apipb.Cmd_CmdGet, &m)
}

func (x *xdpServer) isApiAllow(api string, remotes ...*apipb.Remote) bool {
	x.RLock()
	defer x.RUnlock()

	for _, remote := range remotes {
		setting, ok := x.hosts[remote.Appid]
		if ok && !setting.isAllow(remote, api) {
			return false
		}
	}
	return true
}

func (x *xdpServer) process(p *network.Packet) {
	go func() {
		switch p.Cmd {
		case uint32(apipb.Cmd_CmdRecv):
			x.processRecv(p)
		default:
			x.GetLogger().Warn("unknown cmd %d", p.Cmd)
		}
	}()
}

func (x *xdpServer) signURL(vals url.Values) {
	md5str := fmt.Sprintf("%s%s", vals.Encode(), x.addr.AppSecret)
	m := md5.New()
	token := hex.EncodeToString(m.Sum([]byte(md5str)))
	vals.Set("token", token)
}

// AccessPoint xcloud return access_point info
type AccessPoint struct {
	Addr      string `json:"addr"`
	AccessKey string `json:"access_key"`
}

func (x *xdpServer) getAccessPoint() (*AccessPoint, error) {
	values := make(url.Values)
	values.Set("appid", x.addr.AppID)
	x.signURL(values)
	url := fmt.Sprintf("%s%s?%s", config.Env.XcloudAddr, config.APIAccessPoint, values.Encode())

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	type AccessPointResult struct {
		Status int    `json:"status"`
		Msg    string `json:"msg"`
		AccessPoint
	}
	var ret AccessPointResult
	if err := json.Unmarshal(data, &ret); err != nil {
		return nil, err
	}
	if ret.Status != 0 {
		return nil, fmt.Errorf("response error(%d):%s", ret.Status, ret.Msg)
	}
	return &ret.AccessPoint, nil
}

func (x *xdpServer) call(cmd apipb.Cmd, pm proto.Message) ([]byte, error) {
	return call(x.conn, cmd, pm)
}

func call(conn *network.Connection, cmd apipb.Cmd, pm proto.Message) ([]byte, error) {
	bs, err := proto.Marshal(pm)
	if err != nil {
		return nil, err
	}

	var p network.Packet
	p.Cmd = uint32(cmd)
	p.Data = bs
	rp, err := conn.Call(context.Background(), &p)
	if err != nil {
		return nil, err
	}
	return rp.Data, nil
}

func (x *xdpServer) write(cmd apipb.Cmd, pm proto.Message) error {
	bs, err := proto.Marshal(pm)
	if err != nil {
		return err
	}

	var p network.Packet
	p.Cmd = uint32(apipb.Cmd_CmdSend)
	p.Data = bs
	return x.writePacket(&p)
}

func (x *xdpServer) writePacket(p *network.Packet) error {
	return x.conn.Write(p)
}

func (x *xdpServer) processRecv(p *network.Packet) {
	var notify apipb.Message
	if err := proto.Unmarshal(p.Data, &notify); err != nil {
		x.opts.Logger.Debug("unmarshal handleData error:%s", err)
		return
	}

	var req Request
	req.Remote = (*Remote)(notify.Remote)
	req.Data = (*Data)(notify.Data)
	req.reqTime = time.Now()
	req.pid = p.ID
	x.opts.Handler.Serve(x, &req)
}
