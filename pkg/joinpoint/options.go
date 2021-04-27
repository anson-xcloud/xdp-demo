package joinpoint

import "time"

type Option func(*Options)

type Options struct {
	Worker Worker

	MaxConnectTime time.Duration

	MaxHandlerTime time.Duration
}

var defaultOptions = Options{
	Worker: NewGoWorker(),
}
