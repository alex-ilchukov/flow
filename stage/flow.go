package stage

import (
	"context"
	"sync"

	"github.com/alex-ilchukov/flow"
)

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

	// Spread is amount of forming go-routines, which produce or transform
	// values in concurrent way within the same [Joint] (see [Flow.Flow]
	// method). Any value less than one is interpretted as one.
	Spread int
}

// Flow takes a context and, depending on origin in the flow, does the
// following. If origin is nil, it launches emitting stage: creates logistics
// system of [Joint] abstract type and starts forming new values of type W in
// non-blocking way. If origin is not nil, it launches transforming stage:
// starts its flow of values of type V, creates logistics system of [Joint]
// abstract type, and performs transforming of the values to new values of
// type W in non-blocking way.
//
// In either case, it returns a read-only channel of the new values with
// a slice of channels of error values for reporting on possible forming
// errors. The function takes care of closing of all the channels returned and
// handles gracefully cancellation of data transportation via the provided
// context.
//
// The function supports launching of multiple go-routines for forming of
// values. The amount can be specified via [Flow.Spread] field. If the amount
// is more than one, every go-routine would have its own call to [Former.Form]
// method with the same common instance of [Joint].
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

	if f.Spread <= 1 {
		go f.form(j)
	} else {
		f.spreadForm(j)
	}

	return j.wals, rerrs
}

func (f *Flow[V, W]) form(j *joint[V, W]) {
	defer close(j.wals)
	defer close(j.errs)

	f.Former.Form(j)
}

func (f *Flow[V, W]) spreadForm(j *joint[V, W]) {
	s := spreader[V, W]{former: f.Former, joint: j}
	s.wg.Add(f.Spread)
	for i := 0; i < f.Spread; i++ {
		go s.form()
	}

	go s.wait()
}

type spreader[V, W any] struct {
	wg     sync.WaitGroup
	former Former[V, W]
	joint  *joint[V, W]
}

func (s *spreader[_, _]) form() {
	defer s.wg.Done()
	s.former.Form(s.joint)
}

func (s *spreader[_, _]) wait() {
	s.wg.Wait()
	close(s.joint.wals)
	close(s.joint.errs)
}

var (
	_ flow.Flow[int] = (*Flow[int, int])(nil)
	_ flow.Flow[int] = (*Flow[int, int])(nil)
)
