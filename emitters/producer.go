package emitters

import (
	"context"

	"github.com/alex-ilchukov/flow/errors"
)

// Producer is type of producers of values used by [flow.Emitter] interface
// implementations.
type Producer[V any, E errors.Senders] interface {
	// Produce should produce values of V type and write them into the
	// provided channel with respect to the provided context. Depending on
	// provided collection of error-writing channels of type from set E, it
	// can also support reporting on errors. It is not supposed to manage
	// channels and is assumed to be running in non-blocking go-routine by
	// its emitting owner.
	Produce(context.Context, chan<- V, E)
}
