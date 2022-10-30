package values

import (
	"context"
	"errors"
)

// ErrClosedChannel is used to indicated that channel is closed and reading
// from it is impossible.
var ErrClosedChannel = errors.New("channel is closed")

// Receive takes channel ch and tries to pop a value from the channel within
// the provided context ctx. It return the following values.
//
//  1. If the pop is successful, it returns the value and nil.
//  2. If the channel got closed, it returns default value and
//     [ErrClosedChannel] error.
//  3. If the pop was interrupted by cancellation or deadline event, it returns
//     default value and the corresponding ctx.Err() error.
func Receive[V any](ctx context.Context, ch <-chan V) (v V, err error) {
	select {
	case v, open := <-ch:
		if open {
			return v, nil
		} else {
			return v, ErrClosedChannel
		}

	case <-ctx.Done():
		return v, ctx.Err()
	}
}