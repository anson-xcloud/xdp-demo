package xdp

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
)

type Client struct {
	mtx sync.Mutex

	AppID string

	conn *Connection

	token  string
	cookie *http.Cookie
}

func NewClient(appid string) *Client {
	c := new(Client)
	c.AppID = appid
	return c
}

func (c *Client) Connect() error {
	addr, err := c.getAppAddr()
	if err != nil {
		return err
	}

	cc := newConnection()
	if err = cc.Connect(addr); err != nil {
		return err
	}

	c.conn = cc
	return nil
}

func (c *Client) Send(api string, data []byte) {
	var p Packet
	// p.Cmd=api
	p.Data = data
	c.conn.write(&p)
}

func (c *Client) Get(api string, headers url.Values) ([]byte, error) {
	he := headers.Encode()
	addr := fmt.Sprintf("%s%s%s?%s", XCloudAddr, APIClientXdpPrefix, c.AppID, he)
	resp, err := http.Get(addr)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *Client) Login(user, pwd string) error {
	vals := make(url.Values)
	vals.Add("user", user)
	vals.Add("pwd", pwd)
	addr := fmt.Sprintf("%s%s?%s", XCloudAddr, APIUserLogin, vals.Encode())

	req, err := http.NewRequest(http.MethodGet, addr, nil)
	if err != nil {
		return err
	}
	if c.cookie != nil {
		req.AddCookie(c.cookie)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	cookies := resp.Cookies()
	if len(cookies) > 0 {
		fmt.Println("cookies:::", cookies)
		c.cookie = cookies[0]
	}
	c.token = string(data)
	return nil
}

func (c *Client) getAppAddr() (string, error) {
	addr := fmt.Sprintf("%s%s%s ", XCloudAddr, APIUserGetAccessPoint, c.AppID)
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

func Get(appid, api string, headers url.Values) ([]byte, error) {
	cli := NewClient(appid)
	return cli.Get(api, headers)
}

func Login(appid, user, pwd string) (*Client, error) {
	cli := NewClient(appid)
	if err := cli.Login(user, pwd); err != nil {
		return nil, err
	}
	return cli, nil
}

func Connect(appid string) (*Client, error) {
	cli := NewClient(appid)
	if err := cli.Connect(); err != nil {
		return nil, err
	}

	return cli, nil
}
