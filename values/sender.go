package values

// Sender is abstract type of entities which send values of type V.
type Sender[V any] interface {
	// Send should take value of V type and try to send it. It should
	// return nil, if no error appears, or the error otherwise.
	Send(v V) error
}
