package retrying

import "errors"

var (
	// ErrInterrupt return this error from you function which you are retrying to interrupt the retrying.
	ErrInterrupt = errors.New("interrupted")
	// ErrTimeout this error is returned by default in case of deadline of retrying function is exceeded.
	ErrTimeout = errors.New("timeout")
)

// Interrupt wrap method to force retrying to stop.
//
//	err - (optional) initial error.
func Interrupt(err error) error {
	if err == nil {
		return interruptError{err: ErrInterrupt}
	}
	return interruptError{err: err}
}

type interruptError struct {
	err error
}

func (e interruptError) Error() string {
	return e.err.Error()
}
