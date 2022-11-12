package emitters

import (
	"context"

	"github.com/alex-ilchukov/flow"
	"github.com/alex-ilchukov/flow/errors"
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
//  1. It creates a channel of values of type V.
//  2. It creates as many channels of error values as there are elements in an
//     array of type E.
//  3. It launches a go-routine, where the channels are pushed to the producer
//     to actually handle data processing.
//
// It returns the made channels immediately after go-routine starts up. It
// takes care of closing of all the channels, but delegates graceful
// cancellation of data processing via the provided context to the producer.
func (e *e[V, E]) Emit(ctx context.Context) (<-chan V, []<-chan error) {
	out := make(chan V)
	werrs, rerrs := errors.Make[E]()

	go e.produce(ctx, out, werrs)

	return out, rerrs
}

func (e *e[V, E]) produce(ctx context.Context, out chan<- V, werrs E) {
	defer errors.Close(werrs)
	defer close(out)

	e.p.Produce(ctx, out, werrs)
}

var (
	_ flow.Emitter[int] = (*e[int, errors.No])(nil)
	_ flow.Emitter[int] = (*e[int, errors.One])(nil)
)
