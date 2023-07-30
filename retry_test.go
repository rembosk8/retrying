package retrying_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/rembosk8/retrying"
)

var errTest = errors.New("test")

func TestReturnAlwaysOK(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	f1 := func(context.Context) (int, error) {
		return 5, nil
	}

	ret, err := retrying.Return(ctx, f1)
	assert.NoError(t, err)
	assert.Equal(t, 5, ret)
}

func TestReturnSimpleBad(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	f1 := func(context.Context) (int, error) {
		var i int
		return i, errTest
	}

	ret, err := retrying.Return(ctx, f1, retrying.WithDuration(2*time.Second))
	assert.ErrorIs(t, err, errTest)
	assert.ErrorIs(t, err, retrying.ErrTimeout)
	assert.Equal(t, 0, ret)
}

func TestReturn(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	type args struct {
		fn  func(ctx context.Context) error
		ops []retrying.Option
	}
	tests := []struct {
		name     string
		args     args
		retError error
		fnError  error
	}{
		{
			name: "simple no error",
			args: args{
				fn:  func(ctx context.Context) error { return nil },
				ops: nil,
			},
			retError: nil,
			fnError:  nil,
		},
		{
			name: "simple always error until timeout",
			args: args{
				fn:  func(ctx context.Context) error { return errTest },
				ops: []retrying.Option{retrying.WithDuration(2 * time.Second)},
			},
			retError: retrying.ErrTimeout,
			fnError:  errTest,
		},
		{
			name: "interrupt",
			args: args{
				fn: func(ctx context.Context) error {
					return fmt.Errorf("we are done: %w", retrying.Interrupt(errTest))
				},
				ops: nil,
			},
			retError: retrying.ErrInterrupt,
			fnError:  errTest,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := retrying.Exec(ctx, tt.args.fn, tt.args.ops...)
			if tt.fnError != nil && !errors.Is(err, tt.fnError) {
				t.Errorf("Exec() error = %v, fnError %v", err, tt.fnError)
			}
			if tt.retError != nil && !errors.Is(err, tt.retError) {
				t.Errorf("Exec() error = %v, retError %v", err, tt.retError)
			}
		})
	}
}
