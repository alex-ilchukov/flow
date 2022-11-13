package emitters

import (
	"context"

	"github.com/alex-ilchukov/flow"
	"github.com/alex-ilchukov/flow/errors"
	"github.com/alex-ilchukov/flow/values/sender"
)

type e[V any, E errors.Senders] struct {
	p Producer[V, E]
}

// New creates an implementation of [flow.Emitter] interface with help of the
// provided producer p and returns it.
func New[V any, E errors.Senders](p Producer[V, E]) *e[V, E] {
	return &e[V, E]{p: p}
}

// Emit takes a context and launches the whole emitting process in the
// following way.
//
//  1. It creates a channel of values of type V and wraps it into value sender.
//  2. It creates as many channels of error values as there are elements in an
//     array of type E and wraps them into errors senders.
//  3. It launches a go-routine, where the senders are pushed to the producer
//     to handle data transportation.
//
// It returns the made channels immediately after go-routine starts up. It
// takes care of closing of all the channels.
func (e *e[V, E]) Emit(ctx context.Context) (<-chan V, []<-chan error) {
	out := make(chan V)
	serrs, rerrs, werrs := errors.Make[E](ctx)

	go e.produce(ctx, out, serrs, werrs)

	return out, rerrs
}

func (e *e[V, E]) produce(
	ctx context.Context,
	out chan<- V,
	serrs E,
	werrs []chan<- error,
) {

	defer errors.Close(werrs)
	defer close(out)

	s := sender.New(ctx, out)
	e.p.Produce(ctx, s, serrs)
}

var (
	_ flow.Emitter[int] = (*e[int, errors.No])(nil)
	_ flow.Emitter[int] = (*e[int, errors.One])(nil)
)
