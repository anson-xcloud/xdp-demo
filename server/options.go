package server

import (
	"io/ioutil"
	"os"
)

// Option server option
type Option func(*Options)

// Options server all options
type Options struct {
	Handler2 Handler

	Logger Logger

	// Config xcloud config
	Config string
}

var defaultOptions = Options{
	Handler2: defaultServeMux,
	Logger:   defaultLogger,
}

// WithHandler set handler
// if dont set, default use *ServeMux
func WithHandler(h Handler) Option {
	return func(opts *Options) {
		opts.Handler2 = h
	}
}

// WithLogger set logger
// if dont set, default use *fmtLogger
func WithLogger(l Logger) Option {
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
