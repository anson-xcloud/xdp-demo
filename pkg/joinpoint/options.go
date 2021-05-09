package joinpoint

import (
	"time"

	"github.com/anson-xcloud/xdp-demo/pkg/xlog"
)

type Option func(*Options)

type Options struct {
	worker Worker

	MaxConnectTime time.Duration

	MaxHandlerTime time.Duration

	logger xlog.Logger
}

var defaultOptions = Options{
	worker: NewGoWorker(func(v interface{}) {
		xlog.Errorf("%v", v)
	}),
	logger: xlog.Default,
}

func WithLogger(logger xlog.Logger) Option {
	return func(opts *Options) {
		opts.logger = logger
	}
}
