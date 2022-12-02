package stage

import (
	"context"

	"github.com/alex-ilchukov/flow/errors"
	"github.com/alex-ilchukov/flow/values"
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
	Get() (V, error)

	// Put should try to send the provided value within the [Ctx] context.
	// It should return the corresponding errors from [context] package in
	// case of interruption of the transportation process. It should return
	// nil if the transportation has been successful.
	Put(W) error

	// Report should try to report on the provided error within the [Ctx]
	// context. It should return the corresponding errors from [context]
	// package in case of interruption of the transportation process. It
	// should return nil if the transportation has been successful.
	Report(error) error
}

type joint[V, W any, E errors.Senders] struct {
	ctx  context.Context
	vals <-chan V
	wals chan W
	errs chan error
}

func (j *joint[_, _, _]) Ctx() context.Context {
	return j.ctx
}

func (j *joint[V, _, _]) Get() (V, error) {
	return values.Receive(j.ctx, j.vals)
}

func (j *joint[_, W, _]) Put(w W) error {
	return values.Send(j.ctx, j.wals, w)
}

func (j *joint[_, _, _]) Report(e error) error {
	return values.Send(j.ctx, j.errs, e)
}

var (
	_ Joint[int, int, errors.No]  = (*joint[int, int, errors.No])(nil)
	_ Joint[int, int, errors.One] = (*joint[int, int, errors.One])(nil)
)
