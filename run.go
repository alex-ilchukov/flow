package flow

import (
	"context"

	"github.com/alex-ilchukov/flow/errors"
	"github.com/alex-ilchukov/flow/values"
)

// Run launches the provided flow of values of type V within the provided
// context, discarding all its output values til the values-reading channel is
// closed and listening error-reading channels in blocking way. If an error
// appears in any error-reading channel, it returns it immediately. If no error
// would happen, it returns nil. In both cases it delegates proper closing of
// the channels to the flow and its user. It returns nil immediately, if the
// provided flow is nil.
func Run[V any](ctx context.Context, flow Flow[V]) error {
	if flow == nil {
		return nil
	}

	out, errs := flow.Flow(ctx)
	go values.Discard(out)

	return errors.Any(ctx, errs)
}
