package pit

import (
	"context"

	"github.com/alex-ilchukov/flow/errors"
)

// Pad is abstract type of logistics system, which allows to miners to
// transport their values of type V to flow's implementation. Depending on the
// type E, the system can also allow to report on errors in the mining process.
type Pad[V any, E errors.Senders] interface {
	// Ctx should return the whole context, which the system uses to
	// operate within. It should never return nil. The context can be used
	// to propagate cancellation or deadline events.
	Ctx() context.Context

	// Put should try to transport the provided value within the [Ctx]
	// context. It should return the corresponding errors from [context]
	// package in case of interruption of the transportation process. It
	// should return nil if the transportation has been successful.
	Put(v V) error

	// Errs should return collection of error-senders. The error-senders,
	// if there are any in the collection, should operate within the same
	// [Ctx] context.
	Errs() E
}
