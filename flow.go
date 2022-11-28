package flow

import "context"

// Flow is abstract type of data flows with non-blocking processing.
type Flow[V any] interface {
	// Flow should take a context, launch its processing in non-blocking
	// way, and return a read-only channel of result values of type V with
	// a slice of channels of error values. It is assumed, that the method
	// takes care of closing of all the channels returned and handles
	// gracefully cancellation of data processing via the provided context.
	Flow(context.Context) (<-chan V, []<-chan error)
}
