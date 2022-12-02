package flow

import (
	"context"

	"github.com/alex-ilchukov/flow/chans"
	"github.com/alex-ilchukov/flow/values"
)

// Run launches the provided flow of values of type V within the provided
// context. It supports the following cases.
//
//  1. The provided flow is nil. The function returns nil immediately in this
//     case.
//  2. The provided flow is not nil, but its Flow method returns empty
//     collection of error-reading channels. In this case the function discards
//     anything from the values-reading channel in blocking way til it gets
//     closed and returns nil after.
//  3. The provided flow is not nil, and its Flow method returns non-empty
//     collection of error-reading channels. In this case the function discards
//     anything from the values-reading channel in non-blocking way til it gets
//     closed and listens for an error the error-reading channels in blocking
//     way. If an error appears in any error-reading channel, it returns it
//     immediately. If no error would happen, it returns nil.
//
// In all cases it delegates proper closing of the channels to the flow and its
// user.
func Run[V any](ctx context.Context, flow Flow[V]) error {
	if flow == nil {
		return nil
	}

	out, errs := flow.Flow(ctx)
	if len(errs) == 0 {
		values.Discard(out)
		return nil
	}

	go values.Discard(out)

	errc := chans.Merge(ctx, errs...)
	for err := range errc {
		if err != nil {
			return err
		}
	}

	return nil
}
