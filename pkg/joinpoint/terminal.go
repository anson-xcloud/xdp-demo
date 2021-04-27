package joinpoint

import (
	"context"
	"time"

	"github.com/anson-xcloud/xdp-demo/pkg/xlog"
	"golang.org/x/sync/errgroup"
)

type Terminal struct {
	Provider Provider

	Opts *Options

	connect func(context.Context, string) (Point, []string, error)
}

func Join(ctx context.Context, c *Config, opt ...Option) error {
	if c.Provider == nil {
		return ErrProviderNeed
	}

	var opts = defaultOptions
	for _, o := range opt {
		o(&opts)
	}

	var t Terminal
	t.Provider = c.Provider
	t.Opts = &opts
	t.connect = func(ctx context.Context, addr string) (Point, []string, error) {
		var cctx = ctx
		if opts.MaxConnectTime != 0 {
			var cancel context.CancelFunc
			cctx, cancel = context.WithTimeout(ctx, opts.MaxConnectTime)
			defer cancel()
		}
		return c.Provider.Connect(cctx, addr)
	}
	return t.JoinWithRetry(ctx, c.Addr)
}

func (t *Terminal) JoinWithRetry(ctx context.Context, addr string) error {
	var p Point
	var addrs []string
	var err error
	var nextRetry, maxNextRetry = time.Second, time.Second * 30
	for {
		if p, addrs, err = t.connect(ctx, addr); err == nil {
			break
		}

		if nextRetry = nextRetry * 2; nextRetry > maxNextRetry {
			nextRetry = maxNextRetry
		}
		xlog.Warnf("terminal connect %s fail, wait for retry. error: %s", addr, err)

		select {
		case <-time.After(nextRetry):
		case <-ctx.Done():
			return ErrDone
		}
	}

	go func() {
		eg, egCtx := errgroup.WithContext(ctx)
		eg.Go(func() error { return t.read(egCtx, p, t.Opts.Worker) })
		if err := eg.Wait(); err == nil {
			return
		}

		if err := t.JoinWithRetry(ctx, addr); err != nil {
			for _, addr := range addrs {
				if err := t.JoinWithRetry(ctx, addr); err == nil {
					break
				}
			}
		}
	}()
	return err
}

type responseWriter struct {
	rw ResponseWriter

	cancel context.CancelFunc
}

func (r *responseWriter) WithTimeout(ctx context.Context, du time.Duration) context.Context {
	ctx, cancel := context.WithTimeout(ctx, du)
	go func() {
		<-ctx.Done()
		cancel()
	}()

	r.cancel = cancel
	return ctx
}

func (r *responseWriter) Write(data interface{}) {
	if r.rw != nil {
		r.rw.Write(data)
	}
	if r.cancel != nil {
		r.cancel()
	}
}

func (r *responseWriter) WriteStatus(st *Status) {
	if r.rw != nil {
		r.rw.WriteStatus(st)
	}
	if r.cancel != nil {
		r.cancel()
	}
}

func (t *Terminal) read(ctx context.Context, p Point, worker Worker) error {
	for {
		req, err := p.Recv(ctx)
		if err != nil {
			return err
		}

		worker.Run(func() {
			var rw responseWriter
			rw.rw = req.GetResponseWriter()
			if t.Opts.MaxHandlerTime != 0 {
				ctx = rw.WithTimeout(ctx, t.Opts.MaxHandlerTime)
			}
			t.Provider.Serve(ctx, &rw, req)
		})
	}
}