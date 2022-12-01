package stage

import "github.com/alex-ilchukov/flow/errors"

// Former is abstract type of transformers of values of type V to values of
// type W.
type Former[V, W any, E errors.Senders] interface {
	// Form should receive values of type V from the provided joint, form
	// new values of type W, and send these new values to the joint.
	// Depending on type E, it can also report on errors happened during
	// production.
	Form(Joint[V, W, E])
}
