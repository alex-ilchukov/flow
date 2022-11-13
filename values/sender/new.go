package sender

import (
	"context"

	"github.com/alex-ilchukov/flow/values"
)

// New creates and returns implementation of sender of values to the provided
// channel ch within context ctx.
func New[V any](ctx context.Context, ch chan<- V) *n[V] {
	return &n[V]{ctx: ctx, ch: ch}
}

type n[V any] struct {
	ctx context.Context
	ch  chan<- V
}

// Send tries to push v of type V into the channel within the provided context.
// If the push is successful it returns nil, but if the push was interrupted by
// cancellation or deadline event, it returns the corresponding error.
func (s *n[V]) Send(v V) error {
	return values.Send(s.ctx, s.ch, v)
}

var _ values.Sender[int] = (*n[int])(nil)
