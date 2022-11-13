package emitters

import (
	"context"

	"github.com/alex-ilchukov/flow/errors"
	"github.com/alex-ilchukov/flow/values"
)

// Converters is type of converters of emitted values of type A to new values
// of type B. It is used to create [flow.Emitter] interface implementation.
type Converter[A, B any, E errors.Senders] interface {
	// Convert should receive values of type A from the receiver, convert
	// them in its way to new values of type B, and send those new values
	// with the provided sender. Depending on the provided collection of
	// error senders of type from set E, it can also report on errors. It
	// can assume, that the provided instances are nevel nil.
	Convert(context.Context, values.Receiver[A], values.Sender[B], E)
}
