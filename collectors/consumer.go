package collectors

import (
	"context"

	"github.com/alex-ilchukov/flow/errors"
)

// Consumer is type of consumers of values used by [flow.Collector] interface
// implementations.
type Consumer[V any, E errors.Chans] interface {
	// Consume should read values of type V from the provided channel with
	// respect to the provided context and consume them. Depending on
	// provided collection of error-writing channels of type from set E, it
	// can also support reporting on errors. It is not supposed to manage
	// channels and is assumed to be running in non-blocking go-routine by
	// its collecting owner.
	Consume(context.Context, <-chan V, E)
}
