package values

import "context"

// Send takes channel ch with value v and tries to push the value into the
// channel within the provided context ctx. If the push is successful it
// returns nil, but if the push was interrupted by cancellation or deadline
// event, it returns the corresponding ctx.Err() error.
func Send[V any](ctx context.Context, ch chan<- V, v V) error {
	select {
	case ch <- v:
		return nil

	case <-ctx.Done():
		return ctx.Err()
	}
}
