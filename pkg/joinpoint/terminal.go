package joinpoint

import (
	"context"
	"time"

	"github.com/anson-xcloud/xdp-demo/pkg/xlog"
	"golang.org/x/sync/errgroup"
)

type Terminal struct {
	logger xlog.Logger

	Provider Provider

	Opts *Options

	connect func(context.Context, string) (Transport, []string, error)
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
	t.logger = opts.logger
	t.Provider = c.Provider
	t.Opts = &opts
	t.connect = func(ctx context.Context, addr string) (Transport, []string, error) {
		var cctx = ctx
		if opts.MaxConnectTime != 0 {
			var cancel context.CancelFunc
			cctx, cancel = context.WithTimeout(ctx, opts.MaxConnectTime)
			defer cancel()
		}
		return c.Provider.Connect(cctx, addr)
	}
	return t.JoinWithRetry(ctx, []string{c.ServerAddr})
}

func (t *Terminal) JoinWithRetry(ctx context.Context, addrs []string) error {
	if len(addrs) == 0 {
		return ErrNoneServerAddr
	}

	var addr string
	var nextRetry, maxNextRetry = time.Second, time.Second * 30
	for i := 0; ; i = (i + 1) % len(addrs) {
		addr = addrs[i]
		p, backups, err := t.connect(ctx, addr)
		if err != nil {
			if nextRetry = nextRetry * 2; nextRetry > maxNextRetry {
				nextRetry = maxNextRetry
			}
			t.logger.Warnf("[JOINPOINT] terminal wait for retry, connect %s error: %s", addr, err)

			select {
			case <-time.After(nextRetry):
			case <-ctx.Done():
				return ErrDone
			}
			continue
		}

		t.logger.Debugf("[JOINPOINT] terminal connect %s success", addr)

		go func() {
			start := time.Now()

			eg, egCtx := errgroup.WithContext(ctx)
			eg.Go(func() error { return t.read(egCtx, p, t.Opts.worker) })
			if err := eg.Wait(); err != nil {
				t.logger.Warnf("[JOINPOINT] terminal wait for retry, read %s error: %s", addr, err)

				if nextRetry = nextRetry - time.Since(start); nextRetry > 0 {
					select {
					case <-time.After(nextRetry):
					case <-ctx.Done():
						return
					}
				}
				addrs = []string{addr}
				addrs = append(addrs, backups...)
				if err := t.JoinWithRetry(ctx, addrs); err != nil {
					t.logger.Warnf("[JOINPOINT] terminal join error: %s", err)
				}
			}
		}()
		return nil
	}
}

func (t *Terminal) read(ctx context.Context, p Transport, worker Worker) error {
	for {
		req, err := p.Recv(ctx)
		if err != nil {
			return err
		}

		worker.Run(func() {
			st := time.Now()
			defer func() {
				t.logger.Debugf("[JOINPOINT] terminal serve %s cost %.3fs", req.Discription(), time.Since(st).Seconds())
			}()

			rw := req.GetResponseWriter()
			if t.Opts.MaxHandlerTime != 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, t.Opts.MaxHandlerTime)
				defer cancel()
			}
			t.Provider.Serve(ctx, rw, req)
		})
	}
}
