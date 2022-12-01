package pit

import (
	"context"

	"github.com/alex-ilchukov/flow/errors"
	"github.com/alex-ilchukov/flow/values"
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

type pad[V any, E errors.Senders] struct {
	ctx   context.Context
	vals  chan V
	errs  E
	rerrs []<-chan error
	werrs []chan<- error
}

func (p *pad[_, _]) Ctx() context.Context {
	return p.ctx
}

func (p *pad[V, _]) Put(v V) error {
	return values.Send(p.ctx, p.vals, v)
}

func (p *pad[_, E]) Errs() E {
	return p.errs
}

var (
	_ Pad[int, errors.No]  = (*pad[int, errors.No])(nil)
	_ Pad[int, errors.One] = (*pad[int, errors.One])(nil)
)

func paddify[V any, E errors.Senders](ctx context.Context) *pad[V, E] {
	if ctx == nil {
		ctx = context.Background()
	}

	p := &pad[V, E]{
		ctx:  ctx,
		vals: make(chan V),
	}
	p.errs, p.rerrs, p.werrs = errors.Make[E](ctx)

	return p
}
