package pit

import (
	"context"

	"github.com/alex-ilchukov/flow"
	"github.com/alex-ilchukov/flow/errors"
)

// New takes the provided miner m of values of type V, creates a flow of values
// with the miner as its source, and returns the flow.
func New[V any, E errors.Senders](m Miner[V, E]) *impl[V, E] {
	return &impl[V, E]{m: m}
}

type impl[V any, E errors.Senders] struct {
	m Miner[V, E]
}

// Flow takes a context, creates logistics system in form of [Pad]'s instance,
// and launchs mining of values in non-blocking way. It returns a read-only
// channel of the mined values of type V with a slice of channels of error
// values for reporting on possible mining errors, if E type allows that. The
// function takes care of closing of all the channels returned and handles
// gracefully cancellation of data transportation via the provided context.
func (f *impl[V, E]) Flow(ctx context.Context) (<-chan V, []<-chan error) {
	p := paddify[V, E](ctx)
	go f.mine(p)

	return p.vals, p.rerrs
}

func (f *impl[V, E]) mine(p *pad[V, E]) {
	defer close(p.vals)
	defer errors.Close(p.werrs)

	f.m.Mine(p)
}

var (
	_ flow.Flow[int] = (*impl[int, errors.No])(nil)
	_ flow.Flow[int] = (*impl[int, errors.One])(nil)
)
