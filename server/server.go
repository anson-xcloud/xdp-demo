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
	"strconv"
	"strings"
	"sync"
	"time"

	apipb "github.com/anson-xcloud/xdp-demo/api"
	"github.com/anson-xcloud/xdp-demo/config"
	"github.com/anson-xcloud/xdp-demo/pkg/logger"
	"github.com/anson-xcloud/xdp-demo/pkg/network"
	"github.com/anson-xcloud/xdp-demo/statuscode"
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

type Server interface {
	GetLogger() logger.Logger

	GetAddr() *Address

	// Serve block run server until error or shutdown
	Serve(addr string, opt ...Option) error

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

func SetEnvDebug() {
	config.SetEnv(config.EnvDebugDiscription)
}

func SetEnvDev() {
	config.SetEnv(config.EnvDevDiscription)
}

func SetEnvRelease() {
	config.SetEnv(config.EnvReleaseDiscription)
}

func Serve(addr string, opt ...Option) error {
	return defaultSvr.Serve(addr, opt...)
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

	ID  string
	Rid int

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
	xs.Rid = opts.Rid
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
func (x *xdpServer) Serve(addr string, opt ...Option) error {
	for _, o := range opt {
		o(x.opts)
	}

	address, err := ParseAddress(addr)
	if err != nil {
		return err
	}

	x.addr = address

	for {
		if err = x.run(); err != nil {
			if x.opts.OnceTry {
				return err
			}

			x.GetLogger().Error("run error: %s, retry after 10s", err)
			time.Sleep(time.Second * 10)
		}
	}
}

func (x *xdpServer) run() error {
	ap, err := x.getAccessPoint()
	if err != nil {
		return err
	}
	x.ID = ap.ID

	conn := network.NewConnection()
	conn.Logger = x.opts.Logger
	if err := conn.Connect(ap.Addr); err != nil {
		return err
	}

	go func() {
		data, err := call(conn, "serivce.register", &apipb.ServiceRegisterRequest{
			Id:     ap.ID,
			Rid:    int32(x.Rid),
			Token:  ap.Token,
			Config: x.opts.Config,
		})
		if err != nil {
			conn.Close(err)
			return
		}
		var resp apipb.ServiceRegisterResponse
		if err := proto.Unmarshal(data, &resp); err != nil {
			conn.Close(err)
			return
		}
		x.Rid = int(resp.Rid)
		x.conn = conn

		x.opts.Logger.Info("start serve xdp app %s(%d) ... ", x.addr.AppID, x.Rid)
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
	return x.write("xdp.send", &m)
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
	return x.write("xdp.multisend", &m)
}

func (x *xdpServer) Get(appid string, data *Data) ([]byte, error) {
	pbs := &apipb.Remote{Appid: appid}
	if !x.isApiAllow(data.Api, pbs) {
		return nil, ErrApiNowAllowed
	}

	var m apipb.Message
	m.Remote = pbs
	m.Data = (*apipb.Data)(data)
	return x.call("xdp.get", &m)
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
		// case uint32(apipb.Cmd_CmdRecv):
		// 	x.processRecv(p)
		default:
			x.processRecv(p)
			// x.GetLogger().Warn("unknown cmd %d", p.Cmd)
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
	ID    string `json:"id"`
	Addr  string `json:"addr"`
	Token string `json:"token"`
}

func (x *xdpServer) getAccessPoint() (*AccessPoint, error) {
	values := make(url.Values)
	values.Set("appid", x.addr.AppID)
	values.Set("timestamp", strconv.FormatInt(time.Now().Unix(), 10))
	x.signURL(values)
	url := fmt.Sprintf("%s%s?%s", config.Env.XcloudAddr, config.APIAccessPoint, values.Encode())

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

func (x *xdpServer) call(cmd string, pm proto.Message) ([]byte, error) {
	if x.conn == nil {
		return nil, ErrUnprepared
	}
	return call(x.conn, cmd, pm)
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

func (x *xdpServer) write(cmd string, pm proto.Message) error {
	bs, err := proto.Marshal(pm)
	if err != nil {
		return err
	}

	var p apipb.Packet
	p.Cmd = cmd
	// p.Version=1
	p.Data = bs
	bs, err = proto.Marshal(&p)
	if err != nil {
		return err
	}

	var np network.Packet
	np.Data = bs
	return x.writePacket(&np)
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
	if req.Remote == nil {
		req.Remote = &Remote{}
	}
	if req.Data == nil {
		req.Data = &Data{}
	}
	x.opts.Handler.Serve(x, &req)
}
