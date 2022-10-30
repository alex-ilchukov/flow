package flow

import "context"

// Collector is abstract type of entities which collect values of type V in
// non-blocking way.
type Collector[V any] interface {
	// Collect should take a context with a read-only channel of values of
	// type V and return a slice of channels of error values. It is
	// assumed, that the method takes care of closing of all the channels
	// returned and handles gracefully cancellation of data processing via
	// the provided context.
	Collect(context.Context, <-chan V) []<-chan error
}
