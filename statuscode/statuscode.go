package statuscode

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const (
	HeaderXcloudPrefix = "X-Xcloud-"

	HeaderXcloudAppid       = HeaderXcloudPrefix + "Appid"
	HeaderXcloudStatus      = HeaderXcloudPrefix + "Status"
	HeaderXcloudStatusApp   = HeaderXcloudPrefix + "Status-App"
	HeaderXcloudMessage     = HeaderXcloudPrefix + "Message"
	HeaderXcloudDiscription = HeaderXcloudPrefix + "Discription"
)

const CodeOK = 0

type Response struct {
	Code        int    `json:"code"`
	AppCode     int    `json:"appcode,omitempty"`
	Message     string `json:"msg"`
	Discription string `json:"disc"`
}

func (r *Response) Error() string {
	var arr []string
	arr = append(arr, strconv.Itoa(r.Code))
	if r.Message != "" {
		msg := r.Message
		if r.AppCode != 0 {
			msg = fmt.Sprintf("%s(%d)", r.Message, r.AppCode)
		}
		arr = append(arr, msg)
	}
	if r.Discription != "" {
		arr = append(arr, r.Discription)
	}
	return strings.Join(arr, ",")
}

func Get(url string, obj interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		xs := resp.Header.Get(HeaderXcloudStatus)
		if xs != "" {
			var rerr Response
			rerr.Code, _ = strconv.Atoi(xs)
			rerr.AppCode, _ = strconv.Atoi(resp.Header.Get(HeaderXcloudStatusApp))
			rerr.Message = resp.Header.Get(HeaderXcloudMessage)
			rerr.Discription = resp.Header.Get(HeaderXcloudDiscription)
			return &rerr
		}
		return fmt.Errorf("http code: %d", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, obj)
}
