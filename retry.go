package retrying

import (
	"context"
	"errors"
	"time"
)

// Exec call function fn every 1 sec for 10 sec by default until fn doesn't return error.
//
//	Default behavior for timeout and interval could be changed via Options.
//	To force exit from retry loop fn should return ErrInterrupt error.
func Exec(ctx context.Context, fn func(ctx context.Context) error, ops ...Option) error {
	f := func(ctx context.Context) (struct{}, error) {
		return struct{}{}, fn(ctx)
	}
	_, err := Return(ctx, f, ops...)
	return err
}

// Return call function fn every 1 sec for 10 sec by default until fn doesn't return error and return fn result.
//
//	 If you need to return many values then combine them into a struct.
//		Default behavior for timeout and interval could be changed via Options.
//		To force exit from retry loop fn should wrap/return retrying.Interrupt(err).
func Return[T any](ctx context.Context, fn func(ctx context.Context) (T, error), ops ...Option) (t T, err error) {
	var intErr interruptError
	p := newSettings(ops...)
	ctx, cancel := context.WithTimeout(ctx, p.Duration)
	defer cancel()

	ticker := time.NewTicker(p.Interval)
	defer ticker.Stop()

	for i := uint(0); i < p.MaxNumber; i++ {
		t, err = fn(ctx)
		if err == nil {
			return t, nil
		}
		// in case err is ErrInterrupt - stop
		if errors.As(err, &intErr) {
			return t, intErr.err
		}
		// in case OnErrors is not empty and err not in OnErrors - stop
		if p.OnErrors != nil && !errors.Is(err, p.OnErrors) {
			return t, err
		}

		select {
		case <-ctx.Done():
			return t, errors.Join(err, p.Error)
		case <-ticker.C:
			continue
		}
	}

	return t, err
}
