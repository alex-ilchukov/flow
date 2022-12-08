package stage

import "context"

// Mapper is type of generalized functions, which map values of type V to
// values of type W within the provided context. The functions can return an
// error value.
type Mapper[V, W any] func(context.Context, V) (W, error)
