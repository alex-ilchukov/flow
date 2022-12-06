package stage

import (
	"context"

	"github.com/alex-ilchukov/flow"
)

// New takes a flow of values of type V with a former of the values to new
// values of type W, creates new flow with the former in its core, and returns
// the flow.
func New[V, W any](flow flow.Flow[V], former Former[V, W]) *Flow[V, W] {
	return &Flow[V, W]{Origin: flow, Former: former}
}

// Flow is implementation of [flow.Flow], which can represent emitting stage
// or transforming stage, depending on its fields.
type Flow[V, W any] struct {
	// Origin is flow of original values. It is allowed to be nil, and nil
	// origin means, that the implementation is an emitting stage. If the
	// origin is not nil, the implementation is transforming stage, which
	// makes new values from original values of type V.
	Origin flow.Flow[V]

	// Former is either producer of values of type W, if [Origin] is nil,
	// or transformer of values of type V to new values of type W in other
	// case. It must not be nil.
	Former Former[V, W]
}

// Flow takes a context, launches the flow of values of type V, creates
// logistics system of [Joint] abstract type, and starts transforming of the
// values to new values of type W in non-blocking way. It returns a read-only
// channel of the transformed values with a slice of channels of error values
// for reporting on possible forming errors. The function takes care of closing
// of all the channels returned and handles gracefully cancellation of data
// transportation via the provided context.
func (f *Flow[V, W]) Flow(ctx context.Context) (<-chan W, []<-chan error) {
	if ctx == nil {
		ctx = context.Background()
	}

	j := &joint[V, W]{
		ctx:  ctx,
		wals: make(chan W),
		errs: make(chan error),
	}

	var rerrs []<-chan error
	if f.Origin != nil {
		j.vals, rerrs = f.Origin.Flow(ctx)
	}
	rerrs = append(rerrs, j.errs)

	go f.form(j)

	return j.wals, rerrs
}

func (f *Flow[V, W]) form(j *joint[V, W]) {
	defer close(j.wals)
	defer close(j.errs)

	f.Former.Form(j)
}

var (
	_ flow.Flow[int] = (*Flow[int, int])(nil)
	_ flow.Flow[int] = (*Flow[int, int])(nil)
)
