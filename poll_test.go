package retrying_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/rembosk8/retrying"
)

func TestPollWithCache(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c := 1
	fn := func(ctx2 context.Context) (int, error) {
		c++

		return c, nil
	}

	get, cancel := retrying.PollWithGetFunc(ctx, fn)
	defer cancel()
	for i := 0; i < 5; i++ {
		i2, err := get()
		fmt.Printf("i = %d with err = %v\n", i2, err)
		time.Sleep(500 * time.Millisecond)
	}
	cancel()
	time.Sleep(500 * time.Millisecond)
	i, err := get()
	fmt.Printf("i = %d with err = %v\n", i, err)

	assert.Equal(t, c, i)
	assert.Error(t, err)
}
