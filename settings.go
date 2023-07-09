package retrying

import (
	"time"
)

const (
	defaultDuration = 10 * time.Second
	defaultInterval = time.Second
)

var defaultError = ErrTimeout

type settings struct {
	Duration  time.Duration
	Interval  time.Duration
	Error     error
	OnErrors  error
	MaxNumber uint
}

var defaultSettings = getBaseDefault() //nolint:gochecknoglobals

func getBaseDefault() settings {
	return settings{
		Duration:  defaultDuration,
		Interval:  defaultInterval,
		Error:     defaultError,
		OnErrors:  nil,
		MaxNumber: ^uint(0),
	}
}

func newSettings(ops ...Option) settings {
	s := defaultSettings
	for _, o := range ops {
		o(&s)
	}

	return s
}

func SetDefault(ops ...Option) {
	for _, o := range ops {
		o(&defaultSettings)
	}
}

func ResetDefault() {
	defaultSettings = getBaseDefault()
}
