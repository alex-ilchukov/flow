package flow

import "context"

// Emitter is abstract type of entities which emit values of type V in
// non-blocking way.
type Emitter[V any] interface {
	// Emit should take a context and return a read-only channel of values
	// of type V with a slice of channels of error values. It is assumed,
	// that the method takes care of closing of all the channels returned
	// and handles gracefully cancellation of data processing via the
	// provided context.
	Emit(context.Context) (<-chan V, []<-chan error)
}
