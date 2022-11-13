package emitters

import (
	"context"

	"github.com/alex-ilchukov/flow"
	"github.com/alex-ilchukov/flow/errors"
	"github.com/alex-ilchukov/flow/values/receiver"
	"github.com/alex-ilchukov/flow/values/sender"
)

type t[A, B any, E errors.Senders] struct {
	e flow.Emitter[A]
	c Converter[A, B, E]
}

// Transform creates an instance of [flow.Emitter], which is based on the
// provided emitter and converter of values.
func Transform[A, B any, E errors.Senders](
	e flow.Emitter[A],
	c Converter[A, B, E],

) *t[A, B, E] {

	return &t[A, B, E]{e: e, c: c}
}

// Emit takes a context and launches the whole emitting process in the
// following way.
//
//  1. It calls Emit() method of the provided emitter, getting a channel to
//     read values of type A from (a-channel) and a slice of error-reading
//     channels.
//  2. It creates a channel of values of type B (b-channel).
//  3. It creates as many channels of error values as there are elements in an
//     array of type E and wraps them into errors senders.
//  4. It launches a go-routine, where the a-channel is wrapped to a receiver,
//     b-channel is wrapped to a sender, and pushes the receiver and the
//     senders to the converter to handle data transportation.
//
// It returns the b-channel and combined error-reading channels immediately
// after go-routine starts up. It takes care of closing of all the channels.
func (t *t[A, B, E]) Emit(ctx context.Context) (<-chan B, []<-chan error) {
	a, aerrs := t.e.Emit(ctx)
	b := make(chan B)
	serrs, rerrs, werrs := errors.Make[E](ctx)

	go t.convert(ctx, a, b, serrs, werrs)

	return b, append(aerrs, rerrs...)
}

func (t *t[A, B, E]) convert(
	ctx context.Context,
	a <-chan A,
	b chan<- B,
	serrs E,
	werrs []chan<- error,
) {

	defer errors.Close(werrs)
	defer close(b)

	r := receiver.New(ctx, a)
	s := sender.New(ctx, b)
	t.c.Convert(ctx, r, s, serrs)
}

var (
	_ flow.Emitter[int] = (*t[int, int, errors.No])(nil)
	_ flow.Emitter[int] = (*t[int, int, errors.One])(nil)
)
