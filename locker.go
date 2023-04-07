package locker

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
)

var (
	ErrAlreadyLocked = errors.New("already locked")
	ErrUnlocked      = errors.New("explicitly unlocked")
)

type cancelCauseFunc = func(cause error)

type Lock struct {
	lock          *redislock.Lock
	ctx           context.Context
	cancelRefresh cancelCauseFunc
	refreshWg     *sync.WaitGroup
}

func (m *Lock) Unlock() error {
	m.cancelRefresh(ErrUnlocked)

	m.refreshWg.Wait()

	if err := m.lock.Release(m.ctx); err != nil {
		return fmt.Errorf("can't release lock: %w", err)
	}

	return nil
}

type Locker struct {
	client *redislock.Client
	config *config
}

func (m *Locker) runRefresh(ctx context.Context, l *redislock.Lock, cfg *config, cancel cancelCauseFunc) {
	for {
		select {
		case <-time.After(cfg.refreshPeriod):
		case <-ctx.Done():
			return
		}

		if ctx.Err() != nil {
			return
		}

		if err := l.Refresh(ctx, cfg.ttl, nil); err != nil {
			cancel(fmt.Errorf("can't refresh lock: %w", err))

			return
		}
	}
}

func (m *Locker) Lock(ctx context.Context, key string, options ...Option) (*Lock, context.Context, error) {
	cfg := m.config
	if len(options) > 0 {
		var err error
		cfg, err = makeConfigWithOptions(cfg, options)
		if err != nil {
			return nil, nil, err
		}
	}

	l, err := m.client.Obtain(ctx, key, cfg.ttl, nil)
	if err != nil {
		if errors.Is(err, redislock.ErrNotObtained) {
			return nil, ctx, fmt.Errorf("%w: %s", ErrAlreadyLocked, err) //nolint:errorlint
		}

		return nil, nil, fmt.Errorf("can't obtain lock: %w", err)
	}

	refreshCtx, cancel := withCancelCause(ctx)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		m.runRefresh(refreshCtx, l, cfg, cancel)
	}()

	return &Lock{
		lock:          l,
		ctx:           ctx,
		cancelRefresh: cancel,
		refreshWg:     wg,
	}, refreshCtx, nil
}

func New(client redis.UniversalClient, options ...Option) (*Locker, error) {
	cfg, err := makeConfigWithOptions(defaultConfig, options)
	if err != nil {
		return nil, err
	}

	return &Locker{
		client: redislock.New(client),
		config: cfg,
	}, nil
}
