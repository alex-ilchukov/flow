package stage

import (
	"context"

	"github.com/alex-ilchukov/flow"
)

// New takes a flow of values of type V with a former of the values to new
// values of type W, creates new flow with the former in its core, and returns
// the flow.
func New[V, W any](flow flow.Flow[V], former Former[V, W]) *impl[V, W] {
	return &impl[V, W]{flow: flow, former: former}
}

type impl[V, W any] struct {
	flow   flow.Flow[V]
	former Former[V, W]
}

// Flow takes a context, launches the flow of values of type V, creates
// logistics system of [Joint] abstract type, and starts transforming of the
// values to new values of type W in non-blocking way. It returns a read-only
// channel of the transformed values with a slice of channels of error values
// for reporting on possible forming errors. The function takes care of closing
// of all the channels returned and handles gracefully cancellation of data
// transportation via the provided context.
func (f *impl[V, W]) Flow(ctx context.Context) (<-chan W, []<-chan error) {
	if ctx == nil {
		ctx = context.Background()
	}

	j := &joint[V, W]{
		ctx:  ctx,
		wals: make(chan W),
		errs: make(chan error),
	}

	var rerrs []<-chan error
	if f.flow != nil {
		j.vals, rerrs = f.flow.Flow(ctx)
	}
	rerrs = append(rerrs, j.errs)

	go f.form(j)

	return j.wals, rerrs
}

func (f *impl[V, W]) form(j *joint[V, W]) {
	defer close(j.wals)
	defer close(j.errs)

	f.former.Form(j)
}

var (
	_ flow.Flow[int] = (*impl[int, int])(nil)
	_ flow.Flow[int] = (*impl[int, int])(nil)
)
