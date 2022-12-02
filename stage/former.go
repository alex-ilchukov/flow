package stage

// Former is abstract type of transformers of values of type V to values of
// type W.
type Former[V, W any] interface {
	// Form should receive values of type V from the provided joint, form
	// new values of type W, and send these new values to the joint.
	Form(Joint[V, W])
}
