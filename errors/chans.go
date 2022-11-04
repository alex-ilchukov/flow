package errors

// No is type of zero-elements arrays of error-writing channels, effectively
// meaning that its users are simple enough to need no error handling.
type No [0]chan<- error

// One is type of one-element arrays of error-writing channels used by regular
// cases.
type One [1]chan<- error

// Chans is set of common types of error-writing channels collections. The
// zero-elements array is for those cases, which don't need error handling at
// all, and one-element array is for regular cases.
type Chans interface {
	No | One
}
