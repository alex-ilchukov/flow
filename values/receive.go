package values

import "context"

// Receive takes channel ch and tries to pop a value from the channel within
// the provided context ctx. It returns the following values.
//
//  1. If the pop is successful, it returns the value and nil.
//  2. If the channel got closed, it returns default value and [Over] error.
//  3. If the pop was interrupted by cancellation or deadline event, it returns
//     default value and the corresponding ctx.Err() error.
//
// The function assumes, that the context is never nil, but it supports nil
// channels, interpretting them as closed ones, returning default value and
// [Over] error immediately in such a case.
func Receive[V any](ctx context.Context, ch <-chan V) (v V, err error) {
	if ch == nil {
		return v, Over
	}

	select {
	case v, open := <-ch:
		if open {
			return v, nil
		} else {
			return v, Over
		}

	case <-ctx.Done():
		return v, ctx.Err()
	}
}
