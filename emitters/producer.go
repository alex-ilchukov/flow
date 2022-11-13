package emitters

import (
	"context"

	"github.com/alex-ilchukov/flow/errors"
	"github.com/alex-ilchukov/flow/values"
)

// Producer is type of producers of values used by [flow.Emitter] interface
// implementations.
type Producer[V any, E errors.Senders] interface {
	// Produce should produce values of V type and send them with the
	// provided sender. Depending on the provided collection of
	// error-senders of type from set E, it can also reporting on errors.
	// It can assume, that the provided instances are nevel nil.
	Produce(context.Context, values.Sender[V], E)
}
