package collectors

import (
	"context"

	"github.com/alex-ilchukov/flow/errors"
	"github.com/alex-ilchukov/flow/values"
)

// Consumer is type of consumers of values used by [flow.Collector] interface
// implementations.
type Consumer[V any, E errors.Senders] interface {
	// Consume should receive values of type V wth the provided reveiver
	// and consume them. Depending on the provided collection of error
	// senders of type from set E, it can also report on errrors. It can
	// assume, that the provided instances are never nil.
	Consume(context.Context, values.Receiver[V], E)
}
