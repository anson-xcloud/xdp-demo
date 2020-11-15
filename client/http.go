package client

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"

	"github.com/anson-xcloud/xdp-demo/config"
)

const (
	headerPrefix = "Xcloud-"
	headerAppid  = headerPrefix + "Appid"
)

type httpClient struct {
	sync.RWMutex

	appid string

	token string

	cookie *http.Cookie
}

func NewHttpClient() Client {
	return new(httpClient)
}

func (h *httpClient) Connect(appid string) error {
	h.appid = appid
	return nil
}

func (h *httpClient) Get(req *Request) ([]byte, error) {
	addr := fmt.Sprintf("%s%s%s", config.XCloudAddr, config.APIClientXdpPrefix, req.Appid)
	httpReq, err := http.NewRequest(http.MethodPost, addr, bytes.NewBuffer(req.Data))
	if err != nil {
		return nil, err
	}
	req.Headers[headerAppid] = h.appid
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}
	h.RLock()
	if h.cookie != nil {
		httpReq.AddCookie(h.cookie)
	}
	h.RUnlock()

	resp, err := http.DefaultClient.Do(httpReq)
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

func (h *httpClient) Post(req *Request) error {
	return errors.New("http client donot support Post")
}

func (h *httpClient) Login(user, pwd string) error {
	vals := make(url.Values)
	vals.Add("user", user)
	vals.Add("pwd", pwd)
	addr := fmt.Sprintf("%s%s?%s", config.XCloudAddr, config.APIUserLogin, vals.Encode())
	req, err := http.NewRequest(http.MethodGet, addr, nil)
	if err != nil {
		return err
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

	h.Lock()
	defer h.Unlock()
	if len(cookies) > 0 {
		h.cookie = cookies[0]
	}
	h.token = string(data)
	return nil
}
