package joinpoint

import (
	"github.com/anson-xcloud/xdp-demo/pkg/xflag"
)

type Config struct {
	ServerAddr string

	Provider Provider
}

func DefaultConfig() *Config {
	c := &Config{}
	return c
}

func (c *Config) Flags() *xflag.FlagSet {
	fs := xflag.NewFlagSet()
	fs.StringVar(&c.ServerAddr, "server_addr", "", "joinpoint server addr for connect")
	return fs
}
