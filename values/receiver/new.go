package receiver

import (
	"context"

	"github.com/alex-ilchukov/flow/values"
)

// New creates and returns implementation of receiver of values from the
// provided channel ch within context ctx.
func New[V any](ctx context.Context, ch <-chan V) *n[V] {
	return &n[V]{ctx: ctx, ch: ch}
}

type n[V any] struct {
	ctx context.Context
	ch  <-chan V
}

// Receive tries to receive value of V type. It returns the value and nil, if
// no error appears, or default value and the error otherwise. It returns
// [values.Over] if receiving is impossible anymore.
func (s *n[V]) Receive() (V, error) {
	return values.Receive(s.ctx, s.ch)
}

var _ values.Receiver[int] = (*n[int])(nil)
