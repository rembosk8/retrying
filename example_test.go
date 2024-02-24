package retrying_test

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rembosk8/retrying"
)

func ExampleReturn() {
	ctx := context.Background()

	i := 1
	first2CallsFailFn := func(ctx context.Context) (int, error) {
		if i > 2 {
			return 100, nil
		}
		fmt.Printf("current i is %d, not bigger than 2...\n", i)
		i++

		return 0, errors.New("not now")
	}

	res1, err := retrying.Return(ctx, first2CallsFailFn, retrying.WithInterval(1*time.Microsecond))
	fmt.Printf("Error is: %v\n", err)
	fmt.Printf("Result is: %d\n", res1)

	// Output:
	// current i is 1, not bigger than 2...
	// current i is 2, not bigger than 2...
	// Error is: <nil>
	// Result is: 100
}

func ExampleReturn_second() {
	ctx := context.Background()
	alwaysFailsFn := func(ctx context.Context) (struct{}, error) {
		fmt.Println("always fails")
		return struct{}{}, errors.New("always fails error")
	}

	res2, err := retrying.Return(
		ctx,
		alwaysFailsFn,
		retrying.WithInterval(500*time.Millisecond),
		retrying.WithDuration(1*time.Second),
		retrying.WithRetError(errors.New("i know that it fails always")),
	)
	fmt.Printf("Error is: %v\n", err)
	fmt.Printf("Result is: %v\n", res2)

	// Output:
	// always fails
	// always fails
	// always fails
	// Error is: always fails error
	// i know that it fails always
	// Result is: {}
}

func ExampleReturn_third() {
	ctx := context.Background()

	interruptFn := func(ctx context.Context) (bool, error) {
		if time.Now().After(time.Now().Add(-1 * time.Hour)) {
			fmt.Println("There is no reason to continue to retry the function")
			return false, retrying.Interrupt(nil)
		}
		fmt.Println("success")

		return true, nil
	}
	res3, err := retrying.Return(ctx, interruptFn)
	fmt.Printf("Error is: %v\n", err)
	fmt.Printf("Result is: %v\n", res3)

	// Output:
	// There is no reason to continue to retry the function
	// Error is: interrupted
	// Result is: false
}
