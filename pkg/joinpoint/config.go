package joinpoint

import "github.com/anson-xcloud/xdp-demo/pkg/xlog"

type Config struct {
	Addr string

	Provider Provider

	Logger xlog.Logger
}

func DefaultConfig() *Config {
	c := &Config{}
	c.Logger = xlog.Default
	return c
}
