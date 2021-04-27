package group

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/oklog/run"
)

type ErrSignal struct {
	Sig os.Signal
}

func (e *ErrSignal) Error() string {
	return fmt.Sprintf("signal %v", e.Sig)
}

type Group struct {
	run.Group
}

func (g *Group) With(ctx context.Context) *Group {
	ctx, cancel := context.WithCancel(ctx)
	g.Add(func() error {
		<-ctx.Done()
		return ctx.Err()
	}, func(error) { cancel() })
	return g
}

func (g *Group) Signal() *Group {
	ctx, cancel := context.WithCancel(context.Background())
	g.Add(SignalHandler(ctx), func(error) { cancel() })
	return g
}

func SignalHandler(ctx context.Context) func() error {
	return func() error {
		sch := make(chan os.Signal)
		signal.Notify(sch, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(sch)
		select {
		case sig := <-sch:
			return &ErrSignal{Sig: sig}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
