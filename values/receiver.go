package values

// Receiver is abstract type of entities which receive values of type V.
type Receiver[V any] interface {
	// Receive should try to receive value of V type. It should return the
	// value and nil, if no error appears, or default value and the error
	// otherwise. It must return [Over] if receiving is impossible anymore.
	Receive() (V, error)
}
