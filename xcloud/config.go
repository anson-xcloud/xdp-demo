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

	Handler Handler

	Agent bool
}

func DefaultConfig() *Config {
	return &Config{
		Env:     EnvDefault,
		Logger:  xlog.Default,
		Handler: defaultServeMux,
	}
}

type EnvConfig struct {
	Discription string

	XcloudAddrs []string
}

const (
	EnvDevDiscription     = "dev"
	EnvDebugDiscription   = "debug"
	EnvReleaseDiscription = "release"
)

var (
	enves = map[string]EnvConfig{
		EnvDevDiscription: {
			Discription: EnvDevDiscription,
			XcloudAddrs: []string{"http://localhost:4000"},
		},
		EnvDebugDiscription: {
			Discription: EnvDebugDiscription,
			XcloudAddrs: []string{"http://localhost:4000"},
		},
		EnvReleaseDiscription: {
			Discription: EnvReleaseDiscription,
			XcloudAddrs: []string{"http://xcloud.singularityfuture.com.cn"},
		},
	}

	EnvDev     = GetEnv(EnvDevDiscription)
	EnvDebug   = GetEnv(EnvDebugDiscription)
	EnvRelease = GetEnv(EnvReleaseDiscription)
	EnvDefault = EnvRelease
)

func GetEnv(env string) EnvConfig {
	c, ok := enves[env]
	// debug.PanicIf(ok, "")
	if !ok {
		panic("invalid env")
	}

	return c
}

func SetEnv(env string) {
	EnvDefault = GetEnv(env)
}
