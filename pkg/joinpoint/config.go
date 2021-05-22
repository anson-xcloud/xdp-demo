package joinpoint

type Config struct {
	ServerAddr string `xconf:"server_addr" xconfusage:"joinpoint server addr for connect"`

	Provider Provider `xconf:"-"`
}

func DefaultConfig() *Config {
	c := &Config{}
	return c
}
