package flow

import (
	"context"

	"github.com/alex-ilchukov/flow/errors"
)

// Run takes emitter e with collector c and tries to create a pipeline with
// them in non-blocking way. It calls e.Emit() and c.Collect() methods just
// once within the provided context ctx. It merges the error channels, which
// the methods returned, into one error channel, and listens the channel for
// error in blocking manner. If an error appears, it returns the error,
// otherwise it returns nil. The function takes care of closing of the error
// channel.
//
// Users should usually send a cancel signal via context machinery if the
// function returns an error.
func Run[V any](ctx context.Context, e Emitter[V], c Collector[V]) error {
	values, eErrs := e.Emit(ctx)
	cErrs := c.Collect(ctx, values)

	return errors.Any(ctx, eErrs, cErrs)
}
