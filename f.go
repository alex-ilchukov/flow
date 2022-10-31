package flow

import (
	"context"

	"github.com/alex-ilchukov/flow/errors"
)

type f[V any] struct {
	e Emitter[V]
	c Collector[V]
}

// New returns new instance of [Flow] implementation for values of type V going
// from emitter e to collector c. It panics, if any of the provided argument is
// nil.
func New[V any](e Emitter[V], c Collector[V]) *f[V] {
	if e == nil || c == nil {
		panic("nil argument")
	}

	return &f[V]{e: e, c: c}
}

// Emitter returns emitter of values of the flow.
func (f *f[V]) Emitter() Emitter[V] {
	return f.e
}

// Collector returns collector of values of the flow.
func (f *f[V]) Collector() Collector[V] {
	return f.c
}

// Link constructs the flow of values from emitter to collector within the
// provided context in non-blocking way, merging their error channels into one
// resulting channel. The method takes care of closing the channel when the
// original error channels are closed.
func (f *f[_]) Link(ctx context.Context) <-chan error {
	values, eErrs := f.e.Emit(ctx)
	cErrs := f.c.Collect(ctx, values)

	return errors.Merge(ctx, eErrs, cErrs)
}

// Run constructs the flow of values from emitter to collector within the
// provided context and listens to their error channels in blocking way. If an
// error appears in at least one channel, it returns it immediately. If the
// flow is over without any error, it returns nil.
func (f *f[_]) Run(ctx context.Context) error {
	errs := f.Link(ctx)
	for err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}

var _ Flow[int] = (*f[int])(nil)
