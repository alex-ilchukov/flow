package plant

import (
	"context"

	"github.com/alex-ilchukov/flow"
	"github.com/alex-ilchukov/flow/errors"
)

// New takes a flow of values of type V with a former of the values to new
// values of type W, creates new flow with the former in its core, and returns
// the flow.
func New[V, W any, E errors.Senders](
	flow   flow.Flow[V],
	former Former[V, W, E],

) *impl[V, W, E] {

	return &impl[V, W, E]{flow: flow, former: former}
}

type impl[V, W any, E errors.Senders] struct {
	flow   flow.Flow[V]
	former Former[V, W, E]
}

// Flow takes a context, launches the flow of values of type V, creates
// logistics system of [Joint] abstract type, and starts transforming of the
// values to new values of type W in non-blocking way. It returns a read-only
// channel of the transformed values with a slice of channels of error values
// for reporting on possible forming errors, if E type allows that. The
// function takes care of closing of all the channels returned and handles
// gracefully cancellation of data transportation via the provided context.
func (f *impl[V, W, E]) Flow(ctx context.Context) (<-chan W, []<-chan error) {
	if ctx == nil {
		ctx = context.Background()
	}

	j := &joint[V, W, E]{
		ctx:  ctx,
		wals: make(chan W),
	}
	j.errs, j.rerrs, j.werrs = errors.Make[E](ctx)

	var rerrs []<-chan error
	if f.flow != nil {
		j.vals, rerrs = f.flow.Flow(ctx)
	}
	rerrs = append(rerrs, j.rerrs...)

	go f.form(j)

	return j.wals, rerrs
}

func (f *impl[V, W, E]) form(j *joint[V, W, E]) {
	defer close(j.wals)
	defer errors.Close(j.werrs)

	f.former.Form(j)
}

var (
	_ flow.Flow[int] = (*impl[int, int, errors.No])(nil)
	_ flow.Flow[int] = (*impl[int, int, errors.One])(nil)
)
