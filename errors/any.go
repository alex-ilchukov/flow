package errors

import "context"

// Any merges all the error channels in the provided err parameter within
// context ctx (see errors.Merge for that) into one error channel, listens it
// in blocking manner, and returns the first error appeared in it or nil if the
// channel is closed without any error popped (and that happens when all the
// original error channels are closed).
func Any(ctx context.Context, errs ...[]<-chan error) error {
	errc := Merge(ctx, errs...)
	for err := range errc {
		if err != nil {
			return err
		}
	}

	return nil
}
