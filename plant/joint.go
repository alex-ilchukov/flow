package pit

import (
	"context"

	"github.com/alex-ilchukov/flow/errors"
)

// Joint is abstract type of logistics system, which allows to formers to
// transport their input values of type V and their output values of type W.
// Depending on the type E, the system can also allow to report on errors in
// the forming process.
type Joint[V, W any, E errors.Senders] interface {
	// Ctx should return the whole context, which the system uses to
	// operate within. It should never return nil. The context can be used
	// to propagate cancellation or deadline events.
	Ctx() context.Context

	// Get should try to receive a value within the [Ctx] context. In case
	// of success, it should resturn the value and nil. Otherwise, it
	// should return the corresponding errors from [context] package in
	// case of interruption of the transportation process or [values.Over]
	// error if receiving ability gets closed.
	Get() (v V, error)

	// Put should try to send the provided value within the [Ctx] context.
	// It should return the corresponding errors from [context] package in
	// case of interruption of the transportation process. It should return
	// nil if the transportation has been successful.
	Put(w W) error

	// Errs should return collection of error-senders. The error-senders,
	// if there are any in the collection, should operate within the same
	// [Ctx] context.
	Errs() E
}
