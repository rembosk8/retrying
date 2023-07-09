package retrying

import (
	"errors"
	"time"
)

type Option func(p *settings)

// WithDuration option set timeout of retrying function.
func WithDuration(d time.Duration) Option {
	return func(p *settings) {
		p.Duration = d
	}
}

// WithInterval option set interval between reps.
func WithInterval(i time.Duration) Option {
	return func(p *settings) {
		p.Interval = i
	}
}

// WithRetError option set the error which is going to be return in case of ctx.Done() instead of default ErrTimeout.
func WithRetError(err error) Option {
	return func(p *settings) {
		p.Error = err
	}
}

// WithOnError option could be used to tell a retrying function on which errors to retry. By default, retry for all errors.
func WithOnError(errs ...error) Option {
	return func(p *settings) {
		p.OnErrors = errors.Join(errs...)
	}
}

// WithMaxNumber option set the error which is going to be return in case of ctx.Done() instead of default ErrTimeout.
func WithMaxNumber(n uint) Option {
	return func(p *settings) {
		p.MaxNumber = n
	}
}
