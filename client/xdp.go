package client

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/anson-xcloud/xdp-demo/config"
	"github.com/anson-xcloud/xdp-demo/pkg/network"
)

type xdpClient struct {
	mtx sync.Mutex

	appid string

	conn *network.Connection
}

func New() Client {
	return new(xdpClient)
}

func (x *xdpClient) Serve(appid string) error {
	x.appid = appid
	return nil
}

func (x *xdpClient) connect() error {
	addr, err := x.getAppAddr()
	if err != nil {
		return err
	}

	cc := network.NewConnection()
	if err = cc.Connect(addr); err != nil {
		return err
	}

	x.conn = cc
	return nil
}

func (c *xdpClient) Send(api string, data []byte) {
	var p network.Packet
	// p.Cmd=api
	p.Data = data
	c.conn.Write(&p)
}

func (c *xdpClient) Get(req2 *Request) ([]byte, error) {
	return nil, errors.New("unimplement")
}

func (c *xdpClient) Login(user, pwd string) error {
	return errors.New("unimplement")
}

func (c *xdpClient) getAppAddr() (string, error) {
	addr := fmt.Sprintf("%s%s%s ", config.XCloudAddr, config.APIUserGetAccessPoint, c.appid)
	resp, err := http.Get(addr)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
