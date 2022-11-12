package collectors

import (
	"context"

	"github.com/alex-ilchukov/flow"
	"github.com/alex-ilchukov/flow/errors"
)

type i[V any, E errors.Senders] struct {
	c Consumer[V, E]
}

// New creates an implementation of [flow.Collector] interface with help of the
// provided consumer c and returns it.
func New[V any, E errors.Senders](c Consumer[V, E]) *i[V, E] {
	return &i[V, E]{c: c}
}

// Collect takes a context and launches the whole collecting process in the
// following way.
//
//  1. It creates as many channels of error values as there are elements in an
//     array of type E.
//  2. It launches a go-routine, where the channels are pushed to the consumer
//     to actually handle data processing.
//
// It returns the made channels immediately after go-routine starts up. It
// takes care of closing of all the channels, but delegates graceful
// cancellation of data processing via the provided context to the consumer.
func (i *i[V, E]) Collect(ctx context.Context, in <-chan V) []<-chan error {
	werrs, rerrs := errors.Make[E]()

	go i.consume(ctx, in, werrs)

	return rerrs
}

func (i *i[V, E]) consume(ctx context.Context, in <-chan V, werrs E) {
	defer errors.Close(werrs)

	i.c.Consume(ctx, in, werrs)
}

var (
	_ flow.Collector[int] = (*i[int, errors.No])(nil)
	_ flow.Collector[int] = (*i[int, errors.One])(nil)
)
