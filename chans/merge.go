package chans

import (
	"context"

	"github.com/alex-ilchukov/flow/chans/merge"
)

// Merge takes context ctx with read-only channels of values of type V. It
// creates new channel and launches non-blocking concurrent reading of the
// channels, redirecting any appeared value to the new channel with respect of
// cancellation within the provided context. It returns the created channel.
// The function takes care of closing of the created channel in distinct
// go-routine automatically, when all the original channels are closed.
func Merge[V any](ctx context.Context, chs ...<-chan V) <-chan V {
	return merge.New(ctx, chs).Call()
}
