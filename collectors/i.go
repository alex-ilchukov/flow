package collectors

import (
	"context"

	"github.com/alex-ilchukov/flow"
	"github.com/alex-ilchukov/flow/errors"
	"github.com/alex-ilchukov/flow/values/receiver"
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
//     array of type E and wraps them into error senders.
//  2. It wraps the provided channel into receiver of values.
//  3. It launches a go-routine, where the receiver and the senders are pushed
//     to the consumer to actually handle data transport.
//
// It returns the made channels immediately after go-routine starts up. It
// takes care of closing of all the channels
func (i *i[V, E]) Collect(ctx context.Context, in <-chan V) []<-chan error {
	serrs, rerrs, werrs := errors.Make[E](ctx)

	go i.consume(ctx, in, serrs, werrs)

	return rerrs
}

func (i *i[V, E]) consume(
	ctx context.Context,
	in <-chan V,
	serrs E,
	werrs []chan<- error,
) {

	defer errors.Close(werrs)

	r := receiver.New(ctx, in)
	i.c.Consume(ctx, r, serrs)
}

var (
	_ flow.Collector[int] = (*i[int, errors.No])(nil)
	_ flow.Collector[int] = (*i[int, errors.One])(nil)
)
