package flow

import "context"

// Flow is interface for flow of values of any type V, which are emitted by
// emitter and collected by collector. It supposes, that the flow could be
// constructed by linking the emitter, the collector, and their error channels
// in non-blocking way. It also supposes, that the flow could be linked _and_
// run in blocking way.
type Flow[V any] interface {
	// Emitter should return the emitter of values of the type V. The
	// emitter can't be nil.
	Emitter() Emitter[V]

	// Collector should return the collector of values of the type V. The
	// collector can't be nil.
	Collector() Collector[V]

	// Link should construct the flow of values from emitter to collector
	// within the provided context in non-blocking way, merging their error
	// channels into one resulting channel. The channel can't be nil or
	// closed. The method should take care of closing the channel when the
	// original error channels are closed.
	Link(context.Context) <-chan error

	// Run should construct the flow of values from emitter to collector
	// within the provided context and listen to their error channels in
	// blocking way. If an error appears in at least one channel, it should
	// return it immediately. If the flow is over without any error, it
	// should return nil.
	Run(context.Context) error
}
