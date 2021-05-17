package joinpoint

import (
	"context"
	"time"

	"github.com/anson-xcloud/xdp-demo/pkg/xlog"
	"golang.org/x/sync/errgroup"
)

type Terminal struct {
	logger xlog.Logger

	ctx    context.Context
	cancel context.CancelFunc

	Provider Provider

	Opts *Options

	connect func(context.Context, string) (Transport, []string, error)
}

func Join(ctx context.Context, c *Config, opt ...Option) (*Terminal, error) {
	if c.Provider == nil {
		return nil, ErrProviderNeed
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
	t.ctx, t.cancel = context.WithCancel(ctx)

	go t.joinWithRetry(t.ctx, []string{c.ServerAddr}, 0)
	return &t, nil
}

const maxNextRetry = time.Second * 30

func (t *Terminal) Context() context.Context {
	return t.ctx
}

func (t *Terminal) joinWithRetry(ctx context.Context, addrs []string, nextRetry time.Duration) {
	if len(addrs) == 0 {
		return
	}

	fnNextRetry := func(nr time.Duration) time.Duration {
		if nr = nr*2 + time.Second; nr > maxNextRetry {
			nr = maxNextRetry
		}
		return nr
	}

	for i := 0; ; i = (i + 1) % len(addrs) {
		select {
		case <-time.After(nextRetry):
		case <-ctx.Done():
			t.cancel()
			return
		}

		addr := addrs[i]
		p, backups, err := t.connect(ctx, addr)
		if err != nil {
			if IsStatus(err, CodeUnauthenticated) {
				t.logger.Errorf("[JOINPOINT] terminal stopped, read %s error: %s", addr, err)
				t.cancel()
				return
			}

			nextRetry = fnNextRetry(nextRetry)
			t.logger.Warnf("[JOINPOINT] terminal wait for retry, connect %s error: %s", addr, err)
			continue
		}

		t.logger.Debugf("[JOINPOINT] terminal connect %s success", addr)
		go func() {
			start := time.Now()

			eg, egCtx := errgroup.WithContext(ctx)
			eg.Go(func() error { return t.read(egCtx, p, t.Opts.worker) })
			if err := eg.Wait(); err != nil {
				t.logger.Warnf("[JOINPOINT] terminal wait for retry, read %s error: %s", addr, err)
				nextRetry = fnNextRetry(nextRetry)
				if nextRetry = nextRetry - time.Since(start); nextRetry < 0 {
					nextRetry = 0
				}
				addrs = []string{addr}
				addrs = append(addrs, backups...)
				t.joinWithRetry(ctx, addrs, nextRetry)
			}
		}()
		break
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
