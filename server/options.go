package server

import (
	"io/ioutil"
	"os"

	"github.com/anson-xcloud/xdp-demo/pkg/logger"
)

// Option server option
type Option func(*Options)

// Options server all options
type Options struct {
	Rid int

	Handler Handler

	Logger logger.Logger

	// Config xcloud config
	Config string

	OnceTry bool
}

var defaultOptions = Options{
	Handler: defaultServeMux,
	Logger:  logger.Default,
}

// WithHandler set handler
// if dont set, default use *ServeMux
func WithHandler(h Handler) Option {
	return func(opts *Options) {
		opts.Handler = h
	}
}

func WithRid(rid int) Option {
	return func(opts *Options) {
		opts.Rid = rid
	}
}

// WithLogger set logger
// if dont set, default use *fmtLogger
func WithLogger(l logger.Logger) Option {
	return func(opts *Options) {
		opts.Logger = l
	}
}

// WithConfig set config
func WithConfig(cfg string) Option {
	return func(opts *Options) {
		opts.Config = cfg
	}
}

func WithOnceTry(once bool) Option {
	return func(opts *Options) {
		opts.OnceTry = once
	}
}

// WithConfigFile set config by file path,
// any read error will panic
func WithConfigFile(fp string) Option {
	return func(opts *Options) {
		f, err := os.Open(fp)
		if err != nil {
			panic(err)
		}

		data, err := ioutil.ReadAll(f)
		if err != nil {
			panic(err)
		}

		opts.Config = string(data)
	}
}
