package joinpoint

import (
	"xcloud/pkg/xflag"

	"github.com/anson-xcloud/xdp-demo/pkg/xlog"
)

type Config struct {
	ServerAddr string

	Provider Provider

	Logger xlog.Logger
}

func DefaultConfig() *Config {
	c := &Config{}
	c.Logger = xlog.Default
	return c
}

func (c *Config) Flags() *xflag.FlagSet {
	fs := xflag.NewFlagSet()
	fs.StringVar(&c.ServerAddr, "server_addr", "", "joinpoint server addr for connect")
	return fs
}
