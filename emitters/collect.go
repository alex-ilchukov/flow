package emitters

import (
	"context"

	"github.com/alex-ilchukov/flow"
	"github.com/alex-ilchukov/flow/collectors"
	"github.com/alex-ilchukov/flow/errors"
	"github.com/alex-ilchukov/flow/values/receiver"
)

type col[V any, E errors.Senders] struct {
	e flow.Emitter[V]
	c collectors.Consumer[V, E]
}

// Collect creates an instance of [flow.Emitter], which is based on the
// provided emitter and just consumes values with help of the provided
// consumer. The created instance emits no values and uses [Unit] type for the
// output.
func Collect[V any, E errors.Senders](
	e flow.Emitter[V],
	c collectors.Consumer[V, E],

) *col[V, E] {

	return &col[V, E]{e: e, c: c}
}

// Emit takes a context and launches the whole emitting process in the
// following way.
//
//  1. It calls Emit() method of the provided emitter, getting a channel to
//     read values of type V from and a slice of error-reading channels.
//  2. It creates a channel of Unit values, but the channel is a fake and for
//     return only.
//  3. It creates as many channels of error values as there are elements in an
//     array of type E and wraps them into errors senders.
//  4. It launches a go-routine, where the values-reading channel is wrapped to
//     a receiver, and pushes the receiver and the senders to the consumer to
//     handle data transportation.
//
// It returns the fake channel and combined error-reading channels immediately
// after go-routine starts up. It takes care of closing of all the channels.
func (c *col[V, E]) Emit(ctx context.Context) (<-chan Unit, []<-chan error) {
	v, errs := c.e.Emit(ctx)
	u := make(chan Unit)
	serrs, rerrs, werrs := errors.Make[E](ctx)

	go c.consume(ctx, v, u, serrs, werrs)

	return u, append(errs, rerrs...)
}

func (c *col[V, E]) consume(
	ctx context.Context,
	v <-chan V,
	u chan<- Unit,
	serrs E,
	werrs []chan<- error,
) {

	defer errors.Close(werrs)
	defer close(u)

	r := receiver.New(ctx, v)
	c.c.Consume(ctx, r, serrs)
}

var (
	_ flow.Emitter[Unit] = (*col[int, errors.No])(nil)
	_ flow.Emitter[Unit] = (*col[int, errors.One])(nil)
)
