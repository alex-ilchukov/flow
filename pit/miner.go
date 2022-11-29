package pit

import "github.com/alex-ilchukov/flow/errors"

// Miner is abstract type of producers of values of type V.
type Miner[V any, E errors.Senders] interface {
	// Mine should produce values and put them to the provided pad.
	// Depending on type E, it can also report on errors happened during
	// production.
	Mine(Pad[V, E])
}
