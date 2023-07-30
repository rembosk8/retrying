package retrying

import (
	"context"
	"errors"
	"sync"
	"time"
)

// ErrChan a chanel for handling errors from polling function.
type ErrChan <-chan error

// ValChan a chanel for handling value result from polling function.
type ValChan[T any] <-chan T

// GetFunc function to get results (value of type T and error) from polling function.
type GetFunc[T any] func() (T, error)

// CancelFunc function to cancel polling.
type CancelFunc = context.CancelFunc

// Poll is calling function fn N seconds ( 1 second by default, could be changed with WithInterval option) until
//
//	fn returns ErrInterrupt or CancelFunc is called.
func Poll(ctx context.Context, fn func(ctx context.Context) error, ops ...Option) CancelFunc {
	p := newSettings(ops...)
	ctx, cancel := context.WithCancel(ctx)

	defer cancel()

	ticker := time.NewTicker(p.Interval)
	defer ticker.Stop()

	go func() {
		for {
			err := fn(ctx)
			if errors.Is(err, ErrInterrupt) {
				return
			}

			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				continue
			}
		}
	}()

	return cancel
}

func updateChan[T any](v T, ch chan T) {
	select {
	case <-ch:
	default:
	}
	ch <- v
}

func PollReturn[T any](ctx context.Context, fn func(ctx context.Context) (T, error), ops ...Option) (ValChan[T], ErrChan, CancelFunc) {
	ch := make(chan T, 1)
	er := make(chan error, 1)
	ctx, cancel := context.WithCancel(ctx)

	defer cancel()

	go func() {
		var (
			t   T
			err error
		)
		p := newSettings(ops...)
		ticker := time.NewTicker(p.Interval)
		defer ticker.Stop()
		defer close(ch)
		defer close(er)
		for {
			t, err = fn(ctx)
			if err == nil {
				updateChan(t, ch)
			}
			if errors.Is(err, ErrInterrupt) {
				updateChan(err, er)
				return
			}

			select {
			case <-ctx.Done():
				updateChan(ctx.Err(), er)
				return
			case <-ticker.C:
				continue
			}
		}
	}()

	return ch, er, cancel
}

func PollWithGetFunc[T any](
	ctx context.Context,
	fn func(ctx context.Context) (T, error),
	ops ...Option,
) (get GetFunc[T], cancel CancelFunc) {
	var mux sync.RWMutex
	var t T
	var err error

	ctx, cancel = context.WithCancel(ctx)
	go func() {
		p := newSettings(ops...)
		ticker := time.NewTicker(p.Interval)
		defer cancel()
		defer ticker.Stop()

		for {
			tt, errt := fn(ctx)
			mux.Lock()
			t, err = tt, errt
			mux.Unlock()
			if errors.Is(err, ErrInterrupt) {
				return
			}

			select {
			case <-ctx.Done():
				err = errors.Join(err, ctx.Err())
				return
			case <-ticker.C:
				continue
			}
		}
	}()

	get = func() (T, error) {
		mux.RLock()
		defer mux.RUnlock()
		return t, err
	}

	return get, cancel
}
