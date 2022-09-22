package lifecycle

import (
	"context"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
)

// Hook defines life-cycle's hook,
// OnStart calls on application start
// OnStop calls on application stop
type Hook struct {
	OnStart func(ctx context.Context) error
	OnStop  func(ctx context.Context) error
}

type LifeCycle struct {
	hooks []Hook
	eg    *errgroup.Group
}

func New() *LifeCycle {
	return &LifeCycle{}
}

func (l *LifeCycle) Add(hooks ...Hook) {
	l.hooks = append(l.hooks, hooks...)
}

// Start 启动程序，并发执行 start hook
// 结束的时候，顺序执行结束的 hook
func (l *LifeCycle) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	l.eg, ctx = errgroup.WithContext(ctx)

	for _, hook := range l.hooks {
		l.startHook(ctx, hook)
	}

	go func() {
		ch := make(chan os.Signal)
		signal.Notify(ch, os.Kill, os.Interrupt)
		<-ch
		cancel()
	}()

	var (
		err error
	)

	defer func() {
		for _, hook := range l.hooks {
			if err = l.stopHook(ctx, hook); err != nil {
				return
			}
		}
	}()

	if err = l.eg.Wait(); err != nil {
		return err
	}

	return nil
}

func (l *LifeCycle) startHook(ctx context.Context, hook Hook) {
	l.eg.Go(func() error {
		return hook.OnStart(ctx)
	})
}

func (l *LifeCycle) stopHook(ctx context.Context, hook Hook) error {
	if hook.OnStop != nil {
		if err := hook.OnStop(ctx); err != nil {
			return err
		}
	}
	return nil
}
