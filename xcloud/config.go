package xcloud

import (
	"fmt"
	"strings"

	"github.com/anson-xcloud/xdp-demo/pkg/xlog"
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

type Config struct {
	Address Address

	Env EnvConfig

	Logger xlog.Logger
}

func DefaultConfig() *Config {
	return &Config{
		Env:    Env(EnvReleaseDiscription),
		Logger: xlog.Default,
	}
}

type EnvConfig struct {
	Discription string

	XcloudAddr string
}

const (
	EnvDevDiscription     = "dev"
	EnvDebugDiscription   = "debug"
	EnvReleaseDiscription = "release"
)

var enves = map[string]EnvConfig{
	EnvDevDiscription: {
		Discription: EnvDevDiscription,
		XcloudAddr:  "http://localhost:31181",
	},
	EnvDebugDiscription: {
		Discription: EnvDebugDiscription,
		XcloudAddr:  "http://localhost:31181",
	},
	EnvReleaseDiscription: {
		Discription: EnvReleaseDiscription,
		XcloudAddr:  "http://xcloud.singularityfuture.com.cn",
	},
}

func Env(env string) EnvConfig {
	c, ok := enves[env]
	// debug.PanicIf(ok, "")
	if !ok {
		panic("invalid env")
	}

	return c
}
