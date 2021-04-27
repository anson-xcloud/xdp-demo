package joinpoint

type Config struct {
	Addr string

	Provider Provider
}

func DefaultConfig() *Config {
	c := &Config{}
	return c
}
