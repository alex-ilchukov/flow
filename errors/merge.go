package errors

import (
	"context"

	"github.com/alex-ilchukov/flow/errors/merge"
)

// Merge takes context ctx with collection errs of read-only error channels. It
// creates new error channel and launches non-blocking concurrent reading of
// the channels, redirecting any appeared error to the new error channel with
// respect of cancellation within the provided context. It returns the error
// channel. The function takes care of closing of the error channel in distinct
// go-routine automatically, when all the original error channels are closed.
func Merge(ctx context.Context, errs ...[]<-chan error) <-chan error {
	return merge.New(ctx, errs...).Call()
}
