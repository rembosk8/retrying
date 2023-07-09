package retrying

import "errors"

var (
	// ErrInterrupt return this error from you function which you are retrying to interrupt the retrying.
	ErrInterrupt = errors.New("interrupted")
	// ErrTimeout this error is returned by default in case of deadline of retrying function is exceeded.
	ErrTimeout = errors.New("timeout")
)

func Interrupt(err error) error {
	return errors.Join(err, ErrInterrupt)
}

type errInterrupt struct {
	err error
}

func (e errInterrupt) Error() string {
	return e.err.Error()
}
